package jobService

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/apps/server/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/apps/server/internal/dto"
	"github.com/slugger7/exorcist/apps/server/internal/environment"
	errs "github.com/slugger7/exorcist/apps/server/internal/errors"
	"github.com/slugger7/exorcist/apps/server/internal/ffmpeg"
	"github.com/slugger7/exorcist/apps/server/internal/logger"
	"github.com/slugger7/exorcist/apps/server/internal/media"
	"github.com/slugger7/exorcist/apps/server/internal/repository"
	mediaService "github.com/slugger7/exorcist/apps/server/internal/service/media"
)

type JobService interface {
	Create(dto.CreateJobDTO) (*model.Job, error)
	StartJobRunner()
}

type jobService struct {
	env          *environment.EnvironmentVariables
	repo         repository.Repository
	logger       logger.Logger
	jobCh        chan bool
	ctx          context.Context
	mediaService mediaService.MediaService
}

var jobServiceInstance *jobService

func New(repo repository.Repository, env *environment.EnvironmentVariables, jobCh chan bool, ctx context.Context, mediaService mediaService.MediaService) JobService {
	if jobServiceInstance == nil {
		jobServiceInstance = &jobService{
			env:          env,
			repo:         repo,
			logger:       logger.New(env),
			jobCh:        jobCh,
			ctx:          ctx,
			mediaService: mediaService,
		}

		jobServiceInstance.logger.Info("UserService instance created")
	}
	return jobServiceInstance
}

func (s *jobService) Create(m dto.CreateJobDTO) (*model.Job, error) {
	defaultJobPriority := dto.JobPriority_Medium
	if m.Priority == nil {
		m.Priority = &(defaultJobPriority)
	}
	data, err := json.Marshal(m.Data)
	strData := string(data)
	if err != nil {
		return nil, errs.BuildError(err, "could not marhsal data field")
	}
	var j *model.Job
	var e error
	switch m.Type {
	case model.JobTypeEnum_ScanPath:
		j, e = s.scanPath(strData, *m.Priority)
	case model.JobTypeEnum_GenerateThumbnail:
		j, e = s.generateThumbnail(strData, *m.Priority)
	case model.JobTypeEnum_RefreshMetadata:
		j, e = s.refreshMetadata(strData, *m.Priority)
	case model.JobTypeEnum_RefreshLibraryMetadata:
		j, e = s.refreshLibraryMetadata(strData, *m.Priority)
	case model.JobTypeEnum_GenerateChapters:
		j, e = s.generateChapters(strData, *m.Priority)
	case model.JobTypeEnum_Convert:
		j, e = s.convert(strData, *m.Priority)
	default:
		return nil, fmt.Errorf("job type not implemented: %v", m.Type)
	}
	if e != nil {
		return nil, errs.BuildError(err, "error encountered while creating job")
	}

	job := model.Job{
		JobType:  m.Type,
		Status:   model.JobStatusEnum_NotStarted,
		Data:     j.Data,
		Priority: j.Priority,
	}

	jobs, err := s.repo.Job().CreateAll([]model.Job{job})
	if err != nil {
		return nil, errs.BuildError(err, "creating job")
	}

	if len(jobs) == 0 {
		return nil, fmt.Errorf("no jobs were returned after creating a job")
	}

	go s.StartJobRunner()

	return &jobs[0], nil
}

func (i *jobService) convert(data string, priority int16) (*model.Job, error) {
	var jobData dto.ConvertData
	if err := json.Unmarshal([]byte(data), &jobData); err != nil {
		return nil, errs.BuildError(err, "unmarshalling data for convert: %v", data)
	}

	priority = dto.JobPriority_Lowest

	media, err := i.repo.Media().GetById(jobData.MediaId)
	if err != nil {
		return nil, errs.BuildError(err, "getting media by id: %v", jobData.MediaId.String())
	}

	if media == nil {
		return nil, fmt.Errorf("no media with id: %v", jobData.MediaId.String())
	}

	if !media.Exists {
		return nil, fmt.Errorf("media does not exist: %v", jobData.MediaId)
	}

	if media.Deleted {
		i.logger.Warningf("Conversion requested on deleted media: %v", jobData.MediaId)
	}

	// TODO: probably need some more filepath sanitization
	filePath := filepath.Clean(
		filepath.Join(
			filepath.Dir(media.Path),
			filepath.Base(jobData.Filename),
		))

	if _, err := os.Stat(filePath); err == nil {
		return nil, fmt.Errorf("destination for conversion already existst: %v", filePath)
	}

	jobData.Path = filePath

	currentDimension := ffmpeg.Dimension{
		Height: new(int),
		Width:  new(int),
	}
	*currentDimension.Width = int(media.Video.Width)
	*currentDimension.Height = int(media.Video.Height)

	dimension := ffmpeg.DetermineDimensions(*jobData.Dimension.ToFfmpegDto(), currentDimension)

	jobData.Dimension = *new(dto.Dimension).FromFfmpegDto(&dimension)

	bytes, err := json.Marshal(jobData)
	if err != nil {
		return nil, errs.BuildError(err, "could not remarshall convert data")
	}

	data = string(bytes)
	return &model.Job{
		Data:     &data,
		Priority: priority,
	}, nil
}

