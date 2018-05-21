package apptest

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/htmltemplate"
	"salsa.debian.org/autodeb-team/autodeb/internal/log"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

//AppTest contains an app and all of its dependencies for testing
type AppTest struct {
	t           *testing.T
	App         *app.App
	DataFS      filesystem.FS
	StaticFS    filesystem.FS
	DB          *database.Database
	authBackend *fakeAuthBackend
}

//SetupTest will create a test App
func SetupTest(t *testing.T) *AppTest {
	config := &app.Config{
		ServerURL: "https://test.auto.debian.net",
	}

	db, err := database.NewDatabase("sqlite3", ":memory:")
	require.NoError(t, err)

	dataFS := filesystem.NewMemMapFs()

	templatesFS, err := filesystem.NewFS(
		filepath.Join(projectDirectory(), "web", "templates"),
	)
	require.NoError(t, err)

	tmplRenderer := htmltemplate.NewRenderer(templatesFS, true)

	staticFS := filesystem.NewMemMapFs()

	authBackend := newFakeAuthBackend()

	logger := log.New(ioutil.Discard)

	app, err := app.NewApp(
		config,
		db,
		dataFS,
		tmplRenderer,
		filesystem.NewHTTPFS(staticFS),
		authBackend,
		logger,
	)
	require.NoError(t, err)

	appTest := &AppTest{
		t:           t,
		App:         app,
		DataFS:      dataFS,
		StaticFS:    staticFS,
		DB:          db,
		authBackend: authBackend,
	}

	return appTest
}

// Login will create a test user and future requests will be authenticated
// as this user
func (appTest *AppTest) Login() *models.User {
	user, err := appTest.DB.GetUser(uint(1))
	require.NoError(appTest.t, err)

	if user == nil {
		user, err = appTest.DB.CreateUser(1, "testuser3579")
	}

	appTest.authBackend.User = user

	return user
}

// Logout will logout the currently logged user
func (appTest *AppTest) Logout() {
	appTest.authBackend.User = nil
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
