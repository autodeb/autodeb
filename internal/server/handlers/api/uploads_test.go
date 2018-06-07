package api_test

import (
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/routertest"

	"github.com/stretchr/testify/assert"
)

func TestUploadDSCGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	apiClient := testRouter.APIClient
	uploadsService := testRouter.AppCtx.UploadsService()

	uploadDir := filepath.Join(uploadsService.UploadsDirectory(), "1")

	err := uploadsService.FS().MkdirAll(uploadDir, 0644)
	assert.NoError(t, err)

	dsc, err := uploadsService.FS().Create(filepath.Join(uploadDir, "test.dsc"))
	assert.NoError(t, err)
	dsc.Write([]byte("Hello"))
	dsc.Close()

	testRouter.DB.CreateFileUpload("test.dsc", "shasum", time.Now())

	fileUpload, err := testRouter.DB.GetFileUpload(uint(1))
	assert.NotNil(t, fileUpload)
	assert.NoError(t, err)

	fileUpload.UploadID = 1
	err = testRouter.DB.UpdateFileUpload(fileUpload)
	assert.NoError(t, err)

	_, err = apiClient.GetUploadDSC(uint(1))
	assert.NoError(t, err)

	response := apiClient.LastRecorder()
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, "text/plain", response.Result().Header.Get("Content-Type"))
	assert.Equal(t, "Hello", response.Body.String())
}

func TestUploadGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	apiClient := testRouter.APIClient

	upload, err := testRouter.DB.CreateUpload(1, "testSource", "testVersion", "testMaintainer", "testChangedBy", true, true)
	assert.NoError(t, err)

	returnedUpload, err := apiClient.GetUpload(upload.ID)
	assert.NoError(t, err)
	assert.Equal(t, upload, returnedUpload)

	response := apiClient.LastRecorder()
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, "application/json", response.Result().Header.Get("Content-Type"))
}

func TestUploadChangesGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	apiClient := testRouter.APIClient
	uploadsService := testRouter.AppCtx.UploadsService()

	uploadDir := filepath.Join(uploadsService.UploadsDirectory(), "1")

	err := uploadsService.FS().MkdirAll(uploadDir, 0644)
	assert.NoError(t, err)

	dsc, err := uploadsService.FS().Create(filepath.Join(uploadDir, "test.changes"))
	assert.NoError(t, err)
	dsc.Write([]byte("Hello"))
	dsc.Close()

	testRouter.DB.CreateFileUpload("test.changes", "shasum", time.Now())

	fileUpload, err := testRouter.DB.GetFileUpload(uint(1))
	assert.NotNil(t, fileUpload)
	assert.NoError(t, err)

	fileUpload.UploadID = 1
	err = testRouter.DB.UpdateFileUpload(fileUpload)
	assert.NoError(t, err)

	_, err = apiClient.GetUploadChanges(uint(1))
	assert.NoError(t, err)

	response := apiClient.LastRecorder()
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, "text/plain", response.Result().Header.Get("Content-Type"))
	assert.Equal(t, "Hello", response.Body.String())
}

func TestUploadDSCGetHandlerNotFound(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	apiClient := testRouter.APIClient

	dsc, err := apiClient.GetUploadDSC(uint(1))
	assert.NoError(t, err)
	assert.Nil(t, dsc)

	response := apiClient.LastRecorder()
	assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
}

func TestUploadFileGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	apiClient := testRouter.APIClient
	uploadsService := testRouter.AppCtx.UploadsService()

	uploadDir := filepath.Join(uploadsService.UploadsDirectory(), "1")

	err := uploadsService.FS().MkdirAll(uploadDir, 0644)
	assert.NoError(t, err)

	dsc, err := uploadsService.FS().Create(filepath.Join(uploadDir, "test.dsc"))
	assert.NoError(t, err)
	dsc.Write([]byte("Hello"))
	dsc.Close()

	uploadFile, err := apiClient.GetUploadFile(uint(1), "test.dsc")
	assert.NoError(t, err)
	assert.NotNil(t, uploadFile)

	response := apiClient.LastRecorder()
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, "Hello", response.Body.String())
}

func TestUploadFilesGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	apiClient := testRouter.APIClient

	testRouter.DB.CreateFileUpload("test", "sum", time.Now())
	fileUpload, _ := testRouter.DB.GetFileUpload(uint(1))
	fileUpload.UploadID = uint(3)
	fileUpload.Completed = true

	err := testRouter.DB.UpdateFileUpload(fileUpload)
	assert.NoError(t, err)

	fileUploads, err := apiClient.GetUploadFiles(uint(3))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(fileUploads))
	assert.Equal(t, "test", fileUploads[0].Filename)

	response := apiClient.LastRecorder()
	assert.Equal(t, "application/json", response.Result().Header.Get("Content-Type"))
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}
