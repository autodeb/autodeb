package apiclient

import (
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
