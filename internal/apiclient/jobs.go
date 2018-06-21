package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// CreateJob will create a new job
func (c *APIClient) CreateJob(job *models.Job) (*models.Job, error) {
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(job); err != nil {
		return nil, errors.WithMessage(err, "could not encode job")
	}

	var createdJob models.Job

	response, err := c.postJSON(
		fmt.Sprintf("/api/jobs"),
		&body,
		&createdJob,
	)

	if response != nil && response.StatusCode != http.StatusCreated {
		return nil, errors.Errorf("unexpected status code %d", response.StatusCode)
	}

	return &createdJob, err
}

// GetJob will retrieve a job by ID
func (c *APIClient) GetJob(jobID uint) (*models.Job, error) {
	var job models.Job

	response, err := c.getJSON(
		fmt.Sprintf("/api/jobs/%d", jobID),
		&job,
	)

	if response != nil && response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &job, nil
}

// GetJobLogContent will retrieve the content of a job's log
func (c *APIClient) GetJobLogContent(jobID uint) (io.Reader, error) {
	response, body, err := c.get(
		fmt.Sprintf("/api/jobs/%d/log.txt", jobID),
	)

	if response != nil && response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(body), nil
}

// GetJobArtifactContent will retrieve the content of a job's artifact
func (c *APIClient) GetJobArtifactContent(jobID uint, filename string) (io.Reader, error) {
	_, body, err := c.get(
		fmt.Sprintf("/api/jobs/%d/artifacts/%s", jobID, filename),
	)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(body), nil
}

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
	response, _, err := c.post(
		fmt.Sprintf("/api/jobs/%d/status/%d", jobID, status),
		nil,
	)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.Errorf("Unexpected status code %v", response.Status)
	}

	return nil
}

// SubmitJobLog will submit logs for a job
func (c *APIClient) SubmitJobLog(jobID uint, jobLog io.Reader) error {
	response, _, err := c.post(
		fmt.Sprintf("/api/jobs/%d/log", jobID),
		jobLog,
	)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusCreated {
		return errors.Errorf("Unexpected status code %v", response.Status)
	}

	return nil
}
