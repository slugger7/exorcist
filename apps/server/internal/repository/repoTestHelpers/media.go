package repoTestHelpers

import (
	"context"
	"database/sql"
	"log"

	"github.com/slugger7/exorcist/apps/server/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/apps/server/internal/db/exorcist/public/table"
)

func CreateMedia(ctx context.Context, db *sql.DB, mediaModel model.Media) model.Media {
	media := table.Media
	stmnt := media.INSERT(
		media.LibraryPathID,
		media.Path,
		media.MediaType,
		media.Title,
		media.Size,
	).
		MODEL(mediaModel).
		RETURNING(media.AllColumns)

	var createdMediaModel model.Media
	if err := stmnt.QueryContext(ctx, db, &createdMediaModel); err != nil {
		log.Fatal(err)
	}

	return createdMediaModel
}

func CreateStubMedia(ctx context.Context, db *sql.DB) model.Media {
	libPath := CreateStubLibraryPath(ctx, db)
	return CreateMedia(ctx, db, model.Media{
		LibraryPathID: libPath.ID,
		Path:          "stub",
		MediaType:     model.MediaTypeEnum_Primary,
		Title:         "stub",
		Size:          69420,
	})
}
