package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/table"
	er "github.com/slugger7/exorcist/internal/errors"
	ff "github.com/slugger7/exorcist/internal/ffmpeg"
	"github.com/slugger7/exorcist/internal/media"
)

func main() {
	err := godotenv.Load()
	er.CheckError(err)

	path := os.Getenv("MEDIA_PATH")

	db := setupDB()
	defer db.Close()

	libraryPath := getOrCreateLibraryPath(db, path)
	fmt.Printf("Library path id %v\n", libraryPath.ID)

	existingVideos := getVideosInLibraryPath(db, libraryPath.ID)

	fmt.Printf("Existing video count %v\n", len(existingVideos))

	values, err := media.GetFilesByExtensions(path, []string{".mp4", ".m4v", ".mkv", ".avi", ".wmv", ".flv", ".webm", ".f4v", ".mpg", ".m2ts", ".mov"})
	er.CheckError(err)

	nonExsistentVideos := media.FindNonExistentVideos(existingVideos, values)
	if len(nonExsistentVideos) > 0 {
		fmt.Println("Found some videos that do not exist any more on disk. Marking them as deleted.")
		removeVideos(db, nonExsistentVideos)
	}

	fmt.Println("Printing out results")
	videoModels := []model.Video{}
	for i, v := range values {
		printPercentage(i, len(values))
		relativePath := media.GetRelativePath(libraryPath.Path, v.Path)

		if videoExsists(existingVideos, relativePath) {
			continue
		}

		data, err := ff.UnmarshalledProbe(v.Path)
		if err != nil {
			fmt.Printf("Unmarshaling failed for %v\nThe error was %v", v.Path, err.Error())
			continue
		}

		width, height, err := ff.GetDimensions(data.Streams)
		if err != nil {
			fmt.Printf("Colud not extract dimensions. Setting to 0 %v\n", err.Error())
		}

		runtime, err := strconv.ParseFloat(data.Format.Duration, 32)
		if err != nil {
			fmt.Printf("Could not convert duration from string (%v) to float for video %v. Setting runtime to 0\n", data.Format.Duration, v)
			runtime = 0
		}
		size, err := strconv.Atoi(data.Format.Size)
		if err != nil {
			fmt.Printf("Could not convert size from string (%v) to int for video %v. Setting size to 0\n", data.Format.Size, v)
			size = 0
		}

		videoModels = append(videoModels, model.Video{
			LibraryPathID: libraryPath.ID,
			RelativePath:  relativePath,
			Title:         v.Name,
			FileName:      v.FileName,
			Height:        int32(height),
			Width:         int32(width),
			Runtime:       int64(runtime),
			Size:          int64(size),
			Checksum:      nil,
		})

		if i%5 == 0 {
			writeModelsToDatabaseBatch(db, videoModels)

			videoModels = []model.Video{}
		}
	}

	writeModelsToDatabaseBatch(db, videoModels)
}

func removeVideos(db *sql.DB, nonExistentVideos []model.Video) {
	for _, v := range nonExistentVideos {
		updateStmnt := table.Video.UPDATE().
			SET(table.Video.Deleted.SET(postgres.Bool(true))).
			MODEL(v).
			WHERE(table.Video.ID.EQ(postgres.UUID(v.ID)))
		dbgSql := updateStmnt.DebugSql()
		fmt.Println(dbgSql)
		_, err := updateStmnt.Exec(db)
		if err != nil {
			fmt.Printf("Could not update video %v to be deleted: %v", v.ID, err.Error())
			continue
		}
	}
}

func videoExsists(existingVideos []struct{ model.Video }, relativePath string) bool {
	return slices.ContainsFunc(existingVideos, func(existingVideo struct{ model.Video }) bool {
		return existingVideo.RelativePath == relativePath
	})
}

func getVideosInLibraryPath(db *sql.DB, libraryPathId uuid.UUID) []struct{ model.Video } {
	findStatement := table.Video.SELECT(table.Video.RelativePath, table.Video.ID).
		FROM(table.Video.Table).
		WHERE(table.Video.LibraryPathID.EQ(postgres.UUID(libraryPathId)))

	var videos []struct {
		model.Video
	}
	err := findStatement.Query(db, &videos)
	er.CheckError(err)

	return videos
}

func writeModelsToDatabaseBatch(db *sql.DB, models []model.Video) {
	if len(models) == 0 {
		return
	}
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
	fmt.Printf("Index: %v Total: %v Progress: %v\n", index, total, int(float64(index)/float64(total)*100.0))
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
