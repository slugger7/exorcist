package repoTestHelpers

import (
	"context"
	"database/sql"
	"log"

	"github.com/slugger7/exorcist/apps/server/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/apps/server/internal/db/exorcist/public/table"
)

func CreateLibraryPath(ctx context.Context, db *sql.DB, libPath model.LibraryPath) model.LibraryPath {
	libraryPath := table.LibraryPath
	stmnt := libraryPath.INSERT(libraryPath.LibraryID, libraryPath.Path).
		MODEL(libPath).
		RETURNING(libraryPath.AllColumns)

	var createdLibPath model.LibraryPath
	if err := stmnt.QueryContext(ctx, db, &createdLibPath); err != nil {
		log.Fatal(err)
	}

	return createdLibPath
}

func CreateStubLibraryPath(ctx context.Context, db *sql.DB) model.LibraryPath {
	lib := CreateStubLibrary(ctx, db)
	return CreateLibraryPath(ctx, db, model.LibraryPath{LibraryID: lib.ID, Path: "stub path"})
}
