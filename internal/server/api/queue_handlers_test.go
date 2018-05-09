package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"

	"github.com/stretchr/testify/assert"
)

func TestNextJobNoJob(t *testing.T) {
	testAPI := setupTest(t)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/api/queue/next", nil)

	testAPI.API.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNoContent, response.Result().StatusCode)
	assert.Equal(t, response.Body.String(), "")
}

func TestNextJob(t *testing.T) {
	testAPI := setupTest(t)
	testAPI.Database.CreateJob(models.JobTypeBuild, uint(3))

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/api/queue/next", nil)

	testAPI.API.ServeHTTP(response, request)

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
