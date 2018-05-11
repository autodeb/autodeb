package worker

import (
	"fmt"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func (w *Worker) submitFailure(job *models.Job, err error) {
	// TODO: do something with the error
	w.submitResult(job, models.JobStatusFailed)
}

func (w *Worker) submitSuccess(job *models.Job) {
	w.submitResult(job, models.JobStatusSuccess)
}

func (w *Worker) submitResult(job *models.Job, status models.JobStatus) {
	if err := w.apiClient.SetJobStatus(job.ID, status); err != nil {
		fmt.Fprintf(
			w.writerOutput,
			"Could not set job %d to status %s", job.ID, status,
		)
	}
	return
}
