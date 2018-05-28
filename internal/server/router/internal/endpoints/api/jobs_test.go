package api_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	job, err = testRouter.AppCtx.JobsService().GetJob(1)

	assert.NoError(t, err)
	assert.Equal(t, models.JobStatusFailed, job.Status)
}

func TestJobLogPostHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	job, err := testRouter.DB.CreateJob(models.JobTypeBuild, 1)
	assert.NoError(t, err)

	job.Status = models.JobStatusSuccess
	err = testRouter.DB.UpdateJob(job)
	assert.NoError(t, err)

	request, _ := http.NewRequest(
		http.MethodPost,
		"/api/jobs/1/log",
		strings.NewReader("log content test"),
	)

	response := testRouter.ServeHTTP(request)

	assert.Equal(t, http.StatusCreated, response.Result().StatusCode)
	assert.Equal(t, "", response.Body.String())

	log, err := testRouter.AppCtx.JobsService().GetJobLog(uint(1))
	assert.NoError(t, err)
	defer log.Close()

	b, err := ioutil.ReadAll(log)
	assert.Equal(t, "log content test", string(b))
}

func TestJobArtifactPostHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	job, err := testRouter.DB.CreateJob(models.JobTypeBuild, 1)
	assert.NoError(t, err)

	job.Status = models.JobStatusSuccess
	err = testRouter.DB.UpdateJob(job)
	assert.NoError(t, err)

	request, _ := http.NewRequest(
		http.MethodPost,
		"/api/jobs/1/artifacts/test.txt",
		strings.NewReader("test txt content"),
	)

	response := testRouter.ServeHTTP(request)

	assert.Equal(t, http.StatusCreated, response.Result().StatusCode)
	assert.Equal(t, "", response.Body.String())

	log, err := testRouter.AppCtx.JobsService().GetJobArtifact(uint(1), "test.txt")
	assert.NoError(t, err)
	defer log.Close()

	b, err := ioutil.ReadAll(log)
	assert.Equal(t, "test txt content", string(b))
}

func TestJobLogTxtGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	err := testRouter.AppCtx.JobsService().SaveJobLog(
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

func TestJobArtifactGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	err := testRouter.AppCtx.JobsService().SaveJobArtifact(
		uint(1),
		"test.txt",
		strings.NewReader("test content"),
	)
	require.NoError(t, err)

	request, _ := http.NewRequest(
		http.MethodGet,
		"/api/jobs/1/artifacts/test.txt",
		nil,
	)

	response := testRouter.ServeHTTP(request)

	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, "text/plain", response.Result().Header.Get("Content-Type"))
	assert.Equal(t, "test content", response.Body.String())
}

func TestJobsArtifactsGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	err := testRouter.AppCtx.JobsService().SaveJobArtifact(
		uint(1),
		"test.txt",
		strings.NewReader("test content"),
	)
	require.NoError(t, err)

	request, _ := http.NewRequest(
		http.MethodGet,
		"/api/jobs/1/artifacts",
		nil,
	)

	response := testRouter.ServeHTTP(request)
	assert.Equal(t, "application/json", response.Result().Header.Get("Content-Type"))
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)

	var jobArtifacts []*models.JobArtifact
	err = json.Unmarshal(response.Body.Bytes(), &jobArtifacts)
	assert.NoError(t, err)

	expected := &models.JobArtifact{
		ID:       uint(1),
		JobID:    uint(1),
		Filename: "test.txt",
	}
	assert.NoError(t, err)
	assert.Equal(t, expected, jobArtifacts[0])
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
