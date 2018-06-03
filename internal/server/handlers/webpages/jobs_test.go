package webpages_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/routertest"
)

func TestJobGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	request := httptest.NewRequest(http.MethodGet, "/jobs/1", nil)
	response := testRouter.ServeHTTP(request)

	assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
	assert.Equal(t, "", response.Body.String())

	job, _ := testRouter.Services.Jobs().CreateBuildJob(1)
	testRouter.DB.CreateArtifact(job.ID, "testartifactfilename")

	response = testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Contains(t, response.Body.String(), "testartifactfilename")
	assert.Contains(t, response.Body.String(), "queued")
}
