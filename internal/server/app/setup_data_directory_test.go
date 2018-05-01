package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupDataDirectory(t *testing.T) {
	app, fs, _ := setupTest(t)

	_, err := fs.Stat(app.UploadedFilesDirectory())
	assert.NoError(t, err, "the directory should have been created")

	_, err = fs.Stat(app.UploadsDirectory())
	assert.NoError(t, err, "the directory should have been created")
}
