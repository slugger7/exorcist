package libraryRepository

import (
	"database/sql"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/environment"
	"github.com/slugger7/exorcist/internal/repository/util"
)

type ILibraryRepository interface {
	CreateLibraryStatement(name string) postgres.InsertStatement
}

type LibraryRepository struct {
	db  *sql.DB
	Env *environment.EnvironmentVariables
}

var libraryRepoInstance *LibraryRepository

func New(db *sql.DB, env *environment.EnvironmentVariables) ILibraryRepository {
	if libraryRepoInstance != nil {
		return libraryRepoInstance
	}
	libraryRepoInstance = &LibraryRepository{
		db:  db,
		Env: env,
	}
	return libraryRepoInstance
}

func (ls *LibraryRepository) CreateLibraryStatement(name string) postgres.InsertStatement {
	newLibrary := model.Library{
		Name: name,
	}

	insertStatement := table.Library.INSERT(table.Library.Name).
		MODEL(newLibrary).
		RETURNING(table.Library.ID)

	util.DebugCheck(ls.Env, insertStatement)

	return insertStatement
}
