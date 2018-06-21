package jobs

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// CreateArchiveUpgrade creates a new ArchiveUpgrade
func (service *Service) CreateArchiveUpgrade(userID uint, packageCount int) (*models.ArchiveUpgrade, error) {
	archiveUpgrade, err := service.db.CreateArchiveUpgrade(userID, packageCount)
	if err != nil {
		return nil, err
	}

	if _, err := service.db.CreateJob(
		models.JobTypeSetupArchiveUpgrade,
		"",
		0,
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

// GetAllArchiveUpgradesPageLimit returns all ArchiveUpgrades with pagination
func (service *Service) GetAllArchiveUpgradesPageLimit(page, limit int) ([]*models.ArchiveUpgrade, error) {
	return service.db.GetAllArchiveUpgradesPageLimit(page, limit)
}

// GetAllArchiveUpgradesByUserID returns all ArchiveUpgrades for a User ID
func (service *Service) GetAllArchiveUpgradesByUserID(userID uint) ([]*models.ArchiveUpgrade, error) {
	return service.db.GetAllArchiveUpgradesByUserID(userID)
}

// GetAllJobsByArchiveUpgradeID returns all jobs for an ArchiveUpgrade
func (service *Service) GetAllJobsByArchiveUpgradeID(id uint) ([]*models.Job, error) {
	return service.db.GetAllJobsByArchiveUpgradeID(id)
}

// GetAllJobsByArchiveUpgradeIDPageLimit returns all jobs for an ArchiveUpgrade
func (service *Service) GetAllJobsByArchiveUpgradeIDPageLimit(id uint, page, limit int) ([]*models.Job, error) {
	return service.db.GetAllJobsByArchiveUpgradeIDPageLimit(id, page, limit)
}

// GetArchiveUpgradeSuccessfulBuilds returns all successful builds for an ArchiveUpgrade.
// Successful builds are build that have passed all tests.
func (service *Service) GetArchiveUpgradeSuccessfulBuilds(archiveUpgradeID uint) ([]*models.Job, error) {
	jobs, err := service.GetAllJobsByArchiveUpgradeID(archiveUpgradeID)
	if err != nil {
		return nil, err
	}

	// Create a [jobID]*models.Job map
	jobMap := make(map[uint]*models.Job)
	for _, job := range jobs {
		jobMap[job.ID] = job
	}

	// Find successful autopkgtest jobs and add their corresponding build jobs to the list
	var successfulBuilds []*models.Job
	for _, job := range jobs {
		if job.Type == models.JobTypeAutopkgtest && job.Status == models.JobStatusSuccess {
			successfulBuilds = append(successfulBuilds, jobMap[job.BuildJobID])
		}
	}

	return successfulBuilds, nil
}
