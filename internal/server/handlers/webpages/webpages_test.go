package webpages_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/routertest"

	"github.com/stretchr/testify/assert"
)

func TestWebPagesRender(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	// Pages that render both authenticated and unauthenticated
	pagesNoAuth := []string{
		"/",
		"/jobs",
		"/uploads",
		"/instructions",
		"/archive-upgrades",
		"/archive-backports",
	}

	// Pages that only render when authenticated
	pagesAuth := []string{
		"/profile",
		"/profile/pgp-keys",
		"/profile/access-tokens",
		"/new-archive-upgrade",
	}

	// Test that the pages render when unauthenticated
	for _, page := range pagesNoAuth {
		request := httptest.NewRequest(http.MethodGet, page, nil)
		response := testRouter.ServeHTTP(request)

		assert.Equal(
			t, http.StatusOK, response.Result().StatusCode,
			"this page should render successfully even when unauthenticated",
		)
	}

	// Test that the pages don't render when unauthenticated
	for _, page := range pagesAuth {
		request := httptest.NewRequest(http.MethodGet, page, nil)
		response := testRouter.ServeHTTP(request)

		assert.Equal(
			t, http.StatusSeeOther, response.Result().StatusCode,
			"this page should redirect when unauthenticated",
		)
	}

	testRouter.Login()

	// Test that all pages render when authenticated
	for _, page := range append(pagesNoAuth, pagesAuth...) {
		request := httptest.NewRequest(http.MethodGet, page, nil)
		response := testRouter.ServeHTTP(request)

		assert.Equal(
			t, http.StatusOK, response.Result().StatusCode,
			"this page should render successfully when authenticated",
		)
	}

}
