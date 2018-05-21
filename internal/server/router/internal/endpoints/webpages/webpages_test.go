package webpages_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/routertest"

	"github.com/stretchr/testify/assert"
)

func TestWebPagesRenderUnauthenticated(t *testing.T) {
	testRouter := routertest.SetupTest(t)

	pages := []string{
		"/",
		"/jobs",
		"/uploads",
	}

	for _, page := range pages {
		request := httptest.NewRequest(http.MethodGet, page, nil)
		response := testRouter.ServeHTTP(request)

		assert.Equal(
			t,
			http.StatusOK,
			response.Result().StatusCode,
			"this page should render successfully even when unauthenticated",
		)
	}
}
