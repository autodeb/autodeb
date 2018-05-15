package jobrunner

import (
	"context"
	"fmt"
	"io"
	"os"

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
			jobRunner.logger.Infof("Received job %+v", job)
			jobRunner.setupAndExecJob(job)
		case <-jobRunner.quit:
			break JOB_LOOP
		}
	}

	jobRunner.logger.Infof("quitting")
	close(jobRunner.done)
}

func (jobRunner *JobRunner) setupAndExecJob(job *models.Job) {
	// Setup the job
	workingDirectory, logFile, err := jobRunner.setupJob(job)
	if err != nil {
		jobRunner.setJobStatus(job, models.JobStatusQueued)
		jobRunner.logger.Errorf("failed job setup: %+v", err)
		return
	}
	defer os.RemoveAll(workingDirectory)
	defer logFile.Close()

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

	jobError := jobRunner.execJob(ctx, job, workingDirectory, logFile)

	// If we canceled the job, requeue
	select {
	case <-ctx.Done():
		jobRunner.setJobStatus(job, models.JobStatusQueued)
		return
	default:
		// Else, continue
	}

	// Submit the job result
	logFile.Seek(0, 0)
	jobRunner.submitResult(job, jobError, logFile)
}

func (jobRunner *JobRunner) execJob(ctx context.Context, job *models.Job, workingDirectory string, logFile io.Writer) error {
	switch job.Type {
	case models.JobTypeBuild:
		return jobRunner.execBuild(ctx, job, workingDirectory, logFile)
	default:
		jobRunner.logger.Errorf("Unknown job type: %s", job.Type)
		return fmt.Errorf("unknown job type: %s", job.Type)
	}
}

// Shutdown stops the JobRunner gracefully, interrupting all running jobs but
// requeueing them on the master node.
func (jobRunner *JobRunner) Shutdown() {
	close(jobRunner.quit)
	<-jobRunner.done
}
