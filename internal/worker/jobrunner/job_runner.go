package jobrunner

import (
	"context"
	"fmt"

	"salsa.debian.org/autodeb-team/autodeb/internal/apiclient"
	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
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
	// Setup the job directory
	jobDirectory, err := jobRunner.setupJobDirectory(job)
	if err != nil {
		jobRunner.setJobStatus(job, models.JobStatusQueued)
		jobRunner.logger.Errorf("failed job directory setup: %+v", err)
		return
	}
	defer jobDirectory.Close()

	// Create a cancelable context for the job
	ctx, cancelCtx := jobRunner.getJobContext()
	defer cancelCtx()

	// Execute the job
	jobError := jobRunner.execJob(ctx, job, jobDirectory)

	// Include the job error at the end of the log
	if jobError != nil {
		fmt.Fprintf(jobDirectory.logFile, "\nError: %+v", jobError)
	}

	// If we canceled the job, requeue
	select {
	case <-ctx.Done():
		jobRunner.setJobStatus(job, models.JobStatusQueued)
		return
	default:
		// Else, continue
	}

	// Submit the job result
	jobDirectory.logFile.Seek(0, 0)
	jobRunner.submitJobResult(job, jobError, jobDirectory.logFile, jobDirectory.artifactsDirectory)
}

// getJobContext returns a context that will be canceled if the JobRunner
// needs to quit. It is the caller's responsibility to cancel the context after
// using it
func (jobRunner *JobRunner) getJobContext() (context.Context, context.CancelFunc) {
	ctx, cancelFunc := context.WithCancel(context.Background())

	// Cancel the context if we should quit
	go func() {
		select {
		case <-jobRunner.quit:
			jobRunner.logger.Infof("Canceling the job context")
			cancelFunc()
		case <-ctx.Done():
		}
	}()

	return ctx, cancelFunc
}

func (jobRunner *JobRunner) execJob(ctx context.Context, job *models.Job, jobDirectory *jobDirectory) error {
	switch job.Type {
	case models.JobTypeBuildUpload:
		return jobRunner.execBuildUpload(
			ctx,
			job,
			jobDirectory.workingDirectory,
			jobDirectory.artifactsDirectory,
			jobDirectory.logFile,
		)
	case models.JobTypeAutopkgtest:
		return jobRunner.execAutopkgtest(
			ctx,
			job,
			jobDirectory.workingDirectory,
			jobDirectory.artifactsDirectory,
			jobDirectory.logFile,
		)
	case models.JobTypeForwardUpload:
		return jobRunner.execForwardUpload(
			ctx,
			job,
			jobDirectory.workingDirectory,
			jobDirectory.artifactsDirectory,
			jobDirectory.logFile,
		)
	case models.JobTypeSetupArchiveUpgrade:
		return jobRunner.execSetupArchiveUpgrade(
			ctx,
			job,
			jobDirectory.workingDirectory,
			jobDirectory.artifactsDirectory,
			jobDirectory.logFile,
		)
	case models.JobTypePackageUpgrade:
		return jobRunner.execPackageUpgrade(
			ctx,
			job,
			jobDirectory.workingDirectory,
			jobDirectory.artifactsDirectory,
			jobDirectory.logFile,
		)
	case models.JobTypeCreateArchiveUpgradeRepository:
		return jobRunner.execCreateArchiveUpgradeRepository(
			ctx,
			job,
			jobDirectory.workingDirectory,
			jobDirectory.artifactsDirectory,
			jobDirectory.logFile,
		)
	default:
		jobRunner.logger.Errorf("Unknown job type: %s", job.Type)
		return errors.Errorf("unknown job type: %s", job.Type)
	}
}

// Shutdown stops the JobRunner gracefully, interrupting all running jobs but
// requeueing them on the master node.
func (jobRunner *JobRunner) Shutdown() {
	close(jobRunner.quit)
	<-jobRunner.done
}
