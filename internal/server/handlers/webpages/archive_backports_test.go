package webpages_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/routertest"
)

func TestArchiveBackportsGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	_, err := testRouter.Services.Jobs().CreateArchiveBackport(14, 54)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, "/archive-backports", nil)
	response := testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Contains(
		t, response.Body.String(),
		"14",
		"54",
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

	archiveBackports, err := testRouter.AppCtx.JobsService().GetAllArchiveBackportsByUserID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(archiveBackports))

	form := &url.Values{}
	form.Add("package-count", "43")

	response := testRouter.PostForm("/new-archive-backport", form)
	assert.Equal(t, http.StatusSeeOther, response.Result().StatusCode)

	archiveBackports, err = testRouter.AppCtx.JobsService().GetAllArchiveBackportsByUserID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(archiveBackports))
}

func TestArchiveBackportGetHandler(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	request := httptest.NewRequest(http.MethodGet, "/archive-backports/1", nil)
	response := testRouter.ServeHTTP(request)

	assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
	assert.Equal(t, "", response.Body.String())

	testRouter.Services.Jobs().CreateArchiveBackport(997, 24)

	response = testRouter.ServeHTTP(request)
	assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	assert.Contains(
		t,
		response.Body.String(),
		"997",
		"24",
	)
}
