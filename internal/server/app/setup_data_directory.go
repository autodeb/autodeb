package app

import (
	"os"
)

func (app *App) setupDataDirectory() error {

	// UploadedFilesdirectory
	if _, err := app.dataFS.Stat(app.UploadedFilesDirectory()); os.IsNotExist(err) {
		if err := app.dataFS.Mkdir(app.UploadedFilesDirectory(), 0744); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// UploadsDirectory
	if _, err := app.dataFS.Stat(app.UploadsDirectory()); os.IsNotExist(err) {
		if err := app.dataFS.Mkdir(app.UploadsDirectory(), 0744); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}
