package userRepository_test

import (
	"testing"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
	userRepository "github.com/slugger7/exorcist/internal/repository/user"
)

var s = userRepository.UserRepository{
	Env: &environment.EnvironmentVariables{DebugSql: false},
}

func Test_GetUserByUsernameAndPassword(t *testing.T) {
	actual, _ := s.GetUserByUsernameAndPassword("someUsername", "somePassword").Sql()

	exected := "\nSELECT \"user\".id AS \"user.id\",\n     \"user\".username AS \"user.username\"\nFROM public.\"user\"\nWHERE ((\"user\".username = $1::text) AND (\"user\".password = $2::text)) AND \"user\".active IS TRUE;\n"
	if exected != actual {
		t.Errorf("Expected %v but got %v", exected, actual)
	}
}

func Test_GetUserByUsername(t *testing.T) {
	actual, _ := s.GetUserByUsername("someUsername").Sql()

	exected := "\nSELECT \"user\".username AS \"user.username\"\nFROM public.\"user\"\nWHERE (\"user\".username = $1::text) AND \"user\".active IS TRUE;\n"
	if exected != actual {
		t.Errorf("Expected %v but got %v", exected, actual)
	}
}

func Test_CreateUser(t *testing.T) {
	user := model.User{
		Username: "someUsername",
		Password: "somePassword",
	}
	actual, _ := s.CreateUser(user).Sql()

	exected := "\nINSERT INTO public.\"user\" (username, password)\nVALUES ($1, $2)\nRETURNING \"user\".id AS \"user.id\",\n          \"user\".username AS \"user.username\",\n          \"user\".active AS \"user.active\",\n          \"user\".created AS \"user.created\",\n          \"user\".modified AS \"user.modified\";\n"
	if exected != actual {
		t.Errorf("Expected %v but got %v", exected, actual)
	}
}
