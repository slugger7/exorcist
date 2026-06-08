package personService

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/apps/server/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/apps/server/internal/dto"
	"github.com/slugger7/exorcist/apps/server/internal/environment"
	errs "github.com/slugger7/exorcist/apps/server/internal/errors"
	"github.com/slugger7/exorcist/apps/server/internal/logger"
	"github.com/slugger7/exorcist/apps/server/internal/models"
	"github.com/slugger7/exorcist/apps/server/internal/repository"
)

type PersonService interface {
	Upsert(name string) (*model.Person, error)
	GetMedia(id, userId uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error)
	Delete(id uuid.UUID) error
}

type personService struct {
	env    *environment.EnvironmentVariables
	repo   repository.Repository
	logger logger.Logger
}

func (p *personService) Delete(id uuid.UUID) error {
	person, err := p.repo.Person().GetById(id)
	if err != nil {
		return errs.BuildError(err, "could not find person by id")
	}
	if person == nil {
		return fmt.Errorf("person with id %v does not exist", id.String())
	}

	if err := p.repo.Person().Delete(id); err != nil {
		return errs.BuildError(err, "error deleting person by id")
	}

	return nil
}

func (p *personService) GetMedia(id, userId uuid.UUID, search dto.MediaSearchDTO) (*dto.PageDTO[models.MediaOverviewModel], error) {
	person, err := p.repo.Person().GetById(id)
	if err != nil {
		return nil, errs.BuildError(err, "could not get person by id from repo: %v", id)
	}

	if person == nil {
		return nil, fmt.Errorf("no person found with id: %v", id)
	}

	media, err := p.repo.Person().GetMedia(id, userId, search)
	if err != nil {
		return nil, errs.BuildError(err, "colud not get media by person id from repo: %v", id)
	}

	return media, nil
}

// Upsert implements IPersonService.
func (p *personService) Upsert(name string) (*model.Person, error) {
	person, err := p.repo.Person().GetByName(name)
	if err != nil {
		return nil, errs.BuildError(err, "could not get person by name from repo")
	}

	if person == nil {
		people, err := p.repo.Person().Create([]string{name})
		if err != nil {
			return nil, errs.BuildError(err, "could not create person by name")
		}
		person = &people[0]
	}

	return person, nil
}

var personServiceInstance *personService

func New(repo repository.Repository, env *environment.EnvironmentVariables) PersonService {
	if personServiceInstance == nil {
		personServiceInstance = &personService{
			env:    env,
			repo:   repo,
			logger: logger.New(env),
		}

		personServiceInstance.logger.Info("PersonService instance created")
	}

	return personServiceInstance
}
