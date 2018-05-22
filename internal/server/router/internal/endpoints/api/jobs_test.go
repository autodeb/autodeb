package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/routertest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobsNextPostHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	testRouter.DB.CreateJob(models.JobTypeBuild, uint(3))

	request, _ := http.NewRequest(http.MethodPost, "/api/jobs/next", nil)

	response := testRouter.ServeHTTP(request)
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

	request, _ := http.NewRequest(http.MethodPost, "/api/jobs/next", nil)

	response := testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusNoContent, response.Result().StatusCode)
	assert.Equal(t, "application/json", response.Result().Header.Get("Content-Type"))
	assert.Equal(t, "", response.Body.String())
}

func TestJobGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	_, err := testRouter.DB.CreateJob(models.JobTypeBuild, 1)
	assert.NoError(t, err)

	request, _ := http.NewRequest(http.MethodGet, "/api/jobs/1", nil)

	response := testRouter.ServeHTTP(request)
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

	request, _ := http.NewRequest(http.MethodGet, "/api/jobs/1", nil)

	response := testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
	assert.Equal(t, "application/json", response.Result().Header.Get("Content-Type"))
	assert.Equal(t, "", response.Body.String())
}

func TestJobStatusPostHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	job, err := testRouter.DB.CreateJob(models.JobTypeBuild, 1)
	assert.NoError(t, err)
	assert.NotEqual(t, models.JobStatusFailed, job.Status)

	job.Status = models.JobStatusAssigned
	testRouter.DB.UpdateJob(job)

	request, _ := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("/api/jobs/1/status/%d", models.JobStatusFailed),
		nil,
	)

	response := testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)

	job, err = testRouter.App.JobsService().GetJob(1)

	assert.NoError(t, err)
	assert.Equal(t, models.JobStatusFailed, job.Status)
}

func TestJobLogTxtGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	err := testRouter.App.JobsService().SaveJobLog(
		uint(1),
		strings.NewReader("hello"),
	)
	require.NoError(t, err)

	request, _ := http.NewRequest(
		http.MethodGet,
		"/api/jobs/1/log.txt",
		nil,
	)

	response := testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, "text/plain", response.Result().Header.Get("Content-Type"))
	assert.Equal(t, "hello", response.Body.String())
}

func TestJobLogTxtGetHandlerNotFound(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	request, _ := http.NewRequest(
		http.MethodGet,
		"/api/jobs/1/log.txt",
		nil,
	)

	response := testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
	assert.Equal(t, "text/plain", response.Result().Header.Get("Content-Type"))
	assert.Equal(t, "", response.Body.String())
}
