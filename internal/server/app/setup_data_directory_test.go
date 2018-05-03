package app_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app/apptest"
)

func TestSetupDataDirectory(t *testing.T) {
	testApp, fs, _ := apptest.SetupTest(t)

	_, err := fs.Stat(testApp.UploadedFilesDirectory())
	assert.NoError(t, err, "the directory should have been created")

	_, err = fs.Stat(testApp.UploadsDirectory())
	assert.NoError(t, err, "the directory should have been created")
}
