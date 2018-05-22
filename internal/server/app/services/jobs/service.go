package jobs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

//Service manages jobs
type Service struct {
	db     *database.Database
	dataFS filesystem.FS
}

//New creates a jobs service
func New(db *database.Database, dataFS filesystem.FS) *Service {
	service := &Service{
		db:     db,
		dataFS: dataFS,
	}
	return service
}

// JobsDirectory contains saved data for jobs such as logs
func (service *Service) JobsDirectory() string {
	return "/jobs"
}

// GetAllJobs returns all jobs
func (service *Service) GetAllJobs() ([]*models.Job, error) {
	return service.db.GetAllJobs()
}

// UnqueueNextJob returns the next job and marks it as assigned
func (service *Service) UnqueueNextJob() (*models.Job, error) {
	job, err := service.db.GetNextJob()
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, nil
	}

	job.Status = models.JobStatusAssigned
	err = service.db.UpdateJob(job)
	if err != nil {
		return nil, err
	}

	return job, err
}

// GetJob returns the job with the given id
func (service *Service) GetJob(id uint) (*models.Job, error) {
	job, err := service.db.GetJob(id)
	if err != nil {
		return nil, err
	}
	return job, nil
}

// UpdateJob will update a job
func (service *Service) UpdateJob(job *models.Job) error {
	return service.db.UpdateJob(job)
}

// GetJobLog returns the log of a job
func (service *Service) GetJobLog(jobID uint) (io.ReadCloser, error) {
	logPath := filepath.Join(
		service.JobsDirectory(),
		fmt.Sprint(jobID),
		"log.txt",
	)

	file, err := service.dataFS.Open(logPath)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return file, nil
}

// SaveJobLog will save logs for a job
func (service *Service) SaveJobLog(jobID uint, content io.Reader) error {
	jobDirectory := filepath.Join(
		service.JobsDirectory(),
		fmt.Sprint(jobID),
	)

	if err := service.dataFS.Mkdir(jobDirectory, 0744); err != nil {
		return err
	}

	logFilePath := filepath.Join(jobDirectory, "log.txt")

	logFile, err := service.dataFS.Create(logFilePath)
	if err != nil {
		return err
	}
	defer logFile.Close()

	if _, err := io.Copy(logFile, content); err != nil {
		return err
	}

	return nil
}
