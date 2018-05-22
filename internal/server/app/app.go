// Package app implements most of the application logic, it contains
// everything that is needed to serve a request.
package app

import (
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/htmltemplate"
	"salsa.debian.org/autodeb-team/autodeb/internal/log"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app/services/jobs"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app/services/pgp"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app/services/uploads"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/auth"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
)

// App is an autodeb server application
type App struct {
	config      *Config
	dataFS      filesystem.FS
	renderer    *htmltemplate.Renderer
	staticFS    http.FileSystem
	authBackend auth.Backend
	logger      log.Logger

	// Services
	uploadsService *uploads.Service
	pgpService     *pgp.Service
	jobsService    *jobs.Service
}

// NewApp create an app from a configuration
func NewApp(
	config *Config,
	db *database.Database,
	dataFS filesystem.FS,
	renderer *htmltemplate.Renderer,
	staticFS http.FileSystem,
	authBackend auth.Backend,
	logger log.Logger) (*App, error) {

	app := App{
		config:      config,
		dataFS:      dataFS,
		renderer:    renderer,
		staticFS:    staticFS,
		authBackend: authBackend,
		logger:      logger,

		// Services
		uploadsService: uploads.New(db, dataFS),
		pgpService:     pgp.New(db, config.ServerURL),
		jobsService:    jobs.New(db, dataFS),
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

// AuthBackend returns the authentification service
func (app *App) AuthBackend() auth.Backend {
	return app.authBackend
}

// Config returns the app's config
func (app *App) Config() *Config {
	return app.config
}

// StaticFS contains static files to be served over http
func (app *App) StaticFS() http.FileSystem {
	return app.staticFS
}

// UploadsService returns the uploads service
func (app *App) UploadsService() *uploads.Service {
	return app.uploadsService
}

// PGPService returns the PGP service
func (app *App) PGPService() *pgp.Service {
	return app.pgpService
}

// JobsService returns the jobs service
func (app *App) JobsService() *jobs.Service {
	return app.jobsService
}

// TemplatesRenderer returns the template renderer
func (app *App) TemplatesRenderer() *htmltemplate.Renderer {
	return app.renderer
}

// UploadedFilesDirectory contains files that are not yet associated
// with a package upload.
func (app *App) UploadedFilesDirectory() string {
	return app.uploadsService.UploadedFilesDirectory()
}

// UploadsDirectory contains completed uploads.
func (app *App) UploadsDirectory() string {
	return app.uploadsService.UploadsDirectory()
}

// JobsDirectory contains saved data for jobs such as logs
func (app *App) JobsDirectory() string {
	return app.jobsService.JobsDirectory()
}
