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
