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

// CreateFileUpload will create a FileUpload
func (db *Database) CreateFileUpload(filename, sha256Sum string, uploadedAt time.Time) (*models.FileUpload, error) {
	fileUpload := &models.FileUpload{
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

// GetFileUpload will return the first FileUpload that matches
func (db *Database) GetFileUpload(filename, sha256Sum string, completed bool) (*models.FileUpload, error) {
	var fileUpload models.FileUpload

	query := db.gormDB.Where(
		&models.FileUpload{
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

// UpdateFileUpload will update a file upload
func (db *Database) UpdateFileUpload(fileUpload *models.FileUpload) error {
	err := db.gormDB.Save(fileUpload).Error
	return err
}
