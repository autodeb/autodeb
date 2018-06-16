package database

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"

	"github.com/jinzhu/gorm"
)

// CreateJob will create a job
func (db *Database) CreateJob(jobType models.JobType, input string, parentType models.JobParentType, parentID uint) (*models.Job, error) {
	job := &models.Job{
		Type:       jobType,
		Input:      input,
		Status:     models.JobStatusQueued,
		ParentID:   parentID,
		ParentType: parentType,
	}

	if err := db.gormDB.Create(job).Error; err != nil {
		return nil, err
	}

	return job, nil
}

// GetAllJobs returns all jobs
func (db *Database) GetAllJobs() ([]*models.Job, error) {
	var jobs []*models.Job

	if err := db.gormDB.Model(&models.Job{}).Find(&jobs).Error; err != nil {
		return nil, err
	}

	return jobs, nil
}

// ChangeJobStatus will change a job's status. This is not
// idempotent and will cause an error if the status was not modified.
func (db *Database) ChangeJobStatus(jobID uint, newStatus models.JobStatus) error {

	query := db.gormDB.Model(
		&models.Job{},
	).Where(
		&models.Job{
			ID: jobID,
		},
	).Not(
		&models.Job{
			Status: newStatus,
		},
	).Update(
		&models.Job{
			Status: newStatus,
		},
	)

	if err := query.Error; err != nil {
		return err
	}

	if query.RowsAffected < 1 {
		return errors.Errorf("could not update job id %d to status %s", jobID, newStatus)
	}

	return nil
}

// GetAllJobsByUploadIDStatuses returns all jobs that match the given id and statuses
func (db *Database) GetAllJobsByUploadIDStatuses(uploadID uint, statuses ...models.JobStatus) ([]*models.Job, error) {
	var jobs []*models.Job

	query := db.gormDB.Model(
		&models.Job{},
	)

	// The first status is the first where clause
	if len(statuses) > 0 {
		status := statuses[0]
		statuses = statuses[0:]

		query = query.Where(
			&models.Job{
				ParentID:   uploadID,
				ParentType: models.JobParentTypeUpload,
				Status:     status,
			},
		)
	}

	// All the other statuses are in Or clauses
	for _, status := range statuses[0:] {
		query = query.Or(
			&models.Job{
				ParentID:   uploadID,
				ParentType: models.JobParentTypeUpload,
				Status:     status,
			},
		)
	}

	if err := query.Find(&jobs).Error; err != nil {
		return nil, err
	}

	return jobs, nil
}

// GetAllJobsByUploadID returns all jobs for an upload
func (db *Database) GetAllJobsByUploadID(uploadID uint) ([]*models.Job, error) {
	var jobs []*models.Job

	query := db.gormDB.Model(
		&models.Job{},
	).Where(
		&models.Job{
			ParentType: models.JobParentTypeUpload,
			ParentID:   uploadID,
		},
	)

	if err := query.Find(&jobs).Error; err != nil {
		return nil, err
	}

	return jobs, nil
}

// GetAllJobsByArchiveUpgradeID returns all jobs for an upload
func (db *Database) GetAllJobsByArchiveUpgradeID(id uint) ([]*models.Job, error) {
	var jobs []*models.Job

	query := db.gormDB.Model(
		&models.Job{},
	).Where(
		&models.Job{
			ParentType: models.JobParentTypeArchiveUpgrade,
			ParentID:   id,
		},
	)

	if err := query.Find(&jobs).Error; err != nil {
		return nil, err
	}

	return jobs, nil
}

// GetJob returns the Job with the given id
func (db *Database) GetJob(id uint) (*models.Job, error) {
	var job models.Job

	query := db.gormDB.Where(
		&models.Job{
			ID: id,
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
