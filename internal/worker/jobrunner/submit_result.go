package jobrunner

import (
	"io"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func (jobRunner *JobRunner) submitResult(job *models.Job, jobError error, jobLog io.Reader) {
	// Set the job status
	jobStatus := models.JobStatusSuccess
	if jobError != nil {
		jobStatus = models.JobStatusFailed
	}
	jobRunner.setJobStatus(job, jobStatus)

	// Submit the log
	jobRunner.submitJobLog(job, jobLog)
}

func (jobRunner *JobRunner) submitJobLog(job *models.Job, jobLog io.Reader) {
	if err := jobRunner.apiClient.SubmitJobLog(job.ID, jobLog); err != nil {
		jobRunner.logger.Errorf("Could not submit job log: %+v", err)
	}
}

func (jobRunner *JobRunner) setJobStatus(job *models.Job, status models.JobStatus) {
	if err := jobRunner.apiClient.SetJobStatus(job.ID, status); err != nil {
		jobRunner.logger.Errorf("Could not set job %d to status %s: %+v", job.ID, status, err)
	}
}
