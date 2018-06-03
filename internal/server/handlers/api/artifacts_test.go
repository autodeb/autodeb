package api_test

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/routertest"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"

	"github.com/stretchr/testify/assert"
)

func TestArtifactGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	apiClient := testRouter.APIClient

	artifact, err := testRouter.Services.Artifacts().CreateArtifact(
		uint(3),
		"test.txt",
		strings.NewReader("test artifact"),
	)
	assert.NoError(t, err)
	assert.NotNil(t, artifact)

	artifact, err = apiClient.GetArtifact(artifact.ID)
	assert.NoError(t, err)

	expected := &models.Artifact{
		ID:       uint(1),
		JobID:    uint(3),
		Filename: "test.txt",
	}
	assert.Equal(t, expected, artifact)

	response := apiClient.LastResponse()
	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestArtifactContentGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	apiClient := testRouter.APIClient

	artifact, err := testRouter.Services.Artifacts().CreateArtifact(
		uint(3),
		"test.txt",
		strings.NewReader("test artifact"),
	)
	assert.NoError(t, err)
	assert.NotNil(t, artifact)

	content, err := apiClient.GetArtifactContent(artifact.ID)
	assert.NoError(t, err)
	assert.NotNil(t, content)

	b, err := ioutil.ReadAll(content)
	assert.NoError(t, err)

	assert.Equal(t, "test artifact", string(b))

	response := apiClient.LastRecorder()
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, "test artifact", response.Body.String())
}
