package app

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// GetAllUploads returns all uploads
func (app *App) GetAllUploads() ([]*models.Upload, error) {
	return app.db.GetAllUploads()
}
