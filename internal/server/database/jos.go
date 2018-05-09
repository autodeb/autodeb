package database

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"

	"github.com/jinzhu/gorm"
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

// UpdateJob will update a job
func (db *Database) UpdateJob(job *models.Job) error {
	err := db.gormDB.Save(job).Error
	return err
}

// GetNextJob will return the next job to run
func (db *Database) GetNextJob() (*models.Job, error) {
	var job models.Job

	query := db.gormDB.Where(
		&models.Job{
			Status: models.JobStatusQueued,
		},
	)

	err := query.First(&job).Error

	if gorm.IsRecordNotFoundError(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &job, nil
}
