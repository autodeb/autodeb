package api_test

import (
	"net/http"
	"testing"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/htmltemplate"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/api"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app/apptest"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
)

type APITest struct {
	App         *app.App
	DataFS      filesystem.FS
	Database    *database.Database
	TemplatesFS filesystem.FS
	StaticFS    filesystem.FS
	API         http.Handler
}

func setupTest(t *testing.T) *APITest {
	testApp, dataFS, db := apptest.SetupTest(t)

	tmplFS := filesystem.NewMemMapFs()
	tmplRenderer := htmltemplate.NewRenderer(tmplFS, true)

	staticFS := filesystem.NewMemMapFs()

	router := api.NewRouter(
		tmplRenderer,
		filesystem.NewHTTPFS(staticFS),
		testApp,
	)

	apiTest := &APITest{
		App:         testApp,
		DataFS:      dataFS,
		Database:    db,
		TemplatesFS: tmplFS,
		StaticFS:    staticFS,
		API:         router,
	}

	return apiTest
}
