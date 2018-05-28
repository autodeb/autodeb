package jobs

import (
	"io/ioutil"
	"path/filepath"
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

func TestSaveJobLog(t *testing.T) {
	jobsService := setupTest(t)

	jobDir := filepath.Join(
		jobsService.jobsDirectory(),
		"1",
	)

	_, err := jobsService.fs.Stat(jobDir)
	require.Error(t, err, "the job directory should not exist")

	err = jobsService.SaveJobLog(
		uint(1),
		strings.NewReader("Hello"),
	)

	assert.NoError(t, err)

	_, err = jobsService.fs.Stat(jobDir)
	assert.NoError(t, err)

	logFilePath := filepath.Join(jobDir, "log.txt")

	_, err = jobsService.fs.Stat(logFilePath)
	assert.NoError(t, err)

	logFile, _ := jobsService.fs.Open(logFilePath)
	defer logFile.Close()
	b, _ := ioutil.ReadAll(logFile)
	assert.Equal(t, "Hello", string(b))
}

func TestSaveJobArtifact(t *testing.T) {
	jobsService := setupTest(t)

	artifacts, err := jobsService.GetAllJobArtifactsByJobID(uint(1))
	assert.NoError(t, err)
	assert.Equal(t, 0, len(artifacts))

	jobArtifactsDirectory := filepath.Join(
		jobsService.jobsDirectory(),
		"1",
		"artifacts",
	)
	_, err = jobsService.fs.Stat(jobArtifactsDirectory)
	require.Error(t, err, "the job artifacts directory should not exist")

	err = jobsService.SaveJobArtifact(
		uint(1),
		"artifact.txt",
		strings.NewReader("job artifact"),
	)

	_, err = jobsService.fs.Stat(jobArtifactsDirectory)
	require.NoError(t, err, "the job artifacts directory should exist")

	_, err = jobsService.fs.Stat(filepath.Join(jobArtifactsDirectory, "artifact.txt"))
	require.NoError(t, err, "the job artifact should exist")

	artifact, _ := jobsService.GetJobArtifact(uint(1), "artifact.txt")
	defer artifact.Close()
	b, _ := ioutil.ReadAll(artifact)
	assert.Equal(t, "job artifact", string(b))

	artifacts, err = jobsService.GetAllJobArtifactsByJobID(uint(1))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(artifacts))
}

func TestGetJobLog(t *testing.T) {
	jobsService := setupTest(t)

	err := jobsService.SaveJobLog(
		uint(1),
		strings.NewReader("Hello"),
	)
	assert.NoError(t, err)

	log, err := jobsService.GetJobLog(uint(1))
	defer log.Close()

	assert.NoError(t, err)

	b, err := ioutil.ReadAll(log)

	assert.NoError(t, err)
	assert.Equal(t, "Hello", string(b))
}
