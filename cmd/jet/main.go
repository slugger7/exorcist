package main

import (
	"database/sql"
	"fmt"
	"os"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	. "github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	. "github.com/slugger7/exorcist/internal/errors"
)

func main() {
	err := godotenv.Load()
	CheckError(err)

	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")

	fmt.Printf("host=%s port=%s user=%s password=%s database=%s", host, port, user, password, dbname)
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	fmt.Println("Opening DB")
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()

	err = db.Ping()
	CheckError(err)

	stmnt := SELECT(Library.AllColumns).FROM(Library)

	query := stmnt.DebugSql()

	fmt.Println(query)

	var dest []struct {
		model.Library
	}

	err = stmnt.Query(db, &dest)

	fmt.Println(dest[len(dest)-1].Name)
}
