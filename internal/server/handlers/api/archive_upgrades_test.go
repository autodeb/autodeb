package api_test

import (
	"net/http"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/routertest"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"

	"github.com/stretchr/testify/assert"
)

func TestArchiveUpgradeGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	apiClient := testRouter.APIClient

	expected, err := testRouter.Services.Jobs().CreateArchiveUpgrade(33, 13)
	assert.NoError(t, err)
	assert.NotNil(t, expected)

	result, err := apiClient.GetArchiveUpgrade(expected.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, expected, result)

	response := apiClient.LastRecorder()
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, "application/json", response.Result().Header.Get("Content-Type"))
}

func TestArchiveUpgradeJobsGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	apiClient := testRouter.APIClient

	archiveUpgrade, err := testRouter.Services.Jobs().CreateArchiveUpgrade(33, 13)
	assert.NoError(t, err)
	assert.NotNil(t, archiveUpgrade)

	jobs, err := apiClient.GetArchiveUpgradeJobs(archiveUpgrade.ID)
	assert.NoError(t, err)
	assert.NotNil(t, jobs)
	assert.Equal(t, 1, len(jobs))

	response := apiClient.LastRecorder()
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, "application/json", response.Result().Header.Get("Content-Type"))
}

func TestArchiveUpgradeSuccessfulBuildsGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	apiClient := testRouter.APIClient

	archiveUpgrade, err := testRouter.Services.Jobs().CreateArchiveUpgrade(33, 13)
	assert.NoError(t, err)
	assert.NotNil(t, archiveUpgrade)

	// Create a build that with a successful autopkgtest job
	build, err := testRouter.Services.Jobs().CreateJob(
		models.JobTypePackageUpgrade,
		"gccc",
		0,
		models.JobParentTypeArchiveUpgrade,
		archiveUpgrade.ID,
	)
	assert.NoError(t, err)

	autopkgtest, err := testRouter.Services.Jobs().CreateAutopkgtestJobFromBuildJob(build)
	assert.NoError(t, err)

	err = testRouter.Services.Jobs().ProcessJobStatus(autopkgtest.ID, models.JobStatusSuccess)
	assert.NoError(t, err)

	// Create a second build, without a successful autopkgtest job
	_, err = testRouter.Services.Jobs().CreateJob(
		models.JobTypePackageUpgrade,
		"gccc",
		0,
		models.JobParentTypeArchiveUpgrade,
		archiveUpgrade.ID,
	)
	assert.NoError(t, err)

	jobs, err := apiClient.GetArchiveUpgradeSuccessfulBuilds(archiveUpgrade.ID)
	assert.NoError(t, err)
	assert.NotNil(t, jobs)
	assert.Equal(t, 1, len(jobs))
	assert.Equal(t, build.ID, jobs[0].ID)

	response := apiClient.LastRecorder()
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, "application/json", response.Result().Header.Get("Content-Type"))
}
