// Package app implements most of the application logic, it contains
// everything that is needed to serve a request.
package app

import (
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/htmltemplate"
	"salsa.debian.org/autodeb-team/autodeb/internal/oauth"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app/uploads"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
)

// App is an autodeb server application
type App struct {
	db             *database.Database
	dataFS         filesystem.FS
	uploadsManager *uploads.Manager
	oauthProvider  oauth.Provider
	renderer       *htmltemplate.Renderer
	staticFS       http.FileSystem
}

// NewApp create an app from a configuration
func NewApp(
	db *database.Database,
	dataFS filesystem.FS,
	oauthProvider oauth.Provider,
	renderer *htmltemplate.Renderer,
	staticFS http.FileSystem) (*App, error) {

	app := App{
		db:             db,
		dataFS:         dataFS,
		uploadsManager: uploads.NewManager(db, dataFS),
		oauthProvider:  oauthProvider,
		renderer:       renderer,
		staticFS:       staticFS,
	}

	if err := app.setupDataDirectory(); err != nil {
		return nil, err
	}

	return &app, nil
}

// StaticFS contains static files to be served over http
func (app *App) StaticFS() http.FileSystem {
	return app.staticFS
}

// TemplatesRenderer returns the template renderer
func (app *App) TemplatesRenderer() *htmltemplate.Renderer {
	return app.renderer
}

// OAuthProvider returns the configured OAuth provider
func (app *App) OAuthProvider() oauth.Provider {
	return app.oauthProvider
}

// UploadedFilesDirectory contains files that are not yet associated
// with a package upload.
func (app *App) UploadedFilesDirectory() string {
	return app.uploadsManager.UploadedFilesDirectory()
}

// UploadsDirectory contains completed uploads.
func (app *App) UploadsDirectory() string {
	return app.uploadsManager.UploadsDirectory()
}

// JobsDirectory contains saved data for jobs such as logs
func (app *App) JobsDirectory() string {
	return "/jobs"
}
