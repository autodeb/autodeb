package database

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// CreateJob will create a job
func (db *Database) CreateJob(jobType models.JobType, uploadID uint) (*models.Job, error) {
	upload := &models.Job{
		Type:     jobType,
		UploadID: uploadID,
		Status:   models.JobStatusQueued,
	}

	if err := db.gormDB.Create(upload).Error; err != nil {
		return nil, err
	}

	return upload, nil
}

// GetAllJobs returns all jobs
func (db *Database) GetAllJobs() ([]*models.Job, error) {
	var jobs []*models.Job

	if err := db.gormDB.Model(&models.Job{}).Find(&jobs).Error; err != nil {
		return nil, err
	}

	return jobs, nil
}
