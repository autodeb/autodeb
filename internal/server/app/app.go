// Package app implements most of the application logic, it contains
// everything that is needed to serve a request.
package app

import (
	"net/http"

	"github.com/gorilla/sessions"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/htmltemplate"
	"salsa.debian.org/autodeb-team/autodeb/internal/log"
	"salsa.debian.org/autodeb-team/autodeb/internal/oauth"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app/auth"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app/uploads"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
)

// App is an autodeb server application
type App struct {
	db             *database.Database
	config         *Config
	dataFS         filesystem.FS
	uploadsManager *uploads.Manager
	oauthProvider  oauth.Provider
	renderer       *htmltemplate.Renderer
	staticFS       http.FileSystem
	sessionStore   sessions.Store
	logger         log.Logger
}

// NewApp create an app from a configuration
func NewApp(
	config *Config,
	db *database.Database,
	dataFS filesystem.FS,
	oauthProvider oauth.Provider,
	renderer *htmltemplate.Renderer,
	staticFS http.FileSystem,
	sessionsStore sessions.Store,
	logger log.Logger) (*App, error) {

	app := App{
		config:         config,
		db:             db,
		dataFS:         dataFS,
		uploadsManager: uploads.NewManager(db, dataFS),
		oauthProvider:  oauthProvider,
		renderer:       renderer,
		staticFS:       staticFS,
		sessionStore:   sessionsStore,
		logger:         logger,
	}

	if err := app.setupDataDirectory(); err != nil {
		return nil, err
	}

	return &app, nil
}

// Logger returns the logger
func (app *App) Logger() log.Logger {
	return app.logger
}

// AuthService returns the authentification service
func (app *App) AuthService() *auth.Service {
	return auth.NewService(app.db, app.sessionStore)
}

// Config returns the app's config
func (app *App) Config() *Config {
	return app.config
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
