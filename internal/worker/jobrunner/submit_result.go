package jobrunner

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func (jobRunner *JobRunner) submitJobResult(job *models.Job, jobError error, jobLog io.Reader, artifactsDirectory string) {
	// Submit the log
	jobRunner.submitJobLog(job, jobLog)

	// Submit the artifacts
	jobRunner.submitJobArtifacts(job, artifactsDirectory)

	// Set the job status, only after all of this has been completed
	jobStatus := models.JobStatusSuccess
	if jobError != nil {
		jobStatus = models.JobStatusFailed
	}
	jobRunner.setJobStatus(job, jobStatus)
}

func (jobRunner *JobRunner) submitJobArtifacts(job *models.Job, artifactsDirectory string) {
	files, err := ioutil.ReadDir(artifactsDirectory)
	if err != nil {
		jobRunner.logger.Errorf("Could not read artifacts directory: %+v", err)
		return
	}

	for _, file := range files {
		artifact, err := os.Open(
			filepath.Join(
				artifactsDirectory,
				file.Name(),
			),
		)
		if err != nil {
			jobRunner.logger.Errorf("Could not open artifact: %+v", err)
			return
		}
		defer artifact.Close()

		if _, err := jobRunner.apiClient.SubmitJobArtifact(job.ID, file.Name(), artifact); err != nil {
			jobRunner.logger.Errorf("Could not submit job artifact: %+v", err)
			return
		}
	}
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
