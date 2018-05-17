package apptest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/htmltemplate"
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

	app, err := app.NewApp(
		db,
		dataFS,
		nil,
		tmplRenderer,
		filesystem.NewHTTPFS(staticFS),
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
