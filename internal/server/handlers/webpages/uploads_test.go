package webpages_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/routertest"

	"github.com/stretchr/testify/assert"
)

func TestUploadsGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	testRouter.DB.CreateUpload(1, "u1source", "u1version", "u1maint", "u1changedby", false, false)
	testRouter.DB.CreateUpload(2, "u2source", "u2version", "u2maint", "u2changedby", false, false)

	//Show all uploads
	request := httptest.NewRequest(http.MethodGet, "/uploads/1", nil)
	response := testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Contains(
		t, response.Body.String(),
		"u1source", "u1version", "u1maint", "u1changedby",
		"u2source", "u2version", "u2maint", "u2changedby",
	)

	//Show only a subset of uploads
	request = httptest.NewRequest(http.MethodGet, "/uploads?user_id=2", nil)
	response = testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Contains(
		t, response.Body.String(),
		"u2source", "u2version", "u2maint", "u2changedby",
	)
	assert.NotContains(
		t, response.Body.String(),
		"u1source", "u1version", "u1maint", "u1changedby",
	)
}

func TestUploadGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	request := httptest.NewRequest(http.MethodGet, "/uploads/1", nil)
	response := testRouter.ServeHTTP(request)

	assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
	assert.Equal(t, "", response.Body.String())

	testRouter.DB.CreateUpload(
		1,
		"testsourcename",
		"testversion",
		"testmaintainer",
		"testchangedbyname",
		false,
		false,
	)

	response = testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Contains(
		t,
		response.Body.String(),
		"testsourcename",
		"testversion",
		"testmaintainer",
		"testchangedbyname",
		false,
	)
}
