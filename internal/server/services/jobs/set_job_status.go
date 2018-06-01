package jobs

import (
	"path/filepath"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// SetJobStatus will change the status of a job
func (service *Service) SetJobStatus(jobID uint, status models.JobStatus) error {
	job, err := service.GetJob(jobID)
	if err != nil {
		return err
	}
	job.Status = status
	if err := service.db.UpdateJob(job); err != nil {
		return err
	}

	// Process the new status
	switch job.Type {
	case models.JobTypeBuild:
		return service.processBuildJobStatus(job, status)
	default:
		// Nothing to do.
		return nil
	}

}

func (service *Service) processBuildJobStatus(job *models.Job, status models.JobStatus) error {
	// Nothing to do if the job has not succeeded
	if status != models.JobStatusSuccess {
		return nil
	}

	// Retrieve the corresponding upload.
	upload, err := service.db.GetUpload(job.UploadID)
	if err != nil {
		return errors.WithMessage(err, "could not find corresponding upload")
	}

	if upload.Autopkgtest {
		return service.createAutopkgtestJobFromBuildJob(job, upload)
	}

	// TODO: forward upload?

	return nil
}

func (service *Service) createAutopkgtestJobFromBuildJob(job *models.Job, upload *models.Upload) error {
	// Get the job artifacts
	artifacts, err := service.GetAllJobArtifactsByJobID(job.ID)
	if err != nil {
		return errors.WithMessage(err, "could not find corresponding job artifacts")
	}

	// Create autopkgtest jobs for all debs
	for _, artifact := range artifacts {
		if filepath.Ext(artifact.Filename) == ".deb" {
			_, err := service.CreateAutopkgtestJob(upload.ID, artifact.ID)
			if err != nil {
				return errors.WithMessagef(err, "could not create autopkgtest job for upload %d and artifact %d", upload.ID, artifact.ID)
			}
		}
	}

	return nil
}
