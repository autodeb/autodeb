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
	authDisabled "salsa.debian.org/autodeb-team/autodeb/internal/server/auth/disabled"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
)

//AppTest contains an app and all of its dependencies for testing
type AppTest struct {
	App      *app.App
	DataFS   filesystem.FS
	StaticFS filesystem.FS
	DB       *database.Database
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

	authBackend := authDisabled.NewBackend()

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
		App:      app,
		DataFS:   dataFS,
		StaticFS: staticFS,
		DB:       db,
	}

	return appTest
}
