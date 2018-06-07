package jobs

import (
	"path/filepath"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// ProcessJobStatus will change the status of a job and proceed
// with creating the follow-up jobs depending on the related
// upload configuration
func (service *Service) ProcessJobStatus(jobID uint, status models.JobStatus) error {
	// Retrieve the job
	job, err := service.GetJob(jobID)
	if err != nil {
		return err
	}

	// Set the new status
	job.Status = status
	if err := service.db.UpdateJob(job); err != nil {
		return err
	}

	// If the job has failed, there is nothing to do, stop here.
	if job.Status == models.JobStatusFailed {
		return nil
	}

	// If this is a forward job, there is nothing to do, stop here.
	if job.Type == models.JobTypeForward {
		return nil
	}

	// Retrieve the corresponding upload.
	upload, err := service.db.GetUpload(job.UploadID)
	if err != nil {
		return errors.WithMessage(err, "could not find corresponding upload")
	}

	// If this is a build job and autopkgtest is enabled,
	// create corresponding autopkgtest jobs and stop here.
	if job.Type == models.JobTypeBuild && upload.Autopkgtest == true {
		return service.createAutopkgtestJobFromBuildJob(job, upload)
	}

	// The next step can only be to forward the upload. Don't
	// bother to continue if it was not requested.
	if upload.Forward == false {
		return nil
	}

	// Is this the last expected result? If not, stop here.
	if jobs, err := service.GetAllUncompletedJobsByUploadID(upload.ID); err != nil {
		return err
	} else if len(jobs) > 0 {
		return nil
	}

	// Were there any failing jobs? If yes, stop here.
	if jobs, err := service.GetAllFailedJobsByUploadID(upload.ID); err != nil {
		return err
	} else if len(jobs) > 0 {
		return nil
	}

	// Forward the upload.
	if _, err := service.CreateForwardJob(upload.ID); err != nil {
		return err
	}

	return nil
}

// createAutopkgtestJobFromBuildJob will create an autopkgtest job for
// every binary package produced by a build job
func (service *Service) createAutopkgtestJobFromBuildJob(job *models.Job, upload *models.Upload) error {
	// Get the job artifacts
	artifacts, err := service.artifactsService.GetAllArtifactsByJobID(job.ID)
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
