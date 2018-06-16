package jobs

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// CreateArchiveUpgrade creates a new ArchiveUpgrade
func (service *Service) CreateArchiveUpgrade(userID uint, packageCount uint) (*models.ArchiveUpgrade, error) {
	archiveUpgrade, err := service.db.CreateArchiveUpgrade(userID, packageCount)
	if err != nil {
		return nil, err
	}

	if _, err := service.db.CreateJob(
		models.JobTypeSetupArchiveUpgrade,
		"",
		models.JobParentTypeArchiveUpgrade,
		archiveUpgrade.ID,
	); err != nil {
		return nil, errors.WithMessagef(err, "could not create archive rebuild setup job")
	}

	return archiveUpgrade, err
}

// GetArchiveUpgrade returns the ArchiveUpgrade with a matching ID
func (service *Service) GetArchiveUpgrade(id uint) (*models.ArchiveUpgrade, error) {
	return service.db.GetArchiveUpgrade(id)
}

// GetAllArchiveUpgrades returns all ArchiveUpgrades
func (service *Service) GetAllArchiveUpgrades() ([]*models.ArchiveUpgrade, error) {
	return service.db.GetAllArchiveUpgrades()
}

// GetAllArchiveUpgradesByUserID returns all ArchiveUpgrades for a User ID
func (service *Service) GetAllArchiveUpgradesByUserID(userID uint) ([]*models.ArchiveUpgrade, error) {
	return service.db.GetAllArchiveUpgradesByUserID(userID)
}
