// Package app implements most of the application logic
package app

import (
	"path/filepath"

	"salsa.debian.org/aviau/autodeb/internal/server/database"
)

// App is an autodeb server application
type App struct {
	config    *Config
	dataStore *database.Database
}

// NewApp create an app from a configuration
func NewApp(cfg *Config, dataStore *database.Database) (*App, error) {

	app := App{
		config:    cfg,
		dataStore: dataStore,
	}

	if err := app.setupDataDirectory(); err != nil {
		return nil, err
	}

	return &app, nil
}

// UploadedFilesDirectory contains files that are not yet associated
// with a package upload.
func (app *App) UploadedFilesDirectory() string {
	return filepath.Join(app.config.DataDirectory, "files")
}

// UploadsDirectory contains completed uploads.
func (app *App) UploadsDirectory() string {
	return filepath.Join(app.config.DataDirectory, "uploads")
}
