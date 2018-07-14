package apiclient

import (
	"fmt"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// GetArchiveUpgrade will retrieve an archive upgrade by ID
func (c *APIClient) GetArchiveUpgrade(id uint) (*models.ArchiveUpgrade, error) {
	var archiveUpgrade models.ArchiveUpgrade

	response, err := c.getJSON(
		fmt.Sprintf("/api/archive-upgrades/%d", id),
		&archiveUpgrade,
	)

	if response != nil && response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &archiveUpgrade, nil
}

// GetArchiveUpgradeJobs returns all jobs for an ArchiveUpgrade
func (c *APIClient) GetArchiveUpgradeJobs(id uint) ([]*models.Job, error) {
	var jobs []*models.Job

	response, err := c.getJSON(
		fmt.Sprintf("/api/archive-upgrades/%d/jobs", id),
		&jobs,
	)

	if response != nil && response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return jobs, nil
}
