package webpages_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/routertest"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func TestJobGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	request := httptest.NewRequest(http.MethodGet, "/jobs/1", nil)
	response := testRouter.ServeHTTP(request)

	assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
	assert.Equal(t, "", response.Body.String())

	job, _ := testRouter.Services.Jobs().CreateJob(
		models.JobTypeBuild, "testinput", models.JobParentTypeArchiveUpgrade, 702,
	)
	testRouter.DB.CreateArtifact(job.ID, "testartifactfilename")

	response = testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Contains(
		t, response.Body.String(),
		"testartifactfilename",
		"queued",
		"testinput",
		"archive-upgrade",
		"702",
	)
}
