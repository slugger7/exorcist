package mediaService

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/apps/server/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/apps/server/internal/dto"
	"github.com/slugger7/exorcist/apps/server/internal/environment"
	errs "github.com/slugger7/exorcist/apps/server/internal/errors"
	"github.com/slugger7/exorcist/apps/server/internal/logger"
	"github.com/slugger7/exorcist/apps/server/internal/models"
	"github.com/slugger7/exorcist/apps/server/internal/repository"
	personService "github.com/slugger7/exorcist/apps/server/internal/service/person"
	tagService "github.com/slugger7/exorcist/apps/server/internal/service/tag"
)

type mediaService struct {
	env           *environment.EnvironmentVariables
	repo          repository.Repository
	logger        logger.Logger
	personService personService.PersonService
	tagService    tagService.TagService
}

type MediaService interface {
	AddTag(id uuid.UUID, tagId uuid.UUID) (*model.MediaTag, error)
	AddPerson(id uuid.UUID, personId uuid.UUID) (*model.MediaPerson, error)
	Delete(id uuid.UUID, physical bool) error
	LogProgress(id, userId uuid.UUID, progress dto.ProgressUpdateDTO) (*model.MediaProgress, error)
	GetByIdAndUserIdWithRelations(id, userId uuid.UUID, relationType *model.MediaRelationTypeEnum) (*models.Media, error)
	Relate(id uuid.UUID, relateDto dto.PutMediaRelationDto) ([]model.MediaRelation, error)
	CopyTags(toId, fromId uuid.UUID) error
	CopyPeople(toId, fromId uuid.UUID) error
	DeleteRelations(id uuid.UUID, deleteDto dto.DeleteMediaRelationsDto) error
}

func createRelations(id uuid.UUID, relationDto dto.PutMediaRelationDto) []model.MediaRelation {
	relations := []model.MediaRelation{}
	for _, r := range relationDto.RelatedToIDs {
		relations = append(relations, model.MediaRelation{
			MediaID:      id,
			RelatedTo:    r,
			RelationType: model.MediaRelationTypeEnum_Media,
		})

		if relationDto.Backrelate {
			relations = append(relations, model.MediaRelation{
				MediaID:      r,
				RelatedTo:    id,
				RelationType: model.MediaRelationTypeEnum_Media,
			})
		}

		if relationDto.Interrelate {
			for _, r1 := range relationDto.RelatedToIDs {
				if r1 != r {
					relations = append(relations, model.MediaRelation{
						MediaID:      r,
						RelatedTo:    r1,
						RelationType: model.MediaRelationTypeEnum_Media,
					})
				}
			}
		}
	}

	return relations
}

func (s *mediaService) Relate(id uuid.UUID, relateDto dto.PutMediaRelationDto) ([]model.MediaRelation, error) {
	media, err := s.repo.Media().GetById(id)
	if err != nil {
		return nil, errs.BuildError(err, "error fetching media")
	}

	if media == nil {
		return nil, fmt.Errorf("No media found for id: %v", id)
	}

	relations := createRelations(id, relateDto)

	relationModels, err := s.repo.Media().Relate(relations)
	if err != nil {
		return nil, errs.BuildError(err, "error relating media")
	}

	return relationModels, nil
}

// DeleteRelations implements [MediaService].
func (s *mediaService) DeleteRelations(id uuid.UUID, deleteDto dto.DeleteMediaRelationsDto) error {
	media, err := s.repo.Media().GetById(id)
	if err != nil {
		return errs.BuildError(err, "could not get media by id")
	}
	if media == nil {
		return fmt.Errorf("could not find media by id: %v", id.String())
	}

	if err := s.repo.Media().DeleteRelations(id, deleteDto); err != nil {
		return errs.BuildError(err, "could not delete relations")
	}

	return nil
}

