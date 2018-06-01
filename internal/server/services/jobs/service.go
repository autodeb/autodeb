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
	db *database.Database
	fs filesystem.FS
}

//New creates a jobs service
func New(db *database.Database, fs filesystem.FS) *Service {
	service := &Service{
		db: db,
		fs: fs,
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

// jobArtifactsDirectory returns the path of a job's artifacts directory
func (service *Service) jobArtifactsDirectory(jobID uint) string {
	jobArtifactDirectory := filepath.Join(
		service.jobDirectory(jobID),
		"artifacts",
	)
	return jobArtifactDirectory
}

// jobArtifactPath returns the path of a job's artifact
func (service *Service) jobArtifactPath(jobID uint, filename string) string {
	// Clean the file name
	_, filename = filepath.Split(filename)
	jobArtifactPath := filepath.Join(
		service.jobArtifactsDirectory(jobID),
		filename,
	)
	return jobArtifactPath
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
	return service.db.CreateJob(models.JobTypeBuild, uploadID)
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

// SaveJobArtifact will save a job artifact
func (service *Service) SaveJobArtifact(jobID uint, filename string, content io.Reader) error {
	if _, err := service.db.CreateJobArtifact(jobID, filename); err != nil {
		return err
	}

	if err := service.fs.MkdirAll(service.jobArtifactsDirectory(jobID), 0744); err != nil {
		return err
	}

	artifact, err := service.fs.Create(service.jobArtifactPath(jobID, filename))
	if err != nil {
		return err
	}
	defer artifact.Close()

	if _, err := io.Copy(artifact, content); err != nil {
		return err
	}

	return nil
}

// GetJobArtifact returns a job artifact
func (service *Service) GetJobArtifact(jobID uint, filename string) (io.ReadCloser, error) {
	file, err := service.fs.Open(service.jobArtifactPath(jobID, filename))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return file, nil
}

// GetAllJobArtifactsByJobID returns a list of all artifacts for a job
func (service *Service) GetAllJobArtifactsByJobID(jobID uint) ([]*models.JobArtifact, error) {
	jobArtifacts, err := service.db.GetAllJobArtifactsByJobID(jobID)
	if err != nil {
		return nil, err
	}
	return jobArtifacts, nil
}
