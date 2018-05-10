// Package worker contains the core of the autodeb worker. It creates all
// dependencies injects them at the right place.
package worker

// Worker is the autodeb worker. It retrieves jobs from the main
// server and executes them
type Worker struct {
}

// New creates a Worker
func New(cfg *Config) (*Worker, error) {
	worker := Worker{}

	return &worker, nil
}

// Close will shutdown the worker
func (srv *Worker) Close() error {
	return nil
}
