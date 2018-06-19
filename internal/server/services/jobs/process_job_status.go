package jobs

import (
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

	// If the job has failed, there is nothing to do, no matter the parent type.
	if job.Status == models.JobStatusFailed {
		return nil
	}

	switch job.ParentType {
	case models.JobParentTypeUpload:
		return service.processUploadJobStatus(job)
	case models.JobParentTypeArchiveUpgrade:
		return service.processArchiveUpgradeJobStatus(job)
	default:
		return nil
	}

}

func (service *Service) processArchiveUpgradeJobStatus(job *models.Job) error {
	// If this is a package upgrade job, create a corresponding autopkgtest job
	// and stop here.
	if job.Type == models.JobTypePackageUpgrade {
		if _, err := service.CreateAutopkgtestJobFromBuildJob(job); err != nil {
			return err
		}
	}

	// If this is a CreateArchiveUpgradeRepositoryJob, stop here.
	if job.Type == models.JobTypeCreateArchiveUpgradeRepository {
		return nil
	}

	// Was that the last expected job for this archive upgrade?
	if jobs, err := service.GetAllUncompletedJobsByArchiveUpgradeID(job.ParentID); err != nil {
		return err
	} else if len(jobs) > 0 {
		return nil
	}

	// Create a repository from the archive upgrade
	if _, err := service.CreateArchiveUpgradeRepositoryJob(job.ParentID); err != nil {
		return err
	}

	return nil
}

func (service *Service) processUploadJobStatus(job *models.Job) error {
	// If this is a forward job, there is nothing to do, stop here.
	if job.Type == models.JobTypeForwardUpload {
		return nil
	}

	// Retrieve the corresponding upload.
	upload, err := service.db.GetUpload(job.ParentID)
	if err != nil {
		return errors.WithMessage(err, "could not find corresponding upload")
	}

	// If this is a build job and autopkgtest is enabled,
	// create a corresponding autopkgtest job and stop here.
	if job.Type == models.JobTypeBuildUpload && upload.Autopkgtest == true {
		_, err := service.CreateAutopkgtestJobFromBuildJob(job)
		return err
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
