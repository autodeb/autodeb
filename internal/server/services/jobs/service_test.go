package jobs_test

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx/appctxtest"
)

func TestSaveJobLog(t *testing.T) {
	appCtxTest := appctxtest.SetupTest(t)

	jobDir := filepath.Join(
		appCtxTest.AppCtx.JobsService().JobsDirectory(),
		"1",
	)

	_, err := appCtxTest.AppCtx.JobsService().FS().Stat(jobDir)
	require.Error(t, err, "the job directory should not exist")

	err = appCtxTest.AppCtx.JobsService().SaveJobLog(
		uint(1),
		strings.NewReader("Hello"),
	)

	assert.NoError(t, err)

	_, err = appCtxTest.AppCtx.JobsService().FS().Stat(jobDir)
	assert.NoError(t, err)

	logFilePath := filepath.Join(jobDir, "log.txt")

	_, err = appCtxTest.AppCtx.JobsService().FS().Stat(logFilePath)
	assert.NoError(t, err)

	logFile, _ := appCtxTest.AppCtx.JobsService().FS().Open(logFilePath)
	defer logFile.Close()
	b, _ := ioutil.ReadAll(logFile)
	assert.Equal(t, "Hello", string(b))
}

func TestSaveJobArtifact(t *testing.T) {
	appCtxTest := appctxtest.SetupTest(t)
	jobsService := appCtxTest.AppCtx.JobsService()

	jobArtifactsDirectory := filepath.Join(
		jobsService.JobsDirectory(),
		"1",
		"artifacts",
	)
	_, err := jobsService.FS().Stat(jobArtifactsDirectory)
	require.Error(t, err, "the job artifacts directory should not exist")

	err = jobsService.SaveJobArtifact(
		uint(1),
		"artifact.txt",
		strings.NewReader("job artifact"),
	)

	_, err = jobsService.FS().Stat(jobArtifactsDirectory)
	require.NoError(t, err, "the job artifacts directory should exist")

	artifact, _ := jobsService.GetJobArtifact(uint(1), "artifact.txt")
	defer artifact.Close()
	b, _ := ioutil.ReadAll(artifact)
	assert.Equal(t, "job artifact", string(b))
}

func TestGetJobLog(t *testing.T) {
	appCtxTest := appctxtest.SetupTest(t)

	err := appCtxTest.AppCtx.JobsService().SaveJobLog(
		uint(1),
		strings.NewReader("Hello"),
	)
	assert.NoError(t, err)

	log, err := appCtxTest.AppCtx.JobsService().GetJobLog(uint(1))
	defer log.Close()

	assert.NoError(t, err)

	b, err := ioutil.ReadAll(log)

	assert.NoError(t, err)
	assert.Equal(t, "Hello", string(b))
}
