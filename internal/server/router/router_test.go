package router_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/routertest"

	"github.com/stretchr/testify/assert"
)

func TestStaticFilesNotFound(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/static/test", nil)

	testRouter.Router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
}

func TestStaticFiles(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	testRouter.StaticFS.Create("test")

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/static/test", nil)

	testRouter.Router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
}
