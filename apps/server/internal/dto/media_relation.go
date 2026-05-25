package dto

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/apps/server/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/apps/server/internal/models"
)

type MediaRelationDto struct {
	RelatedToID  uuid.UUID                   `json:"relatedToId"`
	RelationType model.MediaRelationTypeEnum `json:"relationType"`
	// TODO: consider parsing this to the actual meta data type
	Metadata *string `json:"metadata"`
}

func (d *MediaRelationDto) FromModel(m models.MediaRelation) MediaRelationDto {
	d.RelatedToID = m.RelatedTo
	d.RelationType = m.RelationType
	d.Metadata = m.Metadata

	return *d
}
