package worker

import (
	"time"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/worker/jobrunner"
)

func (w *Worker) start() {
	workerQueue := make(chan chan *models.Job)

	// Start job runners
	runnerCount := 1
	for i := 0; i < runnerCount; i++ {
		jobRunner := jobrunner.New(
			workerQueue,
			w.apiClient,
			w.workingDirectory,
			w.logger,
		)
		w.logger.Printf("Starting runner #%d:\n", runnerCount)
		go jobRunner.Start()
	}

	for {
		select {
		case jobs := <-workerQueue: // Wait until a worker asks for a job
			// Loop until we are able to give a job to the runner
			for {
				// Try to get fetch a job
				job, err := w.apiClient.UnqueueNextJob()
				if err != nil {
					w.logger.Printf("Error: could not obtain new job: %v\n", err)
					time.Sleep(10 * time.Second)
					continue
				}
				if job == nil {
					w.logger.Printf("No job available.\n")
					time.Sleep(10 * time.Second)
					continue
				}

				// Give the job to the worker
				w.logger.Printf("Obtained job: %+v\n", job)
				jobs <- job
				break
			}
		}
	}
}
