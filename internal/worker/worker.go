// Package worker contains the core of the autodeb worker. It creates all
// dependencies injects them at the right place.
package worker

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"salsa.debian.org/autodeb-team/autodeb/internal/apiclient"
	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/log"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/worker/jobrunner"
)

// Worker is the autodeb worker. It retrieves jobs from the main
// server and passes them to JobRunners that will execute them.
//
// Workers are not safe for use by multiple goroutines.
type Worker struct {
	apiClient        *apiclient.APIClient
	workingDirectory string
	logger           log.Logger

	// This slice contains all JobRunners that were started by this worker.
	jobRunners []*jobrunner.JobRunner

	// This is a queue of JobRunners waiting for a job
	jobRunnerQueue chan chan *models.Job

	// The quit channel is closed by Shutdown() to signal that we should quit
	quit chan struct{}

	// The done channel is closed by the main loop to signal that it has exited
	done chan struct{}
}

// New creates a Worker
func New(cfg *Config, loggingOutput io.Writer) (*Worker, error) {

	// Check that all fields are present
	if cfg.ServerURL == "" {
		return nil, errors.New("ServerURL is empty")
	}
	if cfg.WorkingDirectory == "" {
		return nil, errors.New("WorkingDirectory is empty")
	}
	if cfg.RunnerCount == 0 {
		return nil, errors.New("RunnerCount cannot be 0")
	}

	// Create the apiClient
	apiClient, err := apiclient.New(cfg.ServerURL, &http.Client{})
	if err != nil {
		return nil, err
	}

	// Set workingDirectory to the absolute path
	workingDirectory, err := filepath.Abs(cfg.WorkingDirectory)
	if err != nil {
		return nil, err
	}

	// Create the workingDirectory
	if err := os.MkdirAll(workingDirectory, 0755); err != nil {
		return nil, err
	}

	// Create the logger
	logger := log.New(loggingOutput)
	logger.SetLevel(cfg.LogLevel)

	worker := Worker{
		apiClient:        apiClient,
		workingDirectory: workingDirectory,
		logger:           logger,
		jobRunnerQueue:   make(chan chan *models.Job),
		quit:             make(chan struct{}),
		done:             make(chan struct{}),
	}

	worker.startJobRunners(cfg.RunnerCount)

	go worker.dispatchJobs()

	return &worker, nil
}

func (w *Worker) startJobRunners(count int) {
	for i := 0; i < count; i++ {
		jobRunner := jobrunner.New(
			w.jobRunnerQueue,
			w.apiClient,
			w.workingDirectory,
			w.logger.PrefixLogger(
				fmt.Sprintf("JobRunner#%d", i),
			),
		)

		w.jobRunners = append(w.jobRunners, jobRunner)

		w.logger.Infof("Starting JobRunner#%d", i)

		go jobRunner.Start()
	}
}

func (w *Worker) dispatchJobs() {
DISPATCH_JOBS_LOOP:
	for {
		select {
		case jobs := <-w.jobRunnerQueue: // Wait until a JobRunner asks for a job
			// Loop until we are able to give a job to the runner
			for {
				if dispatched := w.dispatchJob(jobs); dispatched {
					break
				}

				// Could not dispatch a job to the runner. Wait 10 seconds
				// before trying again or quit if we are asked to.
				select {
				case <-w.quit:
					break DISPATCH_JOBS_LOOP
				case <-time.After(10 * time.Second):
					continue
				}
			}
		case <-w.quit:
			break DISPATCH_JOBS_LOOP
		}
	}
	w.logger.Infof("Quitting - no longer dispatching jobs")
	close(w.done)
}

func (w *Worker) dispatchJob(jobs chan *models.Job) bool {
	job, err := w.apiClient.UnqueueNextJob()
	if err != nil {
		w.logger.Errorf("Could not obtain new job: %v", err)
		return false
	}

	if job == nil {
		w.logger.Infof("No job available.")
		return false
	}

	w.logger.Infof("Obtained job: %+v", job)
	jobs <- job

	return true
}

// Shutdown stop the worker gracefully.
func (w *Worker) Shutdown() error {
	// Stop dispatching jobs
	close(w.quit)
	<-w.done

	// Shutdown remaining runners
	var wg sync.WaitGroup
	for _, jobRunner := range w.jobRunners {
		wg.Add(1)
		go func(jr *jobrunner.JobRunner) {
			jr.Shutdown()
			wg.Done()
		}(jobRunner)
	}
	wg.Wait()

	return nil
}
