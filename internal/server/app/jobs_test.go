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
	testApp, fs, _ := apptest.SetupTest(t)

	jobDir := filepath.Join(
		testApp.JobsDirectory(),
		"1",
	)

	_, err := fs.Stat(jobDir)
	require.Error(t, err, "the job directory should not exist")

	err = testApp.SaveJobLog(
		uint(1),
		strings.NewReader("Hello"),
	)

	assert.NoError(t, err)

	_, err = fs.Stat(jobDir)
	assert.NoError(t, err)

	logFilePath := filepath.Join(jobDir, "log.txt")

	_, err = fs.Stat(logFilePath)
	assert.NoError(t, err)

	logFile, _ := fs.Open(logFilePath)
	defer logFile.Close()
	b, _ := ioutil.ReadAll(logFile)
	assert.Equal(t, "Hello", string(b))
}

func TestGetJobLog(t *testing.T) {
	testApp, _, _ := apptest.SetupTest(t)

	err := testApp.SaveJobLog(
		uint(1),
		strings.NewReader("Hello"),
	)
	assert.NoError(t, err)

	log, err := testApp.GetJobLog(uint(1))
	defer log.Close()

	assert.NoError(t, err)

	b, err := ioutil.ReadAll(log)

	assert.NoError(t, err)
	assert.Equal(t, "Hello", string(b))
}
