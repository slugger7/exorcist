package models

import (
	"github.com/slugger7/exorcist/apps/server/internal/db/exorcist/public/model"
)

type MediaVideo struct {
	model.Media
	model.Video
}

type MediaRelation struct {
	model.MediaRelation
}

type MediaOverviewModel struct {
	model.Media
	model.MediaProgress
	*model.Video
	*model.FavouriteMedia
}

type Media struct {
	model.Media
	*model.Image
	*model.Video
	*model.MediaProgress
	*model.FavouriteMedia
	People         []model.Person
	Tags           []model.Tag
	MediaRelations []MediaRelation
}

type MediaImage struct {
	model.Image
	model.Media
}
