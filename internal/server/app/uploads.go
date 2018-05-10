package app

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/app/uploads"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

//UploadParameters define upload options
type UploadParameters = uploads.UploadParameters

// ProcessUpload processes uploads
func (app *App) ProcessUpload(uploadParameters *UploadParameters, content io.Reader) (*models.Upload, error) {
	return app.uploadsManager.ProcessUpload(uploadParameters, content)
}

// GetAllUploads returns all uploads
func (app *App) GetAllUploads() ([]*models.Upload, error) {
	return app.db.GetAllUploads()
}

//GetAllFileUploadsByUploadID returns all FileUploads associated to an upload
func (app *App) GetAllFileUploadsByUploadID(uploadID uint) ([]*models.FileUpload, error) {
	fileUploads, err := app.db.GetAllFileUploadsByUploadID(uploadID)
	if err != nil {
		return nil, err
	}
	return fileUploads, nil
}

// GetUploadDSC returns the DSC of the upload with a matching id
func (app *App) GetUploadDSC(uploadID uint) (io.ReadCloser, error) {
	fileUploads, err := app.GetAllFileUploadsByUploadID(uploadID)
	if err != nil {
		return nil, err
	}

	for _, fileUpload := range fileUploads {
		if strings.HasSuffix(fileUpload.Filename, ".dsc") {
			return app.GetUploadFile(uploadID, fileUpload.Filename)
		}
	}

	return nil, nil
}

// GetUploadFile returns the file associated with the upload id and filename
func (app *App) GetUploadFile(uploadID uint, filename string) (io.ReadCloser, error) {
	file, err := app.dataFS.Open(
		filepath.Join(
			app.UploadsDirectory(),
			fmt.Sprint(uploadID),
			filename,
		),
	)
	if err != nil {
		return nil, err
	}
	return file, nil
}
