package apiclient

import (
	"fmt"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// GetArchiveBackport will retrieve an ArchiveBackport by ID
func (c *APIClient) GetArchiveBackport(id uint) (*models.ArchiveBackport, error) {
	var archiveBackport models.ArchiveBackport

	response, err := c.getJSON(
		fmt.Sprintf("/api/archive-backports/%d", id),
		&archiveBackport,
	)

	if response != nil && response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &archiveBackport, nil
}
