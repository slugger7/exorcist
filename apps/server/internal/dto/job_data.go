package dto

import (
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/apps/server/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/apps/server/internal/ffmpeg"
)

type ScanPathData struct {
	LibraryPathId uuid.UUID `json:"libraryPathId"`
}

type GenerateThumbnailData struct {
	MediaId uuid.UUID `json:"mediaId"`
	Path    string    `json:"path"`
	// Optional: If set to 0, timestamp at 25% of video playback will be used. Value in seconds
	Timestamp float64 `json:"timestamp"`
	// Optional: If set to 0, video height will be used
	Height int `json:"height"`
	// Optional: If set to 0, video widtch will be used
	Width        int                          `json:"width"`
	RelationType *model.MediaRelationTypeEnum `json:"relationType"`
	Metadata     *ThumbnailMetadataDTO        `json:"metadata"`
}

type RefreshFields struct {
	Size     bool `json:"size"`
	Checksum bool `json:"checksum"`
}

type RefreshMetadata struct {
	MediaId       uuid.UUID      `json:"mediaId"`
	RefreshFields *RefreshFields `json:"refreshFields"`
}

type RefreshLibraryMetadata struct {
	LibraryId     uuid.UUID      `json:"libraryId"`
	BatchSize     int            `json:"batchSize"`
	RefreshFields *RefreshFields `json:"refreshFields"`
}

type GenerateChaptersData struct {
	MediaId      uuid.UUID `json:"mediaId"`
	Interval     float64   `json:"interval"`
	Height       int       `json:"height"`
	Width        int       `json:"width"`
	MaxDimension int       `json:"maxDimension"`
	Overwrite    bool      `json:"overwrite"`
}

type ConvertData struct {
	MediaId            uuid.UUID `json:"mediaId" binding:"required"`
	Dimension          Dimension `json:"dimension"`
	Filename           string    `json:"filename" binding:"required"`
	Path               string    `json:"path" tstype:"-"` // omitted for clients
	ConstantRateFactor *int      `json:"constantRateFactor"`
	VariableBitrate    *int      `json:"variableBitrate"`
	ForcePixelFormat   *string   `json:"forcePixelFormat"`
}

func (d *ConvertData) ToFfmpegDto() *ffmpeg.ConvertDto {
	v := &ffmpeg.ConvertDto{
		Dimension: ffmpeg.Dimension{
			Height: new(int),
			Width:  new(int),
		},
		ConstantRateFactor: d.ConstantRateFactor,
		VariableBitrate:    d.VariableBitrate,
		ForcePixelFormat:   d.ForcePixelFormat,
	}
	*v.Dimension.Height = *d.Dimension.Height
	*v.Dimension.Width = *d.Dimension.Width

	return v
}
