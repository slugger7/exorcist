package mediaRepository

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/slugger7/exorcist/apps/server/internal/environment"
	"github.com/slugger7/exorcist/apps/server/internal/repository/repoTestHelpers"
	"github.com/slugger7/exorcist/apps/server/internal/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MediaRepoTestSuite struct {
	suite.Suite
	pgContainer *testhelpers.PostgresContainer
	repo        *mediaRepository
	ctx         context.Context
	db          *sql.DB
}

func (suite *MediaRepoTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	pgContainer, err := testhelpers.CreatePostgresContainer(suite.ctx)
	if err != nil {
		log.Fatal(err)
	}

	suite.pgContainer = pgContainer
	suite.db = pgContainer.SetupDatabase()

	env := &environment.EnvironmentVariables{}

	repo := New(suite.db, env, suite.ctx)

	suite.repo = repo
}

func (suite *MediaRepoTestSuite) TearDownSuite() {
	if err := suite.pgContainer.Terminate(suite.ctx); err != nil {
		log.Fatalf("error terminating postgres container: %s", err)
	}
}

func TestMediaRepoTestSuite(t *testing.T) {
	suite.Run(t, new(MediaRepoTestSuite))
}

func (suite *MediaRepoTestSuite) TestGetMediaById() {
	t := suite.T()

	stubMedia := repoTestHelpers.CreateStubMedia(suite.ctx, suite.db)

	media, err := suite.repo.GetById(stubMedia.ID)
	assert.NoError(t, err)
	assert.Equal(t, stubMedia.ID, media.Media.ID)
	assert.Equal(t, stubMedia.LibraryPathID, media.Media.LibraryPathID)
	assert.Equal(t, stubMedia.Path, media.Media.Path)
	assert.Equal(t, stubMedia.MediaType, media.Media.MediaType)
	assert.Equal(t, stubMedia.Title, media.Media.Title)
	assert.Equal(t, stubMedia.Size, media.Media.Size)
}
