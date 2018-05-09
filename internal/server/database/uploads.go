package database

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// CreateUpload will create an upload
func (db *Database) CreateUpload(source, version, maintainer, changedBy string) (*models.Upload, error) {
	upload := &models.Upload{
		Source:     source,
		Version:    version,
		Maintainer: maintainer,
		ChangedBy:  changedBy,
	}

	if err := db.gormDB.Create(upload).Error; err != nil {
		return nil, err
	}

	return upload, nil
}

// GetAllUploads returns all uploads
func (db *Database) GetAllUploads() ([]*models.Upload, error) {
	var uploads []*models.Upload

	if err := db.gormDB.Model(&models.Upload{}).Find(&uploads).Error; err != nil {
		return nil, err
	}

	return uploads, nil
}
