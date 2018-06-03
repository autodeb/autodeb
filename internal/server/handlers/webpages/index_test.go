package webpages_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/routertest"

	"github.com/stretchr/testify/assert"
)

func TestIndexGetHandlerUnauthenticated(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := testRouter.ServeHTTP(request)

	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Contains(t, response.Body.String(), "login")
	assert.NotContains(t, response.Body.String(), "logout")
}

func TestIndexGetHandlerAuthenticated(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	testRouter.Login()

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := testRouter.ServeHTTP(request)

	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.NotContains(t, response.Body.String(), "login")
	assert.Contains(t, response.Body.String(), "logout")
}
