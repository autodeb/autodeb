package database

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// CreateJobArtifact will create a job artifact
func (db *Database) CreateJobArtifact(jobID uint, filename string) (*models.JobArtifact, error) {
	jobArtifact := &models.JobArtifact{
		JobID:    jobID,
		Filename: filename,
	}

	if err := db.gormDB.Create(jobArtifact).Error; err != nil {
		return nil, err
	}

	return jobArtifact, nil
}

// GetAllJobArtifactsByJobID returns all job artifacts for a job
func (db *Database) GetAllJobArtifactsByJobID(jobID uint) ([]*models.JobArtifact, error) {
	var jobArtifacts []*models.JobArtifact

	query := db.gormDB.Model(
		&models.JobArtifact{},
	).Where(
		&models.JobArtifact{
			JobID: jobID,
		},
	)

	if err := query.Find(&jobArtifacts).Error; err != nil {
		return nil, err
	}

	return jobArtifacts, nil
}
