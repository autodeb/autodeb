// Package app implements most of the application logic
package app

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app/uploads"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
)

// App is an autodeb server application
type App struct {
	db             *database.Database
	dataFS         filesystem.FS
	uploadsManager *uploads.Manager
}

// NewApp create an app from a configuration
func NewApp(db *database.Database, dataFS filesystem.FS) (*App, error) {

	app := App{
		db:             db,
		dataFS:         dataFS,
		uploadsManager: uploads.NewManager(db, dataFS),
	}

	if err := app.setupDataDirectory(); err != nil {
		return nil, err
	}

	return &app, nil
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
