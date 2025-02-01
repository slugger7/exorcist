package userService

import (
	"errors"
	"log"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	"github.com/slugger7/exorcist/internal/environment"
	errs "github.com/slugger7/exorcist/internal/errors"
	"github.com/slugger7/exorcist/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	CreateUser(username, password string) (*model.User, error)
	ValidateUser(username, password string) (*model.User, error)
}

type UserService struct {
	Env  *environment.EnvironmentVariables
	repo repository.IRepository
}

var userServiceInstance *UserService

func New(repo repository.IRepository, env *environment.EnvironmentVariables) *UserService {
	if userServiceInstance == nil {
		userServiceInstance = &UserService{
			Env:  env,
			repo: repo,
		}

		log.Println("UserService instance created")
	}
	return userServiceInstance
}

func (us *UserService) UserExists(username string) (bool, error) {
	var users []struct {
		model.User
	}
	userRepo := us.repo.UserRepo()
	err := userRepo.GetUserByUsername(username).
		Query(&users)
	if err != nil {
		return false, err
	}

	return len(users) > 0, nil
}

func (us *UserService) CreateUser(username, password string) (*model.User, error) {
	log.Println("Creating user in service")
	userExists, err := us.UserExists(username)
	errs.CheckError(err)

	if userExists {
		return nil, errors.New("user already exists")
	}

	user := model.User{
		Username: username,
		Password: hashPassword(password), // needs to be hashed and salted
	}
	var users []struct {
		model.User
	}
	err = us.repo.UserRepo().CreateUser(user).Query(&users)
	if err != nil {
		return nil, errors.Join(errors.New("could not create a new user"), err)
	}

	return &users[len(users)-1].User, nil
}

func (us *UserService) ValidateUser(username, password string) (*model.User, error) {
	var users []struct {
		model.User
	}
	err := us.repo.UserRepo().
		GetUserByUsername(username, table.User.ID, table.User.Password).
		Query(&users)
	if err != nil {
		return nil, err
	}
	if len(users) > 1 {
		panic("Found more than one active user for a username")
	}

	if len(users) != 1 {
		return nil, nil
	}

	user := users[len(users)-1].User

	if !compareHashedPassword(user.Password, password) {
		log.Printf("Password did not match hashed password in database for user %v", username)
		return nil, nil
	}
	user.Password = "" // do not want to return the hash of the password

	return &user, nil
}

func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	errs.CheckError(err)

	return string(hashedPassword)
}

func compareHashedPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
