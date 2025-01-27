package libraryPathRepository

import (
	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/repository"
)

func GetLibraryPathsSelect() postgres.SelectStatement {
	selectQuery := table.LibraryPath.
		SELECT(table.LibraryPath.ID, table.LibraryPath.Path).
		FROM(table.LibraryPath)

	repository.DebugCheck(selectQuery)
	return selectQuery
}

func CreateLibraryPath(libraryId uuid.UUID, path string) postgres.InsertStatement {
	newLibPath := model.LibraryPath{
		LibraryID: libraryId,
		Path:      path,
	}

	insertStatement := table.LibraryPath.
		INSERT(
			table.LibraryPath.LibraryID,
			table.LibraryPath.Path,
		).
		MODEL(newLibPath).
		RETURNING(table.LibraryPath.ID, table.LibraryPath.Path)

	repository.DebugCheck(insertStatement)

	return insertStatement
}
