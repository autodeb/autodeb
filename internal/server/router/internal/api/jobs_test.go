package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/routertest"

	"github.com/stretchr/testify/assert"
)

func TestJobsNextPostHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	testRouter.Database.CreateJob(models.JobTypeBuild, uint(3))

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/api/jobs/next", nil)

	testRouter.Router.ServeHTTP(response, request)

	assert.Equal(t, "application/json", response.Result().Header.Get("Content-Type"))
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)

	var job models.Job
	err := json.Unmarshal(response.Body.Bytes(), &job)

	expected := models.Job{
		ID:       uint(1),
		Type:     models.JobTypeBuild,
		Status:   models.JobStatusAssigned,
		UploadID: uint(3),
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, job)
}

func TestJobsNextPostHandlerNoJob(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/api/jobs/next", nil)

	testRouter.Router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNoContent, response.Result().StatusCode)
	assert.Equal(t, "application/json", response.Result().Header.Get("Content-Type"))
	assert.Equal(t, response.Body.String(), "")
}

func TestJobGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	_, err := testRouter.Database.CreateJob(models.JobTypeBuild, 1)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/api/jobs/1", nil)

	testRouter.Router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, "application/json", response.Result().Header.Get("Content-Type"))

	var job models.Job
	err = json.Unmarshal(response.Body.Bytes(), &job)
	assert.NoError(t, err)

	assert.Equal(t, uint(1), job.ID)
	assert.Equal(t, models.JobTypeBuild, job.Type)
	assert.Equal(t, uint(1), job.UploadID)
}

func TestJobGetHandlerNotFound(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/api/jobs/1", nil)

	testRouter.Router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
	assert.Equal(t, "application/json", response.Result().Header.Get("Content-Type"))
	assert.Equal(t, response.Body.String(), "")
}
