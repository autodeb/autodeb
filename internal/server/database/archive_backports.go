package database

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"

	"github.com/jinzhu/gorm"
)

// CreateArchiveBackport will create an ArchiveBackport
func (db *Database) CreateArchiveBackport(userID uint, packageCount int) (*models.ArchiveBackport, error) {
	archiveBackport := &models.ArchiveBackport{
		UserID:       userID,
		PackageCount: packageCount,
	}

	if err := db.gormDB.Create(archiveBackport).Error; err != nil {
		return nil, err
	}

	return archiveBackport, nil
}

// GetArchiveBackport returns the ArchiveBackport with the given id
func (db *Database) GetArchiveBackport(id uint) (*models.ArchiveBackport, error) {
	var archiveBackport models.ArchiveBackport

	query := db.gormDB.Where(
		&models.ArchiveBackport{
			ID: id,
		},
	)

	err := query.First(&archiveBackport).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &archiveBackport, nil
}

// GetAllArchiveBackports returns all ArchiceBackports
func (db *Database) GetAllArchiveBackports() ([]*models.ArchiveBackport, error) {
	var archiveBackports []*models.ArchiveBackport

	if err := db.gormDB.Model(&models.ArchiveBackport{}).Find(&archiveBackports).Error; err != nil {
		return nil, err
	}

	return archiveBackports, nil
}

// GetAllArchiveBackportsPageLimit returns all ArchiveBackports with pagination
func (db *Database) GetAllArchiveBackportsPageLimit(page, limit int) ([]*models.ArchiveBackport, error) {
	offset := page * limit

	var archiveBackports []*models.ArchiveBackport

	query := db.gormDB.Model(
		&models.ArchiveBackport{},
	).Order(
		"id desc",
	).Offset(
		offset,
	).Limit(
		limit,
	)

	if err := query.Find(&archiveBackports).Error; err != nil {
		return nil, err
	}

	return archiveBackports, nil
}

// GetAllArchiveBackportsByUserID returns all ArchiveBackports for a UserID
func (db *Database) GetAllArchiveBackportsByUserID(userID uint) ([]*models.ArchiveBackport, error) {
	var archiveBackports []*models.ArchiveBackport

	query := db.gormDB.Model(
		&models.ArchiveBackport{},
	).Where(
		&models.ArchiveBackport{
			UserID: userID,
		},
	)

	if err := query.Find(&archiveBackports).Error; err != nil {
		return nil, err
	}

	return archiveBackports, nil
}
