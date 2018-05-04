package app

import (
	"io"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/app/uploads"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

type UploadParameters = uploads.UploadParameters

func (app *App) ProcessUpload(uploadParameters *UploadParameters, content io.Reader) (*models.Upload, error) {
	return app.uploadsManager.ProcessUpload(uploadParameters, content)
}

// GetAllUploads returns all uploads
func (app *App) GetAllUploads() ([]*models.Upload, error) {
	return app.db.GetAllUploads()
}
