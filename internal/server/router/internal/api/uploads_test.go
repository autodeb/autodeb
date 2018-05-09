package api_test

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/routertest"

	"github.com/stretchr/testify/assert"
)

func TestUploadDSCGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	uploadDir := filepath.Join(testRouter.App.UploadsDirectory(), "1")

	err := testRouter.DataFS.MkdirAll(uploadDir, 0644)
	assert.NoError(t, err)

	dsc, err := testRouter.DataFS.Create(filepath.Join(uploadDir, "test.dsc"))
	assert.NoError(t, err)
	dsc.Write([]byte("Hello"))
	dsc.Close()

	testRouter.Database.CreateFileUpload("test.dsc", "shasum", time.Now())

	fileUpload, err := testRouter.Database.GetFileUpload(uint(1))
	assert.NotNil(t, fileUpload)
	assert.NoError(t, err)

	fileUpload.UploadID = 1
	err = testRouter.Database.UpdateFileUpload(fileUpload)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/api/uploads/1/dsc", nil)

	testRouter.Router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, "text/plain", response.Result().Header.Get("Content-Type"))
	assert.Equal(t, "Hello", response.Body.String())
}

func TestUploadDSCGetHandlerNotFound(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/api/uploads/1/dsc", nil)

	testRouter.Router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
}

func TestUploadFileGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	uploadDir := filepath.Join(testRouter.App.UploadsDirectory(), "1")

	err := testRouter.DataFS.MkdirAll(uploadDir, 0644)
	assert.NoError(t, err)

	dsc, err := testRouter.DataFS.Create(filepath.Join(uploadDir, "test.dsc"))
	assert.NoError(t, err)
	dsc.Write([]byte("Hello"))
	dsc.Close()

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/api/uploads/1/test.dsc", nil)

	testRouter.Router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, "Hello", response.Body.String())
}