// CopyPeople implements [MediaService].
func (s *mediaService) CopyPeople(toId uuid.UUID, fromId uuid.UUID) error {
	toMedia, err := s.repo.Media().GetById(toId)
	if err != nil {
		return errs.BuildError(err, "error getting media to copy people to")
	}
	if toMedia == nil {
		return fmt.Errorf("could not find media to copy people to: %v", toId.String())
	}

	fromMedia, err := s.repo.Media().GetById(fromId)
	if err != nil {
		return errs.BuildError(err, "error getting media to copy people from")
	}
	if fromMedia == nil {
		return fmt.Errorf("could not find media to compy people from: %v", fromId.String())
	}

	var accErrs error
	for _, f := range fromMedia.People {
		if _, err := s.AddPerson(toId, f.ID); err != nil {
			if accErrs == nil {
				accErrs = errs.BuildError(err, "could not copy tag from %v", fromId.String())
			} else {

			}
		}
	}

	return nil
}

// CopyTags implements [MediaService].
func (s *mediaService) CopyTags(toId uuid.UUID, fromId uuid.UUID) error {
	toMedia, err := s.repo.Media().GetById(toId)
	if err != nil {
		return errs.BuildError(err, "error getting media to copy tags to")
	}
	if toMedia == nil {
		return fmt.Errorf("could not find media to copy tags to: %v", toId.String())
	}

	fromMedia, err := s.repo.Media().GetById(fromId)
	if err != nil {
		return errs.BuildError(err, "error getting media to copy tags from")
	}
	if fromMedia == nil {
		return fmt.Errorf("could not find media to compy tags from: %v", fromId.String())
	}

	for _, f := range fromMedia.Tags {
		s.AddTag(toId, f.ID)
	}

	return nil
}

func (s *mediaService) GetByIdAndUserIdWithRelations(id, userId uuid.UUID, relationType *model.MediaRelationTypeEnum) (*models.Media, error) {
	m, err := s.repo.Media().GetByIdAndUserId(id, userId)
	if err != nil {
		return nil, errs.BuildError(err, "")
	}

	r, err := s.repo.Media().GetRelationsFor(id, relationType)
	if err != nil {
		return nil, errs.BuildError(err, "")
	}

	m.MediaRelations = r

	return m, nil
}

// LogProgress implements MediaService.
func (m *mediaService) LogProgress(id, userId uuid.UUID, progress dto.ProgressUpdateDTO) (*model.MediaProgress, error) {
	current, err := m.repo.Media().GetProgressForUser(id, userId)
	if err != nil {
		if !progress.Overwrite {
			return nil, errs.BuildError(err, "could not fetch progress for user %v and video %v", userId.String(), id.String())
		}
	}

	if progress.Progress < 0 {
		m, err := m.repo.Media().GetById(id)
		if err != nil {
			return nil, errs.BuildError(err, "could not get media by id")
		}

		progress.Progress = m.Runtime + progress.Progress
	}

	if current != nil && !progress.Overwrite {
		if current.Timestamp > progress.Progress {
			return current, nil
		}
	}

	prog := &model.MediaProgress{
		UserID:    userId,
		MediaID:   id,
		Timestamp: progress.Progress,
	}

	newProg, err := m.repo.Media().UpsertProgress(*prog)
	if err != nil {
		return nil, errs.BuildError(err, "could not upsert progress for in repo")
	}

	return newProg, nil
}

// Delete implements MediaService.
func (m *mediaService) Delete(id uuid.UUID, physical bool) error {
	mediaEntity, err := m.repo.Media().GetById(id)
	if err != nil {
		return errs.BuildError(err, "could not find media by id: %v", id.String())
	}

	if mediaEntity == nil {
		return fmt.Errorf("media entity with id (%v) does not exist", id.String())
	}

	assets, err := m.repo.Media().GetAssetsFor(id)
	if err != nil {
		return errs.BuildError(err, "could not find assets for: %v", id.String())
	}

	if physical {
		// BE CAREFUL IN THIS SECTION, FILES ARE DELETED OFF DISK
		assetsPath := path.Join(m.env.Assets, mediaEntity.Media.ID.String())
		if err = os.RemoveAll(assetsPath); err != nil {
			m.logger.Errorf("could not remove assets and assets folder (%v): %v", assetsPath, err.Error())
		}

		if err = os.Remove(mediaEntity.Media.Path); err != nil {
			m.logger.Errorf("could not remove media (%v): %v", mediaEntity.Path, err.Error())
		}

		mediaEntity.Media.Exists = false
	}

	mediaEntity.Media.Deleted = true
	for _, a := range assets {
		mediaModel := model.Media{
			ID:      a.MediaID,
			Deleted: true,
			Exists:  !physical,
		}

		if err := m.repo.Media().Delete(mediaModel); err != nil {
			return errs.BuildError(err, "something failed while deleting an asset (%v) in repo: %v", a.MediaID.String(), id.String())
		}
	}

	if err := m.repo.Media().Delete(mediaEntity.Media); err != nil {
		return errs.BuildError(err, "something failed while deleting media in repo: %v", id.String())
	}

	return nil
}