func (i *jobService) generateChapters(data string, priority int16) (*model.Job, error) {
	var jobData dto.GenerateChaptersData
	if err := json.Unmarshal([]byte(data), &jobData); err != nil {
		return nil, errs.BuildError(err, "unmarshalling data for generate chapters data: %v", data)
	}

	if jobData.Interval == 0 {
		jobData.Interval = float64(((time.Minute * 5).Seconds()))
	}

	media, err := i.repo.Media().GetById(jobData.MediaId)
	if err != nil {
		return nil, errs.BuildError(err, "getting media by id: %v", jobData.MediaId.String())
	}

	if media == nil {
		return nil, fmt.Errorf("no media with id: %v", jobData.MediaId.String())
	}

	if media.Video == nil {
		return nil, fmt.Errorf("media is not of type video: %v", jobData.MediaId.String())
	}

	bytes, err := json.Marshal(jobData)
	if err != nil {
		return nil, errs.BuildError(err, "could not remarshall generate chapters data")
	}

	data = string(bytes)

	return &model.Job{
		Data:     &data,
		Priority: priority,
	}, nil
}

func (i *jobService) refreshLibraryMetadata(data string, priority int16) (*model.Job, error) {
	var jobData dto.RefreshLibraryMetadata
	if err := json.Unmarshal([]byte(data), &jobData); err != nil {
		return nil, errs.BuildError(err, "unmarshalling data for refresh library metadata: %v", data)
	}

	if jobData.RefreshFields == nil {
		jobData.RefreshFields = &dto.RefreshFields{
			Size:     true,
			Checksum: false,
		}
	}

	library, err := i.repo.Library().GetById(jobData.LibraryId)
	if err != nil {
		return nil, errs.BuildError(err, "getting library by id: %v", jobData.LibraryId.String())
	}

	if library == nil {
		return nil, fmt.Errorf("no library found with id: %v", jobData.LibraryId.String())
	}

	bytes, err := json.Marshal(jobData)
	if err != nil {
		return nil, errs.BuildError(err, "remarshalling refresh library metadata")
	}

	data = string(bytes)

	return &model.Job{
		Data:     &data,
		Priority: priority,
	}, nil
}

func (i *jobService) refreshMetadata(data string, priority int16) (*model.Job, error) {
	var jobData dto.RefreshMetadata
	if err := json.Unmarshal([]byte(data), &jobData); err != nil {
		return nil, errs.BuildError(err, "unmarshalling data for refresh meta data: %v", data)
	}

	mediaEntity, err := i.repo.Media().GetById(jobData.MediaId)
	if err != nil {
		return nil, errs.BuildError(err, "fetching media entity by id: %v", jobData.MediaId.String())
	}

	if mediaEntity == nil {
		return nil, fmt.Errorf("no media entity found to refresh the metada of: %v", jobData.MediaId.String())
	}

	return &model.Job{
		Data:     &data,
		Priority: priority,
	}, nil

}

func (i *jobService) removeExistingThumbnail(mediaId uuid.UUID) error {
	media, err := i.repo.Media().GetById(mediaId)
	if err != nil {
		return errs.BuildError(err, "could not get assets for media id %v", mediaId)
	}

	for _, m := range media.MediaRelations {
		if m.MediaRelation.RelationType == model.MediaRelationTypeEnum_Thumbnail {
			if err := i.mediaService.Delete(m.MediaRelation.RelatedTo, true); err != nil {
				return errs.BuildError(err, "could not remove thumbnail information")
			}

			if err := i.repo.Media().RemoveRelation(m.MediaRelation.MediaID, m.MediaRelation.RelatedTo); err != nil {
				return errs.BuildError(
					err,
					"could not remove thumbnail relation: %v related to %v",
					m.MediaRelation.MediaID,
					m.MediaRelation.RelatedTo)
			}
		}
	}

	return nil
}

