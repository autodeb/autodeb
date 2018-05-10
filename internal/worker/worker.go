// Package worker contains the core of the autodeb worker. It creates all
// dependencies injects them at the right place.
package worker

import (
	"io"

	"salsa.debian.org/autodeb-team/autodeb/internal/apiclient"
)

// Worker is the autodeb worker. It retrieves jobs from the main
// server and executes them
type Worker struct {
	writerOutput io.Writer
	writerError  io.Writer
	apiClient    *apiclient.APIClient
}

// New creates a Worker
func New(cfg *Config) (*Worker, error) {

	apiClient, err := apiclient.New(cfg.ServerURL)
	if err != nil {
		return nil, err
	}

	worker := Worker{
		apiClient:    apiClient,
		writerOutput: cfg.WriterOutput,
		writerError:  cfg.WriterError,
	}

	go worker.run()

	return &worker, nil
}

// Close will shutdown the worker
func (srv *Worker) Close() error {
	return nil
}
