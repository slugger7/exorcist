package repoTestHelpers

import (
	"context"
	"database/sql"
	"log"

	"github.com/slugger7/exorcist/apps/server/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/apps/server/internal/db/exorcist/public/table"
)

func CreateLibrary(ctx context.Context, db *sql.DB, lib model.Library) model.Library {
	library := table.Library
	stmnt := library.INSERT(library.Name, library.LibraryType).
		MODEL(lib).
		RETURNING(library.AllColumns)

	sql := stmnt.DebugSql()
	_ = sql

	var createdLib model.Library
	if err := stmnt.QueryContext(ctx, db, &createdLib); err != nil {
		log.Fatal(err)
	}

	return createdLib
}

func CreateStubLibrary(ctx context.Context, db *sql.DB) model.Library {
	return CreateLibrary(ctx, db, model.Library{Name: "stub", LibraryType: model.LibraryTypeEnum_Video})
}
