package app_test

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/app/apptest"
)

func TestSaveJobLog(t *testing.T) {
	appTest := apptest.SetupTest(t)

	jobDir := filepath.Join(
		appTest.App.JobsDirectory(),
		"1",
	)

	_, err := appTest.DataFS.Stat(jobDir)
	require.Error(t, err, "the job directory should not exist")

	err = appTest.App.SaveJobLog(
		uint(1),
		strings.NewReader("Hello"),
	)

	assert.NoError(t, err)

	_, err = appTest.DataFS.Stat(jobDir)
	assert.NoError(t, err)

	logFilePath := filepath.Join(jobDir, "log.txt")

	_, err = appTest.DataFS.Stat(logFilePath)
	assert.NoError(t, err)

	logFile, _ := appTest.DataFS.Open(logFilePath)
	defer logFile.Close()
	b, _ := ioutil.ReadAll(logFile)
	assert.Equal(t, "Hello", string(b))
}

func TestGetJobLog(t *testing.T) {
	appTest := apptest.SetupTest(t)

	err := appTest.App.SaveJobLog(
		uint(1),
		strings.NewReader("Hello"),
	)
	assert.NoError(t, err)

	log, err := appTest.App.GetJobLog(uint(1))
	defer log.Close()

	assert.NoError(t, err)

	b, err := ioutil.ReadAll(log)

	assert.NoError(t, err)
	assert.Equal(t, "Hello", string(b))
}
