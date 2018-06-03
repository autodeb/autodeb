package apiclient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// GetArtifact will retrieve an artifact
func (c *APIClient) GetArtifact(artifactID uint) (*models.Artifact, error) {
	var artifact models.Artifact

	response, err := c.getJSON(
		fmt.Sprintf("/api/artifacts/%d", artifactID),
		&artifact,
	)

	if response != nil && response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &artifact, nil
}

// GetArtifactContent will retrieve an artifact's content
func (c *APIClient) GetArtifactContent(artifactID uint) (io.Reader, error) {
	response, body, err := c.get(
		fmt.Sprintf("/api/artifacts/%d/content", artifactID),
	)

	if response != nil && response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(body), nil
}

// GetJobArtifacts retrieves all artifacts associated with a job
func (c *APIClient) GetJobArtifacts(jobID uint) ([]*models.Artifact, error) {
	var jobArtifacts []*models.Artifact

	response, err := c.getJSON(
		fmt.Sprintf("/api/jobs/%d/artifacts", jobID),
		&jobArtifacts,
	)

	if response != nil && response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return jobArtifacts, nil
}

// SubmitJobArtifact will submit a job artifact
func (c *APIClient) SubmitJobArtifact(jobID uint, filename string, content io.Reader) (*models.Artifact, error) {
	var artifact models.Artifact

	response, err := c.postJSON(
		fmt.Sprintf("/api/jobs/%d/artifacts/%s", jobID, filename),
		content,
		&artifact,
	)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, errors.Errorf("Unexpected status code %v", response.Status)
	}

	return &artifact, nil
}
