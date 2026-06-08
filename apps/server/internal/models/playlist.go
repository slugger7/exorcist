package models

import "github.com/slugger7/exorcist/apps/server/internal/db/exorcist/public/model"

type Playlist struct {
	model.Playlist
	Media model.PlaylistMedia
}
