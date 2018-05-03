package app_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"salsa.debian.org/autodeb-team/autodeb/internal/errorchecks"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app/apptest"
)

func TestProcessFileUpload(t *testing.T) {
	testApp, fs, db := apptest.SetupTest(t)

	_, err := fs.Stat(filepath.Join(testApp.UploadedFilesDirectory(), "1"))
	require.Error(t, err, "the file directory should not exist")

	upload, err := testApp.ProcessUpload(
		&app.UploadParameters{
			Filename: "test.dsc",
		},
		strings.NewReader("this is a test file\n"),
	)
	assert.NoError(t, err)
	assert.Nil(t, upload)

	_, err = fs.Stat(filepath.Join(testApp.UploadedFilesDirectory(), "1"))
	assert.NoError(t, err)

	_, err = fs.Stat(filepath.Join(testApp.UploadedFilesDirectory(), "1", "test.dsc"))
	assert.NoError(t, err)

	expectedSHASum := "b6668cf8c46c7075e18215d922e7812ca082fa6cc34668d00a6c20aee4551fb6"

	fileUpload, err := db.GetPendingFileUpload(
		"test.dsc",
		expectedSHASum,
		false,
	)
	assert.NoError(t, err)
	assert.NotNil(t, fileUpload)

	assert.Equal(t, uint(1), fileUpload.ID)
	assert.Equal(t, "test.dsc", fileUpload.Filename)
	assert.Equal(t, expectedSHASum, fileUpload.SHA256Sum)
	assert.Equal(t, false, fileUpload.Completed)
}

const dummyChangesFile = `Format: 1.8
Date: Wed, 04 Apr 2018 14:28:29 -0400
Source: source-package-name
Binary: binary-package-name
Architecture: source
Version: 1.0+ds1-1
Distribution: unstable
Urgency: medium
Maintainer: Alexandre Viau <aviau@debian.org>
Changed-By: Alexandre Viau <aviau@debian.org>
Description:
 binary-package-name - dummy description
Changes:
 syncthing (1.0+ds1-1) unstable; urgency=medium
 .
   * Less bugs.
Checksums-Sha1:
 804d716fc5844f1cc5516c8f0be7a480517fdea2 20 test.dsc
Checksums-Sha256:
 b6668cf8c46c7075e18215d922e7812ca082fa6cc34668d00a6c20aee4551fb6 20 test.dsc
Files:
 66ad00916013ea0f7a6550f762b1de1d 20 utils optional test.dsc
`

func TestProcessChangesBadFormatRejected(t *testing.T) {
	testApp, _, _ := apptest.SetupTest(t)

	upload, err := testApp.ProcessUpload(
		&app.UploadParameters{
			Filename: "test.changes",
		},
		strings.NewReader("test"),
	)
	assert.Error(t, err)
	assert.True(t, errorchecks.IsInputError(err))
	assert.Nil(t, upload)
}

func TestProcessChangesMissingFile(t *testing.T) {
	testApp, _, _ := apptest.SetupTest(t)

	upload, err := testApp.ProcessUpload(
		&app.UploadParameters{
			Filename: "test.changes",
		},
		strings.NewReader(dummyChangesFile),
	)
	assert.Error(t, err)
	assert.True(t, errorchecks.IsInputError(err))
	assert.Nil(t, upload)
}

func TestProcessChanges(t *testing.T) {
	testApp, fs, _ := apptest.SetupTest(t)

	upload, err := testApp.ProcessUpload(
		&app.UploadParameters{
			Filename: "test.dsc",
		},
		strings.NewReader("this is a test file\n"),
	)
	assert.NoError(t, err)
	assert.Nil(t, upload)

	upload, err = testApp.ProcessUpload(
		&app.UploadParameters{
			Filename: "test.changes",
		},
		strings.NewReader(dummyChangesFile),
	)
	assert.NoError(t, err)
	assert.NotNil(t, upload)

	assert.Equal(t, uint(1), upload.ID)

	_, err = fs.Stat(filepath.Join(testApp.UploadedFilesDirectory(), "1"))
	assert.Error(t, err, "the uploaded files directory should be removed")

	_, err = fs.Stat(filepath.Join(testApp.UploadsDirectory(), "1"))
	assert.NoError(t, err)

	_, err = fs.Stat(filepath.Join(testApp.UploadsDirectory(), "1", "test.changes"))
	assert.NoError(t, err)

	_, err = fs.Stat(filepath.Join(testApp.UploadsDirectory(), "1", "test.dsc"))
	assert.NoError(t, err)
}
