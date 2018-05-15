package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

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

// GetJob returns the job with the given id
func (app *App) GetJob(id uint) (*models.Job, error) {
	job, err := app.db.GetJob(id)
	if err != nil {
		return nil, err
	}
	return job, nil
}

// UpdateJob will update a job
func (app *App) UpdateJob(job *models.Job) error {
	return app.db.UpdateJob(job)
}

// GetJobLog returns the log of a job
func (app *App) GetJobLog(jobID uint) (io.ReadCloser, error) {
	logPath := filepath.Join(
		app.JobsDirectory(),
		fmt.Sprint(jobID),
		"log.txt",
	)

	file, err := app.dataFS.Open(logPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return file, nil
}

// SaveJobLog will save logs for a job
func (app *App) SaveJobLog(jobID uint, content io.Reader) error {
	jobDirectory := filepath.Join(
		app.JobsDirectory(),
		fmt.Sprint(jobID),
	)

	if err := app.dataFS.Mkdir(jobDirectory, 0744); err != nil {
		return err
	}

	logFilePath := filepath.Join(jobDirectory, "log.txt")

	logFile, err := app.dataFS.Create(logFilePath)
	if err != nil {
		return err
	}
	defer logFile.Close()

	if _, err := io.Copy(logFile, content); err != nil {
		return err
	}

	return nil
}
