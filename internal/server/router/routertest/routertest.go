package routertest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app/apptest"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router"
)

//RouterTest allows for testing the server's Router
type RouterTest struct {
	*apptest.AppTest
	TemplatesFS filesystem.FS
	Router      http.Handler
}

// SetupTest returns a new RouterTest
func SetupTest(t *testing.T) *RouterTest {
	appTest := apptest.SetupTest(t)

	router := router.NewRouter(appTest.App)

	routerTest := &RouterTest{
		AppTest: appTest,
		Router:  router,
	}

	return routerTest
}

func (routerTest *RouterTest) ServeHTTP(request *http.Request) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	routerTest.Router.ServeHTTP(response, request)
	return response
}
