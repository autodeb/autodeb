package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStaticFilesNotFound(t *testing.T) {
	testAPI := setupTest(t)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/static/test", nil)

	testAPI.API.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
}

func TestStaticFiles(t *testing.T) {
	testAPI := setupTest(t)
	testAPI.StaticFS.Create("test")

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/static/test", nil)

	testAPI.API.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
}