// AddPerson implements MediaService.
func (m *mediaService) AddPerson(id uuid.UUID, personId uuid.UUID) (*model.MediaPerson, error) {
	mediaModel, err := m.repo.Media().GetById(id)
	if err != nil {
		return nil, errs.BuildError(err, "could not get media by id from repo: %v", id.String())
	}

	if mediaModel == nil {
		return nil, fmt.Errorf("could not find media by id: %v", id)
	}

	personModel, err := m.repo.Person().GetById(personId)
	if err != nil {
		return nil, errs.BuildError(err, "could not get person by id from repo: %v", personId.String())
	}

	if personModel == nil {
		return nil, fmt.Errorf("colud not find person by id: %v", personId)
	}

	for _, p := range mediaModel.People {
		if p.ID == personId {
			return &model.MediaPerson{MediaID: id, PersonID: personId}, nil
		}
	}

	mediaPeopleModels := []model.MediaPerson{
		{
			PersonID: personId,
			MediaID:  id,
		},
	}
	createdMediaModelPeople, err := m.repo.Person().AddToMedia(mediaPeopleModels)
	if err != nil {
		return nil, errs.BuildError(err, "could not add person (%v) to media (%v)", personId, id)
	}

	return &createdMediaModelPeople[0], nil
}

// AddTag implements MediaService.
func (m *mediaService) AddTag(id uuid.UUID, tagId uuid.UUID) (*model.MediaTag, error) {
	mediaModel, err := m.repo.Media().GetById(id)
	if err != nil {
		return nil, errs.BuildError(err, "could not get media by id from repo: %v", id.String())
	}

	if mediaModel == nil {
		return nil, fmt.Errorf("could not find media by id: %v", id)
	}

	tagModel, err := m.repo.Tag().GetById(tagId)
	if err != nil {
		return nil, errs.BuildError(err, "could not get tag by id from repo: %v", id.String())
	}

	if tagModel == nil {
		return nil, fmt.Errorf("could not find tag by id: %v", tagId)
	}

	for _, t := range mediaModel.Tags {
		if t.ID == tagId {
			return &model.MediaTag{MediaID: id, TagID: tagId}, nil
		}
	}

	mediaTagModels := []model.MediaTag{
		{
			TagID:   tagId,
			MediaID: id,
		},
	}
	createdMediaModelTags, err := m.repo.Tag().AddToMedia(mediaTagModels)
	if err != nil {
		return nil, errs.BuildError(err, "could not add tag (%v) to media (%v)", tagId, id)
	}

	return &createdMediaModelTags[0], nil
}

func lowerStringComparator(a string) func(string) bool {
	return func(b string) bool {
		return strings.ToLower(a) == strings.ToLower(b)
	}
}

var mediaServiceInstance *mediaService

func New(env *environment.EnvironmentVariables, repo repository.Repository, personService personService.PersonService, tagService tagService.TagService) MediaService {
	if mediaServiceInstance == nil {
		mediaServiceInstance = &mediaService{
			env:           env,
			repo:          repo,
			logger:        logger.New(env),
			personService: personService,
			tagService:    tagService,
		}

		mediaServiceInstance.logger.Info("Created media service instance")
	}

	return mediaServiceInstance
}
