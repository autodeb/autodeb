package jobrunner

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/apiclient"
	"salsa.debian.org/autodeb-team/autodeb/internal/log"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

//JobRunner loops indefinitely and runs one job at a time
type JobRunner struct {
	apiClient        *apiclient.APIClient
	workerQueue      chan chan *models.Job
	workingDirectory string
	logger           log.Logger
}

//New creates a new JobRunner
func New(workerQueue chan chan *models.Job, apiClient *apiclient.APIClient, workingDirectory string, logger log.Logger) *JobRunner {
	jobRunner := &JobRunner{
		workerQueue:      workerQueue,
		apiClient:        apiClient,
		workingDirectory: workingDirectory,
		logger:           logger,
	}
	return jobRunner
}

// Start executing jobs
func (jobRunner *JobRunner) Start() {
	jobs := make(chan *models.Job)

	for {
		// Signal that we are ready
		jobRunner.workerQueue <- jobs

		// Wait for a job
		select {
		case job := <-jobs:
			jobRunner.logger.Infof("received job %v", job)
			jobRunner.execJob(job)
		}
	}

}

func (jobRunner *JobRunner) execJob(job *models.Job) {
	switch job.Type {
	case models.JobTypeBuild:
		jobRunner.execBuild(job)
	default:
		jobRunner.logger.Errorf("Unknown job type: %s", job.Type)
	}
}
