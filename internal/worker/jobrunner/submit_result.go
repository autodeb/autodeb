package jobrunner

import (
	"fmt"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func (jobRunner *JobRunner) submitFailure(job *models.Job, err error) {
	// TODO: do something with the error
	jobRunner.submitResult(job, models.JobStatusFailed)
}

func (jobRunner *JobRunner) submitSuccess(job *models.Job) {
	jobRunner.submitResult(job, models.JobStatusSuccess)
}

func (jobRunner *JobRunner) submitResult(job *models.Job, status models.JobStatus) {
	if err := jobRunner.apiClient.SetJobStatus(job.ID, status); err != nil {
		fmt.Printf("Could not set job %d to status %s", job.ID, status)
	}
}
