package jobrunner

import (
	"context"

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

	// The quit channel is closed by Shutdown() to signal that we should quit
	quit chan struct{}

	// The done channel is closed by the main loop to signal that it has exited
	done chan struct{}
}

//New creates a new JobRunner
func New(workerQueue chan chan *models.Job, apiClient *apiclient.APIClient, workingDirectory string, logger log.Logger) *JobRunner {
	jobRunner := &JobRunner{
		workerQueue:      workerQueue,
		apiClient:        apiClient,
		workingDirectory: workingDirectory,
		logger:           logger,

		quit: make(chan struct{}),
		done: make(chan struct{}),
	}
	return jobRunner
}

// Start executing jobs
func (jobRunner *JobRunner) Start() {
	jobs := make(chan *models.Job)

JOB_LOOP:
	for {
		// Wait on the workerQueue or quit
		select {
		case jobRunner.workerQueue <- jobs:
		case <-jobRunner.quit:
			break JOB_LOOP
		}

		// Wait for a job or quit
		select {
		case job := <-jobs:
			jobRunner.execJob(job)
		case <-jobRunner.quit:
			break JOB_LOOP
		}
	}

	jobRunner.logger.Infof("quitting")
	close(jobRunner.done)
}

func (jobRunner *JobRunner) execJob(job *models.Job) {
	jobRunner.logger.Infof("Executing job %v", job)

	// Create a cancelable context for the job
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	// Cancel the context if we should quit
	go func() {
		select {
		case <-jobRunner.quit:
			jobRunner.logger.Infof("Canceling the job context")
			cancelCtx()
		case <-ctx.Done():
		}
	}()

	switch job.Type {
	case models.JobTypeBuild:
		jobRunner.execBuild(ctx, job)
	default:
		jobRunner.logger.Errorf("Unknown job type: %s", job.Type)
	}
}

// Shutdown stops the JobRunner gracefully, interrupting all running jobs but
// requeueing them on the master node.
func (jobRunner *JobRunner) Shutdown() {
	close(jobRunner.quit)
	<-jobRunner.done
}
