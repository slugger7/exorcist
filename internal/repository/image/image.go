package imageRepository

import (
	"database/sql"
	"fmt"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
)

type IImageRepository interface {
	Create(m *model.Image) (*model.Image, error)
	RelateVideo(m *model.VideoImage) (*model.VideoImage, error)
}

type ImageRepository struct {
	db  *sql.DB
	Env *environment.EnvironmentVariables
}

var imageRepoInstance *ImageRepository

func New(db *sql.DB, env *environment.EnvironmentVariables) IImageRepository {
	if imageRepoInstance == nil {
		imageRepoInstance = &ImageRepository{
			db:  db,
			Env: env,
		}
	}
	return imageRepoInstance
}

func (i *ImageRepository) Create(m *model.Image) (*model.Image, error) {
	var imgs []struct{ model.Image }
	if err := i.createStatement(m).Query(&imgs); err != nil {
		return nil, errs.BuildError(err, "error creating image")
	}

	if len(imgs) == 1 {
		return &imgs[len(imgs)-1].Image, nil
	}

	return nil, fmt.Errorf("no images were returned from query")
}

func (i *ImageRepository) RelateVideo(m *model.VideoImage) (*model.VideoImage, error) {
	var vidImgs []struct{ model.VideoImage }
	if err := i.relateVideoStatement(m).Query(&vidImgs); err != nil {
		return nil, errs.BuildError(err, "error creating relation between video and image")
	}

	if len(vidImgs) == 1 {
		return &vidImgs[len(vidImgs)-1].VideoImage, nil
	}

	return nil, fmt.Errorf("no video image relations were returned from query")
}
