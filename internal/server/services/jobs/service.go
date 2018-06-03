package jobs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/services/artifacts"
)

//Service manages jobs
type Service struct {
	db               *database.Database
	fs               filesystem.FS
	artifactsService *artifacts.Service
}

//New creates a jobs service
func New(db *database.Database, artifactsService *artifacts.Service, fs filesystem.FS) *Service {
	service := &Service{
		db:               db,
		artifactsService: artifactsService,
		fs:               fs,
	}
	return service
}

// jobsDirectory contains saved data for jobs such as logs
func (service *Service) jobsDirectory() string {
	return "/"
}

// jobDirectory returns the directory name for a job id
func (service *Service) jobDirectory(jobID uint) string {
	jobDirectory := filepath.Join(
		service.jobsDirectory(),
		fmt.Sprint(jobID),
	)
	return jobDirectory
}

// jobLogPath returns the path of a job's log
func (service *Service) jobLogPath(jobID uint) string {
	jobLogPath := filepath.Join(
		service.jobDirectory(jobID),
		"log.txt",
	)
	return jobLogPath
}

// GetAllJobs returns all jobs
func (service *Service) GetAllJobs() ([]*models.Job, error) {
	return service.db.GetAllJobs()
}

// GetAllJobsByUploadID returns all jobs for a given upload
func (service *Service) GetAllJobsByUploadID(uploadID uint) ([]*models.Job, error) {
	return service.db.GetAllJobsByUploadID(uploadID)
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

// CreateBuildJob creates a build job
func (service *Service) CreateBuildJob(uploadID uint) (*models.Job, error) {
	return service.db.CreateJob(models.JobTypeBuild, uploadID, 0)
}

// CreateAutopkgtestJob creates an autopkgtest job for the provided .deb artifact id
func (service *Service) CreateAutopkgtestJob(uploadID uint, debJobArtifactID uint) (*models.Job, error) {
	return service.db.CreateJob(models.JobTypeAutopkgtest, uploadID, debJobArtifactID)
}

// GetJob returns the job with the given id
func (service *Service) GetJob(id uint) (*models.Job, error) {
	job, err := service.db.GetJob(id)
	if err != nil {
		return nil, err
	}
	return job, nil
}

// GetJobLog returns the log of a job
func (service *Service) GetJobLog(jobID uint) (io.ReadCloser, error) {
	file, err := service.fs.Open(service.jobLogPath(jobID))
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
	if err := service.fs.Mkdir(service.jobDirectory(jobID), 0744); err != nil {
		return err
	}

	logFile, err := service.fs.Create(service.jobLogPath(jobID))
	if err != nil {
		return err
	}
	defer logFile.Close()

	if _, err := io.Copy(logFile, content); err != nil {
		return err
	}

	return nil
}
