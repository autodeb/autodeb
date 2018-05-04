// Package app implements most of the application logic
package app

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
)

// App is an autodeb server application
type App struct {
	db     *database.Database
	dataFS filesystem.FS
}

// NewApp create an app from a configuration
func NewApp(db *database.Database, dataFS filesystem.FS) (*App, error) {

	app := App{
		db:     db,
		dataFS: dataFS,
	}

	if err := app.setupDataDirectory(); err != nil {
		return nil, err
	}

	return &app, nil
}

// UploadedFilesDirectory contains files that are not yet associated
// with a package upload.
func (app *App) UploadedFilesDirectory() string {
	return "/files"
}

// UploadsDirectory contains completed uploads.
func (app *App) UploadsDirectory() string {
	return "/uploads"
}
