package app

import (
	"os"
)

func (app *App) setupDataDirectory() error {

	// Data directory
	if _, err := os.Stat(app.config.DataDirectory); os.IsNotExist(err) {
		if err := os.Mkdir(app.config.DataDirectory, 0744); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// data/files
	if _, err := os.Stat(app.UploadedFilesDirectory()); os.IsNotExist(err) {
		if err = os.Mkdir(app.UploadedFilesDirectory(), 0744); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// data/uploads
	if _, err := os.Stat(app.UploadsDirectory()); os.IsNotExist(err) {
		if err := os.Mkdir(app.UploadsDirectory(), 0744); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}
