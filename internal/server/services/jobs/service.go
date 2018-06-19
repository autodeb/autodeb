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

// GetAllJobsPageLimit returns all jobs with pagination
func (service *Service) GetAllJobsPageLimit(page, limit int) ([]*models.Job, error) {
	return service.db.GetAllJobsPageLimit(page, limit)
}

// GetAllJobsByUploadID returns all jobs for a given upload
func (service *Service) GetAllJobsByUploadID(uploadID uint) ([]*models.Job, error) {
	return service.db.GetAllJobsByUploadID(uploadID)
}

// GetAllUncompletedJobsByUploadID returns all uncompleted jobs for a given upload
func (service *Service) GetAllUncompletedJobsByUploadID(uploadID uint) ([]*models.Job, error) {
	jobs, err := service.db.GetAllJobsByParentAndStatuses(
		models.JobParentTypeUpload,
		uploadID,
		models.JobStatusAssigned,
		models.JobStatusQueued,
	)
	return jobs, err
}

// GetAllUncompletedJobsByArchiveUpgradeID returns all uncompleted jobs for a given archive rebuild
func (service *Service) GetAllUncompletedJobsByArchiveUpgradeID(archiveRebuildID uint) ([]*models.Job, error) {
	jobs, err := service.db.GetAllJobsByParentAndStatuses(
		models.JobParentTypeArchiveUpgrade,
		archiveRebuildID,
		models.JobStatusAssigned,
		models.JobStatusQueued,
	)
	return jobs, err
}

// GetAllFailedJobsByUploadID returns all failed jobs for a given upload
func (service *Service) GetAllFailedJobsByUploadID(uploadID uint) ([]*models.Job, error) {
	jobs, err := service.db.GetAllJobsByParentAndStatuses(
		models.JobParentTypeUpload,
		uploadID,
		models.JobStatusFailed,
	)
	return jobs, err
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

	newStatus := models.JobStatusAssigned

	if err := service.db.ChangeJobStatus(job.ID, newStatus); err != nil {
		return nil, err
	}

	// Set the status for the caller
	job.Status = newStatus

	return job, err
}

//CreateJob creates a new job
func (service *Service) CreateJob(jobType models.JobType, input string, buildJobID uint, parentType models.JobParentType, parentID uint) (*models.Job, error) {
	return service.db.CreateJob(
		jobType,
		input,
		buildJobID,
		parentType,
		parentID,
	)
}

//CreateBuildUploadJob creates a build job
func (service *Service) CreateBuildUploadJob(uploadID uint) (*models.Job, error) {
	return service.db.CreateJob(
		models.JobTypeBuildUpload,
		"",
		0,
		models.JobParentTypeUpload,
		uploadID,
	)
}

//CreateForwardJob creates an upload forward job
func (service *Service) CreateForwardJob(uploadID uint) (*models.Job, error) {
	return service.db.CreateJob(
		models.JobTypeForwardUpload,
		"",
		0,
		models.JobParentTypeUpload,
		uploadID,
	)
}

//CreateAutopkgtestJobFromBuildJob creates an autopkgtest job from the build job
func (service *Service) CreateAutopkgtestJobFromBuildJob(buildJob *models.Job) (*models.Job, error) {
	return service.CreateJob(
		models.JobTypeAutopkgtest,
		"",
		buildJob.ID,
		buildJob.ParentType,
		buildJob.ParentID,
	)
}

//CreateArchiveUpgradeRepositoryJob creates a CreateArchiveUpgradeRepository job
func (service *Service) CreateArchiveUpgradeRepositoryJob(archiveUpgradeID uint) (*models.Job, error) {
	return service.CreateJob(
		models.JobTypeCreateArchiveUpgradeRepository,
		"",
		0,
		models.JobParentTypeArchiveUpgrade,
		archiveUpgradeID,
	)
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
