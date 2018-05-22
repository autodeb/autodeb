package api_test

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/routertest"

	"github.com/stretchr/testify/assert"
)

func TestUploadDSCGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)
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

	request, _ := http.NewRequest(http.MethodGet, "/api/uploads/1/dsc", nil)

	response := testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, "text/plain", response.Result().Header.Get("Content-Type"))
	assert.Equal(t, "Hello", response.Body.String())
}

func TestUploadDSCGetHandlerNotFound(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	request, _ := http.NewRequest(http.MethodGet, "/api/uploads/1/dsc", nil)

	response := testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
}

func TestUploadFileGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	uploadsService := testRouter.AppCtx.UploadsService()

	uploadDir := filepath.Join(uploadsService.UploadsDirectory(), "1")

	err := uploadsService.FS().MkdirAll(uploadDir, 0644)
	assert.NoError(t, err)

	dsc, err := uploadsService.FS().Create(filepath.Join(uploadDir, "test.dsc"))
	assert.NoError(t, err)
	dsc.Write([]byte("Hello"))
	dsc.Close()

	request, _ := http.NewRequest(http.MethodGet, "/api/uploads/1/test.dsc", nil)

	response := testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, "Hello", response.Body.String())
}

func TestUploadFilesGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	testRouter.DB.CreateFileUpload("test", "sum", time.Now())
	fileUpload, _ := testRouter.DB.GetFileUpload(uint(1))
	fileUpload.UploadID = uint(3)
	fileUpload.Completed = true
	err := testRouter.DB.UpdateFileUpload(fileUpload)

	assert.NoError(t, err)

	request, _ := http.NewRequest(http.MethodGet, "/api/uploads/3/files", nil)

	response := testRouter.ServeHTTP(request)
	assert.Equal(t, "application/json", response.Result().Header.Get("Content-Type"))
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)

	var fileUploads []models.FileUpload
	err = json.Unmarshal(response.Body.Bytes(), &fileUploads)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(fileUploads))
	assert.Equal(t, "test", fileUploads[0].Filename)
}
