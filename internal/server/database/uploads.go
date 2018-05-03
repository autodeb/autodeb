package database

import (
	"time"

	"github.com/jinzhu/gorm"

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

// CreatePendingFileUpload will create a pending file upload
func (db *Database) CreatePendingFileUpload(filename, sha256Sum string, uploadedAt time.Time) (*models.PendingFileUpload, error) {
	fileUpload := &models.PendingFileUpload{
		Filename:   filename,
		SHA256Sum:  sha256Sum,
		UploadedAt: uploadedAt,
		Completed:  false,
	}

	if err := db.gormDB.Create(fileUpload).Error; err != nil {
		return nil, err
	}

	return fileUpload, nil
}

// GetPendingFileUpload will return the first pending file upload that matches
func (db *Database) GetPendingFileUpload(filename, sha256Sum string, completed bool) (*models.PendingFileUpload, error) {
	var fileUpload models.PendingFileUpload

	query := db.gormDB.Where(
		&models.PendingFileUpload{
			Filename:  filename,
			SHA256Sum: sha256Sum,
			Completed: completed,
		},
	)

	// Fields with possible zero values
	query = query.Where("completed = ?", completed)

	err := query.First(&fileUpload).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &fileUpload, nil
}

// UpdatePendingFileUpload will a file upload
func (db *Database) UpdatePendingFileUpload(fileUpload *models.PendingFileUpload) error {
	err := db.gormDB.Save(fileUpload).Error
	return err
}
