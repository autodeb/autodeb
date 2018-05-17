package app_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app/apptest"
)

func TestSetupDataDirectory(t *testing.T) {
	appTest := apptest.SetupTest(t)

	_, err := appTest.DataFS.Stat(appTest.App.UploadedFilesDirectory())
	assert.NoError(t, err, "the directory should have been created")

	_, err = appTest.DataFS.Stat(appTest.App.UploadsDirectory())
	assert.NoError(t, err, "the directory should have been created")

	_, err = appTest.DataFS.Stat(appTest.App.JobsDirectory())
	assert.NoError(t, err, "the directory should have been created")
}
