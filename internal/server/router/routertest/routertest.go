package routertest

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx/appctxtest"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router"
)

//RouterTest allows for testing the server's Router
type RouterTest struct {
	*appctxtest.AppCtxTest
	TemplatesFS filesystem.FS
	Router      http.Handler
}

// SetupTest returns a new RouterTest
func SetupTest(t *testing.T) *RouterTest {
	appCtxTest := appctxtest.SetupTest(t)

	router := router.NewRouter(appCtxTest.AppCtx)

	routerTest := &RouterTest{
		AppCtxTest: appCtxTest,
		Router:     router,
	}

	return routerTest
}

//ServeHTTP serves an http request
func (routerTest *RouterTest) ServeHTTP(request *http.Request) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	routerTest.Router.ServeHTTP(response, request)
	return response
}

//PostForm will post a form
func (routerTest *RouterTest) PostForm(path string, form *url.Values) *httptest.ResponseRecorder {
	request := httptest.NewRequest(
		http.MethodPost,
		path,
		strings.NewReader(form.Encode()),
	)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return routerTest.ServeHTTP(request)
}
