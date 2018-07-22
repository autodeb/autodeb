package api_test

import (
	"net/http"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/routertest"

	"github.com/stretchr/testify/assert"
)

func TestArchiveBackportGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	apiClient := testRouter.APIClient

	expected, err := testRouter.Services.Jobs().CreateArchiveBackport(53, 99)
	assert.NoError(t, err)
	assert.NotNil(t, expected)

	result, err := apiClient.GetArchiveBackport(expected.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	assert.Equal(t, expected, result)

	response := apiClient.LastRecorder()
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Equal(t, "application/json", response.Result().Header.Get("Content-Type"))
}
