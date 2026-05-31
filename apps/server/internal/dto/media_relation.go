package dto

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/apps/server/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/apps/server/internal/models"
)

type MediaRelationDto struct {
	RelatedToID  uuid.UUID                   `json:"relatedToId"`
	RelationType model.MediaRelationTypeEnum `json:"relationType"`
	// Needs to be any as we do not know what type the metadata is at this point
	// however the Relation type does tell us and can be used to cast the data
	// to the correct type if needed.
	// This is only used to give the client a full json object without them needing
	// to parse the json string
	Metadata any `json:"metadata"`
}

func (d *MediaRelationDto) FromModel(m models.MediaRelation) MediaRelationDto {
	d.RelatedToID = m.RelatedTo
	d.RelationType = m.RelationType

	switch d.RelationType {
	case model.MediaRelationTypeEnum_Chapter:
		var chapterMetadata ChapterMetadadataDTO
		if e := json.Unmarshal([]byte(*m.Metadata), &chapterMetadata); e != nil {
			slog.Error("failed to unmarshall chapter metadata", "error", e.Error())
			fmt.Println(e.Error())
			d.Metadata = nil
		} else {
			d.Metadata = chapterMetadata
		}
	case model.MediaRelationTypeEnum_Thumbnail:
		var thumbnailMetadata ThumbnailMetadataDTO
		if e := json.Unmarshal([]byte(*m.Metadata), &thumbnailMetadata); e != nil {
			slog.Error("failed to unmarshall thumbnail metadata", "error", e.Error())
			fmt.Println(e.Error())
			d.Metadata = nil
		} else {
			d.Metadata = thumbnailMetadata
		}
	}

	return *d
}
