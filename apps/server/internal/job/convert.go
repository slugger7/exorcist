package job

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/apps/server/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/apps/server/internal/dto"
	errs "github.com/slugger7/exorcist/apps/server/internal/errors"
	"github.com/slugger7/exorcist/apps/server/internal/ffmpeg"
	"github.com/slugger7/exorcist/apps/server/internal/media"
	"github.com/slugger7/exorcist/apps/server/internal/models"
)

const (
	CONVERT_FOLDER_NAME string = "conversions"
)

func (jr *jobRunner) convert(job *model.Job) error {
	var jobData dto.ConvertData
	if err := json.Unmarshal([]byte(*job.Data), &jobData); err != nil {
		return errs.BuildError(err, "error parssing job data for convert: %v", job.Data)
	}

	mediaModel, err := jr.repo.Media().GetById(jobData.MediaId)
	if err != nil {
		return errs.BuildError(err, "error fetching media")
	}

	if mediaModel == nil {
		return fmt.Errorf("no media found with id: %v", jobData.MediaId.String())
	}

	jr.logger.Infof("Converting video %v to %v", mediaModel.Path, jobData.Filename)

	tempPath := filepath.Join(jr.env.Cache, CONVERT_FOLDER_NAME, mediaModel.Media.ID.String())
	tempFilePath := filepath.Join(tempPath, jobData.Filename)

	if _, err := os.Stat(jobData.Path); err == nil {
		return fmt.Errorf("Path for converted media already exsists: %v", jobData.Path)
	}

	// TODO: figure out the file permissions
	os.MkdirAll(tempPath, 0777)

	convertData := jobData.ToFfmpegDto()
	convertData.InputFilePath = mediaModel.Path
	convertData.OutputFilePath = tempFilePath

	if err = ffmpeg.Convert(*convertData); err != nil {
		return errs.BuildError(err, "conversion failed")
	}

	createdMedia, createdVideo, err := jr.addStubMedia(*mediaModel, jobData.Path, tempFilePath)
	if err != nil {
		return errs.BuildError(err, "error creating media")
	}
	if createdMedia == nil {
		return fmt.Errorf("created media is null")
	}

	if err := media.CopyFile(tempFilePath, jobData.Path); err != nil {
		return errs.BuildError(err, "error copying file")
	}

	if err := os.Remove(tempFilePath); err != nil {
		return errs.BuildError(err, "error removing temp file from %v", tempFilePath)
	}

	if _, err := jr.service.Media().Relate(jobData.MediaId, dto.PutMediaRelationDto{
		RelatedToIDs: []uuid.UUID{createdMedia.ID},
		Backrelate:   true,
		Interrelate:  false,
	}); err != nil {
		jr.logger.Errorf("could not relate converted media with id %v to original media with id %v in job %v: %v",
			createdMedia.ID.String(), jobData.MediaId.String(), job.ID.String(), err.Error())
	}

	jobs := createNewMediaJobs(&job.ID, *createdMedia, *createdVideo, jr.env.Assets)

	_, err = jr.repo.Job().CreateAll(jobs)
	if err != nil {
		return errs.BuildError(err, "could not create jobs for convert job: %v", job.ID.String())
	}

	return nil
}

