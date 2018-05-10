// Package worker contains the core of the autodeb worker. It creates all
// dependencies injects them at the right place.
package worker

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"salsa.debian.org/autodeb-team/autodeb/internal/apiclient"
)

// Worker is the autodeb worker. It retrieves jobs from the main
// server and executes them
type Worker struct {
	apiClient        *apiclient.APIClient
	workingDirectory string
	writerOutput     io.Writer
	writerError      io.Writer
}

// New creates a Worker
func New(cfg *Config) (*Worker, error) {

	// Check that all fields are present
	if cfg.ServerURL == "" {
		return nil, errors.New("ServerURL is empty")
	}
	if cfg.WorkingDirectory == "" {
		return nil, errors.New("WorkingDirectory is empty")
	}
	if cfg.WriterOutput == nil {
		return nil, errors.New("WriterOutput is nil")
	}
	if cfg.WriterError == nil {
		return nil, errors.New("WriterError is nil")
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

	worker := Worker{
		apiClient:        apiClient,
		workingDirectory: workingDirectory,
		writerOutput:     cfg.WriterOutput,
		writerError:      cfg.WriterError,
	}

	go worker.run()

	return &worker, nil
}

// Close will shutdown the worker
func (w *Worker) Close() error {
	return nil
}
