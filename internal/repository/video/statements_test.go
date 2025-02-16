package videoRepository

import (
	"testing"

	"github.com/google/uuid"
	"github.com/slugger7/exorcist/internal/db/exorcist/public/model"
	"github.com/slugger7/exorcist/internal/environment"
)

// TODO: implement [snapshot tests](https://github.com/gkampitakis/go-snaps)

var ds = &VideoRepository{
	Env: &environment.EnvironmentVariables{DebugSql: false},
}

func Test_GetVideoWithoutChecksumStatement(t *testing.T) {
	actual, _ := ds.GetVideoWithoutChecksumStatement().Sql()

	expected := "\nSELECT video.id AS \"video.id\",\n     video.checksum AS \"video.checksum\",\n     video.relative_path AS \"video.relative_path\",\n     library_path.path AS \"library_path.path\"\nFROM public.video\n     INNER JOIN public.library_path ON (library_path.id = video.library_path_id)\nWHERE video.checksum IS NULL;\n"
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func Test_UpdateVideoChecksum(t *testing.T) {
	newUuid, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("Encountered an error while generating a UUID: %v", err)
	}

	checksum := "someChecksum"

	video := model.Video{
		ID:       newUuid,
		Checksum: &checksum,
	}

	actual, _ := ds.UpdateVideoChecksum(video).Sql()

	expected := `
UPDATE public.video
SET checksum = $1::text
WHERE video.id = $2;
`
	if actual != expected {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func Test_MarkVideoAsNotExistingStatement(t *testing.T) {
	newUuid, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("Encountered an error while generating a UUID: %v", err)
	}

	video := model.Video{
		ID:     newUuid,
		Exists: false,
	}

	actual, _ := ds.updateVideoExistsStatement(video).Sql()

	expected := `
UPDATE public.video
SET exists = $1::boolean
WHERE video.id = $2;
`
	if actual != expected {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func Test_GetVideosInLibraryPath(t *testing.T) {
	newUuid, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("Encountered an error while generating a UUID: %v", err)
	}
	actual, _ := ds.getByLibraryPathIdStatement(newUuid).Sql()

	expected := `
SELECT video.relative_path AS "video.relative_path",
     video.id AS "video.id"
FROM public.video
WHERE (video.library_path_id = $1) AND video.exists IS TRUE;
`
	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func Test_InsertVideosStatement_WithNoVideos_ShouldReturnNil(t *testing.T) {
	videos := []model.Video{}
	actual := ds.insertStatement(videos)

	if actual != nil {
		t.Errorf("Expected actual to be nil. Acutal: %v", actual)
	}
}

func Test_InsertVideosStatement_WithVideos_ShouldReturnStatement(t *testing.T) {
	newUuid, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("Encountered an error while generating a UUID: %v", err)
	}
	video := model.Video{
		LibraryPathID: newUuid,
		RelativePath:  "relativePath",
		Title:         "title",
		FileName:      "filename",
		Height:        69,
		Width:         420,
		Runtime:       1337,
		Size:          80085,
	}
	videos := []model.Video{video}
	actual, _ := ds.insertStatement(videos).Sql()

	expected :=
		`
INSERT INTO public.video (library_path_id, relative_path, title, file_name, height, width, runtime, size)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
`
	if actual != expected {
		t.Errorf("Expected \n%v got \n%v", expected, actual)
	}
}

func Test_GetByIdStatement(t *testing.T) {
	id, _ := uuid.NewRandom()

	sql, _ := ds.getByIdStatement(id).Sql()

	expectedSql := "\nSELECT video.id AS \"video.id\",\n     video.library_path_id AS \"video.library_path_id\",\n     video.relative_path AS \"video.relative_path\",\n     video.title AS \"video.title\",\n     video.file_name AS \"video.file_name\",\n     video.height AS \"video.height\",\n     video.width AS \"video.width\",\n     video.runtime AS \"video.runtime\",\n     video.size AS \"video.size\",\n     video.checksum AS \"video.checksum\",\n     video.added AS \"video.added\",\n     video.deleted AS \"video.deleted\",\n     video.exists AS \"video.exists\",\n     video.created AS \"video.created\",\n     video.modified AS \"video.modified\"\nFROM public.video\nWHERE video.id = $1\nLIMIT $2;\n"
	if expectedSql != sql {
		t.Errorf("Expected sql: %v\nGot sql: %v", expectedSql, sql)
	}
}