func createNewMediaJobs(jobId *uuid.UUID, newMedia model.Media, newVideo model.Video, assetPath string) []model.Job {
	jobs := []model.Job{}

	checksumJob, err := CreateGenerateChecksumJob(newMedia.ID, jobId)
	if err != nil {
		slog.Warn("could not create checksum job", "jobId", jobId.String())
	}
	if checksumJob != nil {
		jobs = append(jobs, *checksumJob)
	}

	relationType := model.MediaRelationTypeEnum_Thumbnail

	height := int(newVideo.Height)
	width := int(newVideo.Width)
	dimension := ffmpeg.Dimension{
		Height: &height,
		Width:  &width,
	}
	scaledDimension := ffmpeg.ScaleByMaxDimension(maxDimension, dimension)

	thumbnailPath := filepath.Join(
		assetPath,
		newMedia.ID.String(),
		fmt.Sprintf(
			`%v.%v.%vx%v.webp`,
			filepath.Base(newMedia.Path),
			relationType.String(),
			*scaledDimension.Height,
			*scaledDimension.Width,
		))

	thumbnailJob, err := CreateGenerateThumbnailJob(newMedia.ID, jobId,
		thumbnailPath, 0, *scaledDimension.Height, *scaledDimension.Width, &relationType, nil)
	if err != nil {
		slog.Warn("could not create generate thumbnail job", "jobId", jobId.String())
	}
	if thumbnailJob != nil {
		jobs = append(jobs, *thumbnailJob)
	}

	chaptersJob, err := CreateGenerateChaptersJob(newMedia.ID, jobId,
		nil, *scaledDimension.Height, *scaledDimension.Width, maxDimension, false)
	if err != nil {
		slog.Warn("could not create generate chapters job", "jobId", jobId.String())
	}
	if chaptersJob != nil {
		jobs = append(jobs, *chaptersJob)
	}

	return jobs
}

func (jr *jobRunner) addStubMedia(existing models.Media, filePath, tempPath string) (*model.Media, *model.Video, error) {
	libraryPath, err := jr.repo.LibraryPath().GetContainingPath(existing.Media.Path)
	if err != nil {
		return nil, nil, errs.BuildError(err, "could not get library path containing: %v", existing.Media.Path)
	}
	if len(libraryPath) == 0 {
		return nil, nil, fmt.Errorf("library path was nil for %v", existing.Media.Path)
	}

	fileSize, err := media.GetFileSize(tempPath)
	if err != nil {
		return nil, nil, errs.BuildError(err, "could not get file size")
	}

	newMediaModel := model.Media{
		LibraryPathID: libraryPath[len(libraryPath)-1].ID,
		Title:         existing.Title,
		Size:          fileSize,
		Path:          filePath,
		MediaType:     model.MediaTypeEnum_Primary,
	}

	createdMedia, err := jr.repo.Media().Create([]model.Media{newMediaModel})
	if err != nil {
		return nil, nil, errs.BuildError(err, "could not create stub media for %v", filePath)
	}
	if len(createdMedia) != 1 {
		return nil, nil, fmt.Errorf("incorrect amount of media returned at creation: count %v for %v", len(createdMedia), filePath)
	}

	createdVideo, err := jr.addVideo(tempPath, createdMedia[len(createdMedia)-1].ID)
	if err != nil {
		return nil, nil, errs.BuildError(err, "could not create video for media")
	}

	return &createdMedia[len(createdMedia)-1], createdVideo, nil
}

func (jr *jobRunner) addVideo(path string, mediaId uuid.UUID) (*model.Video, error) {
	ffmpegData, err := ffmpeg.UnmarshalledProbe(path)
	if err != nil {
		return nil, errs.BuildError(err, "could not get probe for %v", path)
	}

	dimension, err := ffmpeg.GetDimensions(ffmpegData.Streams)
	if err != nil {
		jr.logger.Warningf("could not extract dimensions for %v. Setting to 0. Reason: %v", path, err.Error())
		*dimension.Height = 0
		*dimension.Width = 0
	}

	runtime, err := strconv.ParseFloat(ffmpegData.Format.Duration, 32)
	if err != nil {
		jr.logger.Warningf(
			"could not convert duration from string (%v) to float for video %v. Setting runtime to 0",
			ffmpegData.Format.Duration, path)
		runtime = 0.0
	}

	newVideoModel := model.Video{
		MediaID: mediaId,
		Height:  int32(*dimension.Height),
		Width:   int32(*dimension.Width),
		Runtime: runtime,
	}

	createdVideos, err := jr.repo.Video().Insert([]model.Video{newVideoModel})
	if err != nil {
		return nil, errs.BuildError(err, "could not create video")
	}
	if len(createdVideos) != 1 {
		return nil, fmt.Errorf("no videos returned after creation")
	}

	return &createdVideos[len(createdVideos)-1], nil
}
