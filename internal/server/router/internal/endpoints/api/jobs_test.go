package api_test

import (
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
	apiClient := testRouter.APIClient
	testRouter.Login()

	testRouter.Services.Jobs().CreateBuildJob(uint(3))

	job, err := apiClient.UnqueueNextJob()
	assert.NoError(t, err)

	expected := &models.Job{
		ID:       uint(1),
		Type:     models.JobTypeBuild,
		Status:   models.JobStatusAssigned,
		UploadID: uint(3),
	}
	assert.Equal(t, expected, job)

	response := apiClient.LastResponse()
	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestJobsNextPostHandlerNoJob(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	apiClient := testRouter.APIClient
	testRouter.Login()

	job, err := apiClient.UnqueueNextJob()
	assert.NoError(t, err)
	assert.Nil(t, job)

	response := apiClient.LastRecorder()
	assert.Equal(t, http.StatusNoContent, response.Result().StatusCode)
	assert.Equal(t, "application/json", response.Result().Header.Get("Content-Type"))
	assert.Equal(t, "", response.Body.String())
}

func TestJobGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	apiClient := testRouter.APIClient

	_, err := testRouter.Services.Jobs().CreateBuildJob(1)
	assert.NoError(t, err)

	job, err := apiClient.GetJob(uint(1))
	assert.NoError(t, err)
	assert.Equal(t, uint(1), job.ID)
	assert.Equal(t, models.JobTypeBuild, job.Type)
	assert.Equal(t, uint(1), job.UploadID)

	response := apiClient.LastRecorder()
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, "application/json", response.Result().Header.Get("Content-Type"))
}

func TestJobGetHandlerNotFound(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	apiClient := testRouter.APIClient

	job, err := apiClient.GetJob(uint(1))
	assert.Nil(t, job)
	assert.NoError(t, err)

	response := apiClient.LastRecorder()
	assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
	assert.Equal(t, "application/json", response.Result().Header.Get("Content-Type"))
	assert.Equal(t, "", response.Body.String())
}

func TestJobStatusPostHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	apiClient := testRouter.APIClient
	testRouter.Login()

	job, err := testRouter.Services.Jobs().CreateBuildJob(1)
	assert.NoError(t, err)
	assert.NotEqual(t, models.JobStatusFailed, job.Status)

	job.Status = models.JobStatusAssigned
	testRouter.DB.UpdateJob(job)

	err = apiClient.SetJobStatus(uint(1), models.JobStatusFailed)
	assert.NoError(t, err)

	response := apiClient.LastResponse()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	job, err = testRouter.AppCtx.JobsService().GetJob(1)

	assert.NoError(t, err)
	assert.Equal(t, models.JobStatusFailed, job.Status)
}

func TestJobLogPostHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	apiClient := testRouter.APIClient
	testRouter.Login()

	job, err := testRouter.Services.Jobs().CreateBuildJob(1)
	assert.NoError(t, err)

	job.Status = models.JobStatusAssigned
	err = testRouter.DB.UpdateJob(job)
	assert.NoError(t, err)

	err = apiClient.SubmitJobLog(
		uint(1),
		strings.NewReader("log content test"),
	)
	assert.NoError(t, err)

	response := apiClient.LastRecorder()
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
	apiClient := testRouter.APIClient
	testRouter.Login()

	job, err := testRouter.Services.Jobs().CreateBuildJob(1)
	assert.NoError(t, err)

	job.Status = models.JobStatusAssigned
	err = testRouter.DB.UpdateJob(job)
	assert.NoError(t, err)

	err = apiClient.SubmitJobArtifact(
		uint(1),
		"test.txt",
		strings.NewReader("test txt content"),
	)
	assert.NoError(t, err)

	response := apiClient.LastRecorder()
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
	apiClient := testRouter.APIClient

	err := testRouter.AppCtx.JobsService().SaveJobLog(
		uint(1),
		strings.NewReader("hello"),
	)
	require.NoError(t, err)

	_, err = apiClient.GetJobLog(uint(1))
	assert.NoError(t, err)

	response := apiClient.LastRecorder()
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, "text/plain", response.Result().Header.Get("Content-Type"))
	assert.Equal(t, "hello", response.Body.String())
}

func TestJobArtifactGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	apiClient := testRouter.APIClient

	err := testRouter.AppCtx.JobsService().SaveJobArtifact(
		uint(1),
		"test.txt",
		strings.NewReader("test content"),
	)
	require.NoError(t, err)

	_, err = apiClient.GetJobArtifact(uint(1), "test.txt")
	assert.NoError(t, err)

	response := apiClient.LastRecorder()
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, "text/plain", response.Result().Header.Get("Content-Type"))
	assert.Equal(t, "test content", response.Body.String())
}

func TestJobsArtifactsGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	apiClient := testRouter.APIClient

	err := testRouter.AppCtx.JobsService().SaveJobArtifact(
		uint(1),
		"test.txt",
		strings.NewReader("test content"),
	)
	require.NoError(t, err)

	jobArtifacts, err := apiClient.GetJobArtifacts(uint(1))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(jobArtifacts))

	expected := &models.JobArtifact{
		ID:       uint(1),
		JobID:    uint(1),
		Filename: "test.txt",
	}
	assert.Equal(t, expected, jobArtifacts[0])

	response := apiClient.LastRecorder()
	assert.Equal(t, "application/json", response.Result().Header.Get("Content-Type"))
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
}

func TestJobLogTxtGetHandlerNotFound(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	apiClient := testRouter.APIClient

	log, err := apiClient.GetJobLog(uint(1))
	assert.NoError(t, err)
	assert.Nil(t, log)

	response := apiClient.LastRecorder()
	assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
	assert.Equal(t, "text/plain", response.Result().Header.Get("Content-Type"))
	assert.Equal(t, "", response.Body.String())
}
