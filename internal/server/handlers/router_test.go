package router_test

import (
	"net/http"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/routertest"

	"github.com/stretchr/testify/assert"
)

func TestStaticFilesNotFound(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	request, _ := http.NewRequest(http.MethodGet, "/static/test", nil)

	response := testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
}

func TestStaticFiles(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	testRouter.StaticFS.Create("test")

	request, _ := http.NewRequest(http.MethodGet, "/static/test", nil)

	response := testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
}
