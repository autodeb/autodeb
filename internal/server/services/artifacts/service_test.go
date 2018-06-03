package artifacts

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database/databasetest"
)

func setupTest(t *testing.T) *Service {
	db := databasetest.SetupTest(t)
	fs := filesystem.NewMemMapFS()
	service := New(db, fs)
	return service
}

func TestCreateArtifact(t *testing.T) {
	artifactsService := setupTest(t)

	artifacts, err := artifactsService.GetAllArtifactsByJobID(uint(1))
	assert.NoError(t, err)
	assert.Equal(t, 0, len(artifacts))

	artifactPath := "/1"
	_, err = artifactsService.fs.Stat(artifactPath)
	require.Error(t, err, "the job artifacts directory should not exist")

	_, err = artifactsService.CreateArtifact(
		uint(1),
		"artifact.txt",
		strings.NewReader("job artifact"),
	)

	_, err = artifactsService.fs.Stat(artifactPath)
	require.NoError(t, err, "the job artifact should exist")

	artifact, _ := artifactsService.GetArtifactContent(uint(1))
	defer artifact.Close()
	b, _ := ioutil.ReadAll(artifact)
	assert.Equal(t, "job artifact", string(b))

	artifacts, err = artifactsService.GetAllArtifactsByJobID(uint(1))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(artifacts))
}
