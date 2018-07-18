package webpages_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/routertest"
)

func TestArchiveBackportsGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	_, err := testRouter.Services.Jobs().CreateArchiveBackport(14)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, "/archive-backports", nil)
	response := testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Contains(
		t, response.Body.String(),
		"14",
	)

	request = httptest.NewRequest(http.MethodGet, "/archive-upgrades?page=1", nil)
	response = testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.NotContains(
		t, response.Body.String(),
		"421",
	)
}

func TestNewArchiveBackportPostHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	user := testRouter.Login()

	archiveUpgrades, err := testRouter.AppCtx.JobsService().GetAllArchiveBackportsByUserID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(archiveUpgrades))

	request := httptest.NewRequest(http.MethodPost, "/new-archive-backport", nil)
	response := testRouter.ServeHTTP(request)

	assert.Equal(t, http.StatusSeeOther, response.Result().StatusCode)

	archiveUpgrades, err = testRouter.AppCtx.JobsService().GetAllArchiveBackportsByUserID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(archiveUpgrades))
}

func TestArchiveBackportGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	request := httptest.NewRequest(http.MethodGet, "/archive-backports/1", nil)
	response := testRouter.ServeHTTP(request)

	assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
	assert.Equal(t, "", response.Body.String())

	testRouter.Services.Jobs().CreateArchiveBackport(997)

	response = testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Contains(
		t,
		response.Body.String(),
		"997",
	)
}
