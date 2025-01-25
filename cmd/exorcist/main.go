package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	er "github.com/slugger7/exorcist/internal/errors"
	ff "github.com/slugger7/exorcist/internal/ffmpeg"
	"github.com/slugger7/exorcist/internal/media"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func main() {
	err := godotenv.Load()
	er.CheckError(err)

	path := os.Getenv("MEDIA_PATH")

	db := setupDB()
	defer db.Close()

	libraryPath := getOrCreateLibraryPath(db, path)
	fmt.Printf("Library path id %v\n", libraryPath.ID)

	values, err := media.GetFilesByExtensions(path, []string{".mp4", ".m4v", ".mkv", ".avi", ".wmv", ".flv", ".webm", ".f4v", ".mpg", ".m2ts", ".mov"})
	er.CheckError(err)

	fmt.Println("Printing out results")
	videoModels := []model.Video{}
	for i, v := range values {
		printPercentage(i, len(values))
		checksum := "lol"

		probeData, err := ffmpeg.Probe(v.Path)
		if err != nil {
			fmt.Printf("Could not probe the following file %v.\nThis is the error: %v", v.Path, err.Error())
			continue
		}

		var data *ff.Probe
		err = json.Unmarshal([]byte(probeData), &data)
		er.CheckError(err)

		width, height, err := ff.GetDimensions(data.Streams)
		if err != nil {
			fmt.Printf("Colud not extract dimensions. Setting to 0 %v\n", err.Error())
		}

		runtime, err := strconv.ParseFloat(data.Format.Duration, 32)
		if err != nil {
			fmt.Printf("Could not convert duration from string (%v) to int for video %v. Setting runtime to 0\n", data.Format.Duration, v)
			runtime = 0
		}
		size, err := strconv.Atoi(data.Format.Size)
		if err != nil {
			fmt.Printf("Could not convert size from string (%v) to int for video %v. Setting size to 0\n", data.Format.Size, v)
			size = 0
		}

		videoModels = append(videoModels, model.Video{
			LibraryPathID: libraryPath.ID,
			RelativePath:  media.GetRelativePath(libraryPath.Path, v.Path),
			Title:         v.Name,
			FileName:      v.FileName,
			Height:        int32(height),
			Width:         int32(width),
			Runtime:       int64(runtime),
			Size:          int64(size),
			Checksum:      &checksum,
		})

		if i%100 == 0 {
			writeModelsToDatabaseBatch(db, videoModels)

			videoModels = []model.Video{}
		}
	}

	writeModelsToDatabaseBatch(db, videoModels)
}

func writeModelsToDatabaseBatch(db *sql.DB, models []model.Video) {
	fmt.Println("Writing batch")

	insertStatement := table.Video.INSERT(
		table.Video.LibraryPathID,
		table.Video.RelativePath,
		table.Video.Title,
		table.Video.FileName,
		table.Video.Height,
		table.Video.Width,
		table.Video.Runtime,
		table.Video.Size,
		table.Video.Checksum,
	).
		MODELS(models).
		RETURNING(table.Video.ID)

	var newVideos []struct {
		model.Video
	}
	err := insertStatement.Query(db, &newVideos)
	er.CheckError(err)
}

func printPercentage(index, total int) {
	fmt.Printf("Index: %v Total: %v Progress: %v%%\n", index, total, index/total*100)
}

func getOrCreateLibraryPath(db *sql.DB, path string) model.LibraryPath {
	libraryPath, err := getExistingLibraryPathID(db)
	if err != nil {
		libraryPath = createLibWithPath(db, path)
	}
	return libraryPath
}

func getExistingLibraryPathID(db *sql.DB) (model.LibraryPath, error) {
	selectQuery := table.LibraryPath.
		SELECT(table.LibraryPath.ID, table.LibraryPath.Path).
		FROM(table.LibraryPath)

	var libraryPath []struct {
		model.LibraryPath
	}

	err := selectQuery.Query(db, &libraryPath)
	er.CheckError(err)

	if len(libraryPath) == 0 {
		return model.LibraryPath{}, errors.New("no library path was found, first creat a library")
	}

	return libraryPath[0].LibraryPath, nil
}

func createLibWithPath(db *sql.DB, path string) model.LibraryPath {
	newLib := model.Library{
		Name: "New Lib",
	}

	insertStatement := table.Library.INSERT(table.Library.Name).
		MODEL(newLib).
		RETURNING(table.Library.ID)

	var library []struct {
		model.Library
	}

	err := insertStatement.Query(db, &library)
	er.CheckError(err)

	newLibPath := model.LibraryPath{
		LibraryID: library[0].ID,
		Path:      path,
	}

	insertStatement = table.LibraryPath.INSERT(
		table.LibraryPath.LibraryID,
		table.LibraryPath.Path,
	).
		MODEL(newLibPath).
		RETURNING(table.LibraryPath.ID, table.LibraryPath.Path)

	var libraryPath []struct {
		model.LibraryPath
	}

	err = insertStatement.Query(db, &libraryPath)
	er.CheckError(err)

	return libraryPath[0].LibraryPath
}

func setupDB() *sql.DB {
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	user := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")

	fmt.Printf("host=%s port=%s user=%s password=%s database=%s", host, port, user, password, dbname)
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	fmt.Println("Opening DB")
	db, err := sql.Open("postgres", psqlconn)
	er.CheckError(err)

	return db
}
