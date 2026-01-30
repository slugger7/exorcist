package dto

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/apps/server/internal/models"
)

type MediaRelationDto struct {
	MediaID     uuid.UUID `json:"mediaId"`
	ThumbnailID uuid.UUID `json:"thumbnailId"`
}

func (d *MediaRelationDto) FromModel(m models.MediaRelation) MediaRelationDto {
	d.MediaID = m.Media.ID
	d.ThumbnailID, _ = uuid.NewRandom()

	return *d
}
