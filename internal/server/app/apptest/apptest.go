package apptest

import (
	"io/ioutil"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/require"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/htmltemplate"
	"salsa.debian.org/autodeb-team/autodeb/internal/log"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
)

//AppTest contains an app and all of its dependencies for testing
type AppTest struct {
	App      *app.App
	DataFS   filesystem.FS
	StaticFS filesystem.FS
	DB       *database.Database
}

//SetupTest will create a test App
func SetupTest(t *testing.T) *AppTest {
	config := &app.Config{
		ServerURL: "https://test.auto.debian.net",
	}

	db, err := database.NewDatabase(
		"sqlite3",
		":memory:",
	)
	require.NoError(t, err)

	dataFS := filesystem.NewMemMapFs()

	tmplRenderer := htmltemplate.NewRenderer(
		filesystem.NewMemMapFs(),
		true,
	)

	staticFS := filesystem.NewMemMapFs()

	sessionsStore := sessions.NewCookieStore([]byte("something-very-secret"))

	logger := log.New(ioutil.Discard)

	app, err := app.NewApp(
		config,
		db,
		dataFS,
		nil,
		tmplRenderer,
		filesystem.NewHTTPFS(staticFS),
		sessionsStore,
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
