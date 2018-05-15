package apiclient

import (
	"fmt"
	"io"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// UnqueueNextJob will return the next job on the queue
func (c *APIClient) UnqueueNextJob() (*models.Job, error) {
	var job models.Job

	response, err := c.postJSON("/api/jobs/next", nil, &job)

	// no job available, don't bother looking at the json decoding error
	if response != nil && response.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &job, nil
}

// SetJobStatus will set the Job Status
func (c *APIClient) SetJobStatus(jobID uint, status models.JobStatus) error {
	response, err := c.post(
		fmt.Sprintf("/api/jobs/%d/status/%d", jobID, status),
		nil,
	)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status code %v", response.Status)
	}

	return nil
}

// SubmitJobLog will submit logs for a job
func (c *APIClient) SubmitJobLog(jobID uint, jobLog io.Reader) error {
	response, err := c.post(
		fmt.Sprintf("/api/jobs/%d/log", jobID),
		jobLog,
	)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status code %v", response.Status)
	}

	return nil
}
