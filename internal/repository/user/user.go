package userRepository

import (
	"database/sql"
	"log"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
)

type UserStatement struct {
	postgres.Statement
	db *sql.DB
}

type IUserRepository interface {
	GetUserByUsernameAndPassword(username, password string) (*model.User, error)
	GetUserByUsername(username string, columns ...postgres.Projection) (*model.User, error)
	CreateUser(user model.User) (*model.User, error)
}

type UserRepository struct {
	db  *sql.DB
	Env *environment.EnvironmentVariables
}

var userRepositoryInstance *UserRepository

func New(db *sql.DB, env *environment.EnvironmentVariables) IUserRepository {
	if userRepositoryInstance == nil {
		userRepositoryInstance = &UserRepository{
			db:  db,
			Env: env,
		}

		log.Println("User repository instance created")
	}

	return userRepositoryInstance
}

func (us *UserStatement) Query(destination interface{}) error {
	return us.Statement.Query(us.db, destination)
}

func (ur *UserRepository) CreateUser(user model.User) (*model.User, error) {
	var newUser struct{ model.User }
	if err := ur.createUserStatement(user).Query(&newUser); err != nil {
		return nil, errs.BuildError(err, "could not create user %v", user)
	}
	return &newUser.User, nil
}

func (ur *UserRepository) GetUserByUsername(username string, columns ...postgres.Projection) (*model.User, error) {
	var users []struct{ model.User }
	if err := ur.getUserByUsernameStatement(username, columns...).Query(&users); err != nil {
		log.Println("something went wrong getting user by username")
		return nil, errs.BuildError(err, "could not get user by username '%v'", username)
	}
	var user *model.User
	if len(users) > 0 {
		user = &users[len(users)-1].User
	}
	return user, nil
}

func (ur *UserRepository) GetUserByUsernameAndPassword(username, password string) (*model.User, error) {
	var users []struct{ model.User }
	if err := ur.getUserByUsernameAndPasswordStatement(username, password).Query(&users); err != nil {
		log.Println("something went wrong getting user by username")
		return nil, errs.BuildError(err, "could not get user by username '%v' and password", username)
	}
	var user *model.User
	if len(users) > 0 {
		user = &users[len(users)-1].User
	}
	return user, nil
}
