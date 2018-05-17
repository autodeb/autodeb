package routertest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/htmltemplate"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app/apptest"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router"
)

//RouterTest allows for testing the server's Router
type RouterTest struct {
	App         *app.App
	DataFS      filesystem.FS
	DB          *database.Database
	TemplatesFS filesystem.FS
	StaticFS    filesystem.FS
	Router      http.Handler
}

// SetupTest returns a new RouterTest
func SetupTest(t *testing.T) *RouterTest {
	testApp, dataFS, db := apptest.SetupTest(t)

	tmplFS := filesystem.NewMemMapFs()
	tmplRenderer := htmltemplate.NewRenderer(tmplFS, true)

	staticFS := filesystem.NewMemMapFs()

	router := router.NewRouter(
		tmplRenderer,
		filesystem.NewHTTPFS(staticFS),
		testApp,
		nil,
	)

	routerTest := &RouterTest{
		App:         testApp,
		DataFS:      dataFS,
		DB:          db,
		TemplatesFS: tmplFS,
		StaticFS:    staticFS,
		Router:      router,
	}

	return routerTest
}

func (routerTest *RouterTest) ServeHTTP(request *http.Request) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	routerTest.Router.ServeHTTP(response, request)
	return response
}
