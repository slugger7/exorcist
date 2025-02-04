package userRepository

import (
	"strings"
	"testing"

	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
)

var s = UserRepository{
	Env: &environment.EnvironmentVariables{DebugSql: false},
}

func Test_GetUserByUsernameAndPassword(t *testing.T) {
	actual, _ := s.getUserByUsernameAndPasswordStatement("someUsername", "somePassword").Sql()

	exected := "\nSELECT \"user\".id AS \"user.id\",\n     \"user\".username AS \"user.username\"\nFROM public.\"user\"\nWHERE ((\"user\".username = $1::text) AND (\"user\".password = $2::text)) AND \"user\".active IS TRUE;\n"
	if exected != actual {
		t.Errorf("Expected %v but got %v", exected, actual)
	}
}

func Test_GetUserByUsername(t *testing.T) {
	actual, _ := s.getUserByUsernameStatement("someUsername").Sql()

	exected := `
SELECT "user".username AS "user.username"
FROM public."user"
WHERE ("user".username = $1::text) AND "user".active IS TRUE;
`
	if strings.Trim(exected, " ") != strings.Trim(actual, " ") {
		t.Errorf("Expected %s but got %s", exected, actual)
	}
}

func Test_CreateUser(t *testing.T) {
	user := model.User{
		Username: "someUsername",
		Password: "somePassword",
	}
	actual, _ := s.createUserStatement(user).Sql()

	exected := "\nINSERT INTO public.\"user\" (username, password)\nVALUES ($1, $2)\nRETURNING \"user\".id AS \"user.id\",\n          \"user\".username AS \"user.username\",\n          \"user\".active AS \"user.active\",\n          \"user\".created AS \"user.created\",\n          \"user\".modified AS \"user.modified\";\n"
	if exected != actual {
		t.Errorf("Expected %v but got %v", exected, actual)
	}
}