const ErrActionGenerateThumbnailVideoNotFound = "could not find video for generate thumbnail job: %v"

func (i *jobService) generateThumbnail(data string, priority int16) (*model.Job, error) {
	var generateThumbnailData dto.GenerateThumbnailData

	if err := json.Unmarshal([]byte(data), &generateThumbnailData); err != nil {
		return nil, errs.BuildError(err, "could not unmarshal data for job %v", data)
	}

	if generateThumbnailData.RelationType == nil {
		v := model.MediaRelationTypeEnum_Thumbnail
		generateThumbnailData.RelationType = &v
	}

	m, err := i.repo.Video().GetByMediaId(generateThumbnailData.MediaId)
	if err != nil {
		return nil, errs.BuildError(
			err,
			ErrActionGenerateThumbnailVideoNotFound,
			generateThumbnailData.MediaId)
	}

	f, err := media.GetFileInformation(m.Media.Path)
	if err != nil {
		return nil, errs.BuildError(err, "could not get file information")
	}

	if err := i.removeExistingThumbnail(generateThumbnailData.MediaId); err != nil {
		return nil, errs.BuildError(err, "could not remove existing thumbnail")
	}

	w := ffmpeg.Dimension{
		Height: generateThumbnailData.Height,
		Width:  generateThumbnailData.Width,
	}

	c := ffmpeg.Dimension{
		Height: new(int),
		Width:  new(int),
	}
	*c.Height = int(m.Video.Height)
	*c.Width = int(m.Video.Width)

	d := ffmpeg.DetermineDimensions(w, c)

	if generateThumbnailData.Height == nil {
		generateThumbnailData.Height = new(int)
	}
	if generateThumbnailData.Width == nil {
		generateThumbnailData.Width = new(int)
	}
	*generateThumbnailData.Height = *d.Height
	*generateThumbnailData.Width = *d.Width

	generateThumbnailData.Path = filepath.Join(
		i.env.Assets,
		generateThumbnailData.MediaId.String(),
		fmt.Sprintf(
			`%v.%v.%vx%v.webp`,
			f.FileName,
			generateThumbnailData.RelationType.String(),
			*generateThumbnailData.Height,
			*generateThumbnailData.Width,
		))

	bytes, err := json.Marshal(generateThumbnailData)
	if err != nil {
		return nil, errs.BuildError(err, "could not remarshal generate thumbnail data")
	}

	data = string(bytes)

	return &model.Job{
		Data:     &data,
		Priority: priority,
	}, nil

	// Generates the thumbnail but does not remove the previous thumbnail record.
	// The query to fetch the video with the thumbnail is joining with multiple results duplicating the record at the moment
}

const ErrActionScanGetLibraryPaths = "could not get library paths in scan action"
const ErrCreatingJobs = "error creating jobs"

func (i *jobService) scanPath(data string, priority int16) (*model.Job, error) {
	var scanPathData dto.ScanPathData

	if err := json.Unmarshal([]byte(data), &scanPathData); err != nil {
		return nil, errs.BuildError(err, "could not unmarshall data for job %v", data)
	}

	// check to see if the library path actually exists before creating a job
	_, err := i.repo.LibraryPath().GetById(scanPathData.LibraryPathId)
	if err != nil {
		return nil, errs.BuildError(err, ErrActionScanGetLibraryPaths)
	}

	return &model.Job{
		Data:     &data,
		Priority: priority,
	}, nil
}

// We do this at the moment to stack a signal to the job runner if it is already running
func (i *jobService) StartJobRunner() {
	i.logger.Debug("Starting a job runner")
	select {
	case <-i.ctx.Done():
		i.logger.Debug("Shutdown signal recieved. Not starting job runner")
		return
	default:
		i.logger.Debug("Starting job runner")
		if i.env.JobRunner {
			i.jobCh <- true
			i.logger.Debug("Job runner signal sent")
		}
	}
}
