package job

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/apps/server/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/apps/server/internal/dto"
	errs "github.com/slugger7/exorcist/apps/server/internal/errors"
	"github.com/slugger7/exorcist/apps/server/internal/ffmpeg"
	"github.com/slugger7/exorcist/apps/server/internal/models"
)

func CreateGenerateChaptersJob(mediaId uuid.UUID, jobId *uuid.UUID, interval *float64, height int, width int, maxDimension int, overwite bool) (*model.Job, error) {
	d := dto.GenerateChaptersData{
		MediaId:      mediaId,
		Height:       new(int),
		Width:        new(int),
		MaxDimension: maxDimension,
		Overwrite:    overwite,
	}
	*d.Height = height
	*d.Width = width

	if interval == nil {
		d.Interval = 60
	}

	js, err := json.Marshal(d)
	if err != nil {
		return nil, errs.BuildError(err, "could not marshall generate chapters data")
	}

	data := string(js)
	job := &model.Job{
		JobType:  model.JobTypeEnum_GenerateChapters,
		Status:   model.JobStatusEnum_NotStarted,
		Data:     &data,
		Parent:   jobId,
		Priority: dto.JobPriority_MediumLow,
	}

	return job, nil
}

func (jr *jobRunner) removeChapters(id uuid.UUID, chapters []models.MediaRelation) error {
	var accErr error
	for _, i := range chapters {
		if err := jr.service.Media().Delete(i.RelatedTo, true); err != nil {
			accErr = errors.Join(accErr, err)
		}

		if err := jr.repo.Media().RemoveRelation(id, i.RelatedTo); err != nil {
			accErr = errors.Join(accErr, err)
		}
	}

	m, err := jr.repo.Media().GetById(id)
	if err != nil {
		accErr = errors.Join(accErr, err)
		return accErr
	}

	// TODO: this can be simplified seeing as we only need the relations
	mediaDto := new(dto.MediaDTO).FromModel(*m)

	mediaUpdate := dto.MediaDTO{
		Relations: mediaDto.Relations,
	}

	jr.ws.MediaUpdate(mediaUpdate)

	return accErr
}

func (jr *jobRunner) generateChapters(job *model.Job) error {
	var jobData dto.GenerateChaptersData
	if err := json.Unmarshal([]byte(*job.Data), &jobData); err != nil {
		return errs.BuildError(err, "error parsing job data for generate chapters: %v", job.Data)
	}

	media, err := jr.repo.Media().GetById(jobData.MediaId)
	if err != nil {
		return errs.BuildError(err, "could not find media by id for generate chapters job %v", jobData.MediaId.String())
	}

	if media == nil {
		return fmt.Errorf("media was nil for generate chapters job: %v", jobData.MediaId.String())
	}

	relations := []models.MediaRelation{}
	for _, relation := range media.MediaRelations {
		if relation.RelationType == model.MediaRelationTypeEnum_Chapter {
			relations = append(relations, relation)
		}
	}

	if len(relations) > 0 {
		if jobData.Overwrite {
			if err := jr.removeChapters(media.Media.ID, relations); err != nil {
				jr.logger.Warningf("some issues removing previous chapters: %v", err.Error())
			}
		} else {
			jr.logger.Infof("chapters already exist for %v as it already has chapters and overwrite was set to false", jobData.MediaId)
			return nil
		}
	}

	if media.Video == nil {
		return fmt.Errorf("media was not of type video: %v", jobData.MediaId.String())
	}

	runtimeDuration := time.Duration(int64(media.Video.Runtime * float64(time.Second)))
	intervalDuration := time.Duration(int64(jobData.Interval * float64(time.Second)))

	relationType := model.MediaRelationTypeEnum_Chapter

	if jobData.Height == nil {
		jobData.Height = new(int)
	}
	if jobData.Width == nil {
		jobData.Width = new(int)
	}

	if *jobData.Height == 0 {
		*jobData.Height = int(media.Video.Height)
	}

	if *jobData.Width == 0 {
		*jobData.Width = int(media.Video.Width)
	}

	if jobData.MaxDimension != 0 {
		if *jobData.Width > jobData.MaxDimension {
			*jobData.Height = ffmpeg.ScaleHeightByWidth(*jobData.Height, *jobData.Width, jobData.MaxDimension)
			*jobData.Width = jobData.MaxDimension
		}

		if *jobData.Height > jobData.MaxDimension {
			*jobData.Width = ffmpeg.ScaleWidthByHeight(*jobData.Height, *jobData.Width, jobData.MaxDimension)
			*jobData.Height = jobData.MaxDimension
		}
	}

	generateThumbnailJobs := []model.Job{}
	var accErr error
	for i := intervalDuration; i < runtimeDuration; i += intervalDuration {
		metadata := dto.ThumbnailMetadataDTO{
			Timestamp: i.Seconds(),
		}

		assetPath := filepath.Join(
			jr.env.Assets,
			media.Media.ID.String(),
			fmt.Sprintf(
				"%v.%v.%vx%v.%v.webp",
				filepath.Base(media.Media.Path),
				relationType.String(),
				*jobData.Height,
				*jobData.Width,
				i,
			))
		job, err := CreateGenerateThumbnailJob(media.Media.ID, &job.ID, assetPath, i.Seconds(), *jobData.Height, *jobData.Width, &relationType, &metadata)
		if err != nil {
			accErr = errors.Join(accErr, err)
			continue
		}

		generateThumbnailJobs = append(generateThumbnailJobs, *job)
	}

	if accErr != nil {
		jr.logger.Errorf("encountered while creating generate thumbnail jobs: %v", accErr.Error())
	}

	if len(generateThumbnailJobs) != 0 {
		if _, err := jr.repo.Job().CreateAll(generateThumbnailJobs); err != nil {
			return errs.BuildError(err, "creating generate thumbnail jobs")
		}
	}

	return nil
}
