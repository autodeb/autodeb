package appctxtest

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"

	gorillaSessions "github.com/gorilla/sessions"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/htmltemplate"
	"salsa.debian.org/autodeb-team/autodeb/internal/http/sessions"
	"salsa.debian.org/autodeb-team/autodeb/internal/log"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/config"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/services/servicestest"
)

//AppCtxTest contains an app and all of its dependencies for testing
type AppCtxTest struct {
	*servicestest.ServicesTest
	t           *testing.T
	AppCtx      *appctx.Context
	StaticFS    filesystem.FS
	authBackend *fakeAuthBackend
}

//SetupTest will create a test App
func SetupTest(t *testing.T) *AppCtxTest {
	servicesTest := servicestest.SetupTest(t)

	config := &config.Config{
		ServerURL: servicesTest.ServerURL,
	}

	templatesFS := filesystem.NewBasePathFS(
		filesystem.NewOsFS(),
		filepath.Join(projectDirectory(), "web", "templates"),
	)

	tmplRenderer := htmltemplate.NewRenderer(templatesFS, true)

	staticFS := filesystem.NewMemMapFS()

	authBackend := newFakeAuthBackend()

	sessionsManager := sessions.NewManager(
		gorillaSessions.NewCookieStore([]byte("autodeb-tests")),
		"autodeb-tests",
	)

	logger := log.New(ioutil.Discard)

	appCtx := appctx.New(
		config,
		tmplRenderer,
		filesystem.NewHTTPFS(staticFS),
		authBackend,
		sessionsManager,
		servicesTest.Services,
		logger,
	)

	appCtxTest := &AppCtxTest{
		ServicesTest: servicesTest,
		t:            t,
		AppCtx:       appCtx,
		StaticFS:     staticFS,
		authBackend:  authBackend,
	}

	return appCtxTest
}

// Login will create a test user and future requests will be authenticated
// as this user
func (appCtxTest *AppCtxTest) Login() *models.User {
	user := appCtxTest.GetOrCreateTestUser()
	appCtxTest.authBackend.User = user
	return user
}

// Logout will logout the currently logged user
func (appCtxTest *AppCtxTest) Logout() {
	appCtxTest.authBackend.User = nil
}

func projectDirectory() string {
	_, sourceFile, _, _ := runtime.Caller(0)
	// apptest
	dir := filepath.Dir(sourceFile)
	// app
	dir = filepath.Dir(dir)
	// server
	dir = filepath.Dir(dir)
	// internal
	dir = filepath.Dir(dir)
	// autodeb
	dir = filepath.Dir(dir)
	return dir
}
