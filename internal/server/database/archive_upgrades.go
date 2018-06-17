package database

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"

	"github.com/jinzhu/gorm"
)

// CreateArchiveUpgrade will create an ArchiveUpgrade
func (db *Database) CreateArchiveUpgrade(userID uint, packageCount uint) (*models.ArchiveUpgrade, error) {
	archiveUpgrade := &models.ArchiveUpgrade{
		UserID:       userID,
		PackageCount: packageCount,
	}

	if err := db.gormDB.Create(archiveUpgrade).Error; err != nil {
		return nil, err
	}

	return archiveUpgrade, nil
}

// GetArchiveUpgrade returns the ArchiveUpgrade with the given id
func (db *Database) GetArchiveUpgrade(id uint) (*models.ArchiveUpgrade, error) {
	var archiveUpgrade models.ArchiveUpgrade

	query := db.gormDB.Where(
		&models.ArchiveUpgrade{
			ID: id,
		},
	)

	err := query.First(&archiveUpgrade).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &archiveUpgrade, nil
}

// GetAllArchiveUpgrades returns all ArchiveUpgrades
func (db *Database) GetAllArchiveUpgrades() ([]*models.ArchiveUpgrade, error) {
	var archiveUpgrades []*models.ArchiveUpgrade

	if err := db.gormDB.Model(&models.ArchiveUpgrade{}).Find(&archiveUpgrades).Error; err != nil {
		return nil, err
	}

	return archiveUpgrades, nil
}

// GetAllArchiveUpgradesPageLimit returns all ArchiveUpgrades with pagination
func (db *Database) GetAllArchiveUpgradesPageLimit(page, limit int) ([]*models.ArchiveUpgrade, error) {
	offset := page * limit

	var archiveUpgrades []*models.ArchiveUpgrade

	query := db.gormDB.Model(
		&models.ArchiveUpgrade{},
	).Order(
		"id desc",
	).Offset(
		offset,
	).Limit(
		limit,
	)

	if err := query.Find(&archiveUpgrades).Error; err != nil {
		return nil, err
	}

	return archiveUpgrades, nil
}

// GetAllArchiveUpgradesByUserID returns all ArchiveUpgrades for a UserID
func (db *Database) GetAllArchiveUpgradesByUserID(userID uint) ([]*models.ArchiveUpgrade, error) {
	var archiveUpgrades []*models.ArchiveUpgrade

	query := db.gormDB.Model(
		&models.ArchiveUpgrade{},
	).Where(
		&models.ArchiveUpgrade{
			UserID: userID,
		},
	)

	if err := query.Find(&archiveUpgrades).Error; err != nil {
		return nil, err
	}

	return archiveUpgrades, nil
}
