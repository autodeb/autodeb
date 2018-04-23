package database

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// CreateUpload will create an upload
func (db *Database) CreateUpload() (*models.Upload, error) {
	upload := &models.Upload{}

	if err := db.gormDB.Create(upload).Error; err != nil {
		return nil, err
	}

	return upload, nil
}
