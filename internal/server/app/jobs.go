package app

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// GetAllJobs returns all jobs
func (app *App) GetAllJobs() ([]*models.Job, error) {
	return app.db.GetAllJobs()
}

// UnqueueNextJob returns the next job and marks it as assigned
func (app *App) UnqueueNextJob() (*models.Job, error) {
	job, err := app.db.GetNextJob()
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, nil
	}

	job.Status = models.JobStatusAssigned
	err = app.db.UpdateJob(job)
	if err != nil {
		return nil, err
	}

	return job, err
}
