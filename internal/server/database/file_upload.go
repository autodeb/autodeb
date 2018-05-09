package database

import (
	"time"

	"github.com/jinzhu/gorm"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

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

// GetFileUpload returns the FileUpload with the given id
func (db *Database) GetFileUpload(id uint) (*models.FileUpload, error) {
	var fileUpload models.FileUpload

	query := db.gormDB.Where(
		&models.FileUpload{
			ID: id,
		},
	)

	err := query.First(&fileUpload).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &fileUpload, nil
}

// GetFileUploadByFileNameSHASumCompleted returns the first FileUpload that matches
func (db *Database) GetFileUploadByFileNameSHASumCompleted(filename, sha256Sum string, completed bool) (*models.FileUpload, error) {
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

// GetAllFileUploadsByUploadID returns all file uploads for an Upload
func (db *Database) GetAllFileUploadsByUploadID(uploadID uint) ([]*models.FileUpload, error) {
	var fileUploads []*models.FileUpload

	query := db.gormDB.Model(
		&models.FileUpload{},
	).Where(
		&models.FileUpload{
			UploadID: uploadID,
		},
	)

	if err := query.Find(&fileUploads).Error; err != nil {
		return nil, err
	}

	return fileUploads, nil
}
