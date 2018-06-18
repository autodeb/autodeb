package webpages_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/routertest"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func TestJobsGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	testRouter.DB.CreateJob(models.JobTypeAutopkgtest, "", 0, models.JobParentTypeArchiveUpgrade, 711)
	testRouter.Services.Jobs().CreateForwardJob(312)

	//Show all jobs
	request := httptest.NewRequest(http.MethodGet, "/jobs", nil)
	response := testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Contains(
		t, response.Body.String(),
		"autopkgtest", "archive-upgrade", "711",
		"forward", "upload", "312",
	)

	// There is nothing on the next page
	request = httptest.NewRequest(http.MethodGet, "/jobs?page=1", nil)
	response = testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.NotContains(
		t, response.Body.String(),
		"u1source", "u1version", "u1maint", "u1changedby",
		"autopkgtest", "archive-upgrade", "711",
		"forward", "upload", "312",
	)
}

func TestJobGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	request := httptest.NewRequest(http.MethodGet, "/jobs/1", nil)
	response := testRouter.ServeHTTP(request)

	assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
	assert.Equal(t, "", response.Body.String())

	job, _ := testRouter.Services.Jobs().CreateJob(
		models.JobTypeBuildUpload, "testinput", 444, models.JobParentTypeArchiveUpgrade, 702,
	)
	testRouter.DB.CreateArtifact(job.ID, "testartifactfilename")

	response = testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Contains(
		t, response.Body.String(),
		"build-upload",
		"testartifactfilename",
		"queued",
		"testinput",
		"444",
		"archive-upgrade",
		"702",
	)
}
