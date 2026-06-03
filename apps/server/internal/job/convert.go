package job

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/slugger7/exorcist/apps/server/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/apps/server/internal/dto"
	errs "github.com/slugger7/exorcist/apps/server/internal/errors"
	"github.com/slugger7/exorcist/apps/server/internal/ffmpeg"
)

const (
	CONVERT_FOLDER_NAME string = "conversions"
)

func (jr *JobRunner) convert(job *model.Job) error {
	var jobData dto.ConvertData
	if err := json.Unmarshal([]byte(*job.Data), &jobData); err != nil {
		return errs.BuildError(err, "error parssing job data for convert: %v", job.Data)
	}

	media, err := jr.repo.Media().GetById(jobData.MediaId)
	if err != nil {
		return errs.BuildError(err, "error fetching media")
	}

	if media == nil {
		return fmt.Errorf("no media found with id: %v", jobData.MediaId.String())
	}

	jr.logger.Infof("Converting video %v to %v", media.Path, jobData.Filename)

	tempPath := filepath.Join(jr.env.Cache, CONVERT_FOLDER_NAME, media.Media.ID.String())

	if _, err := os.Stat(jobData.Path); err == nil {
		return fmt.Errorf("Path for converted media already exsists: %v", jobData.Path)
	}

	convertData := jobData.ToFfmpegDto()
	convertData.InputFilePath = media.Path
	convertData.OutputFilePath = tempPath

	if err = ffmpeg.Convert(*convertData); err != nil {
		return errs.BuildError(err, "conversion failed")
	}

	// TODO Add to library
	// TODO Link to existing media

	return nil
}

func copyFile(original, destination string) error {
	fin, err := os.Open(original)
	if err != nil {
		return errs.BuildError(err, "colud not open original file: %v", original)
	}
	defer fin.Close()

	fout, err := os.Create(destination)
	if err != nil {
		return errs.BuildError(err, "could not create destination: %v", destination)
	}
	defer fout.Close()

	_, err = io.Copy(fout, fin)
	if err != nil {
		currentErr := errs.BuildError(err, "could not copy %v to %v, cleaning up", original, destination)

		err := os.Remove(destination)
		if err != nil {
			return errs.BuildError(errors.Join(currentErr, err), "could not remove file at %v", destination)
		}

		return currentErr
	}

	return nil
}
