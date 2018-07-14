package jobs

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// CreateArchiveBackport creates a new ArchiveBackport
func (service *Service) CreateArchiveBackport(userID uint) (*models.ArchiveBackport, error) {
	archiveBackport, err := service.db.CreateArchiveBackport(userID)
	if err != nil {
		return nil, err
	}

	if _, err := service.db.CreateJob(
		models.JobTypeSetupArchiveBackport,
		"",
		0,
		models.JobParentTypeArchiveBackport,
		archiveBackport.ID,
	); err != nil {
		return nil, errors.WithMessagef(err, "could not create archive backport setup job")
	}

	return archiveBackport, err
}

// GetArchiveBackport returns the ArchiveBackport with a matching ID
func (service *Service) GetArchiveBackport(id uint) (*models.ArchiveBackport, error) {
	return service.db.GetArchiveBackport(id)
}

// GetAllArchiveBackports returns all ArchiveBackports
func (service *Service) GetAllArchiveBackports() ([]*models.ArchiveBackport, error) {
	return service.db.GetAllArchiveBackports()
}

// GetAllArchiveBackportsPageLimit returns all ArchiveBackports with pagination
func (service *Service) GetAllArchiveBackportsPageLimit(page, limit int) ([]*models.ArchiveBackport, error) {
	return service.db.GetAllArchiveBackportsPageLimit(page, limit)
}

// GetAllArchiveBackportsByUserID returns all ArchiveBackports for a User ID
func (service *Service) GetAllArchiveBackportsByUserID(userID uint) ([]*models.ArchiveBackport, error) {
	return service.db.GetAllArchiveBackportsByUserID(userID)
}

// GetAllJobsByArchiveBackportID returns all jobs for an ArchiveBackport
func (service *Service) GetAllJobsByArchiveBackportID(id uint) ([]*models.Job, error) {
	return service.db.GetAllJobsByArchiveBackportID(id)
}

// GetAllJobsByArchiveBackportIDPageLimit returns all jobs for an ArchiveBackport
func (service *Service) GetAllJobsByArchiveBackportIDPageLimit(id uint, page, limit int) ([]*models.Job, error) {
	return service.db.GetAllJobsByArchiveBackportIDPageLimit(id, page, limit)
}
