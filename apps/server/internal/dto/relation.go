package dto

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/apps/server/internal/models"
)

type RelationDTO struct {
	ID uuid.UUID `json:"id"`
}

func (o *RelationDTO) FromModel(m *models.Relation) *RelationDTO {
	o.ID = m.ID

	return o
}
