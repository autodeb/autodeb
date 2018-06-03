package webpages_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/routertest"

	"github.com/stretchr/testify/assert"
)

func TestProfileGetHandlerAuthenticated(t *testing.T) {
	testRouter := routertest.SetupTest(t)
	user := testRouter.Login()

	request := httptest.NewRequest(http.MethodGet, "/profile", nil)
	response := testRouter.ServeHTTP(request)

	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Contains(t, response.Body.String(), user.Username)

	testRouter.Logout()

	request = httptest.NewRequest(http.MethodGet, "/profile", nil)
	response = testRouter.ServeHTTP(request)

	assert.Equal(t, http.StatusSeeOther, response.Result().StatusCode)
	assert.NotContains(t, response.Body.String(), user.Username)
}
