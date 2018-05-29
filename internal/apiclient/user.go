package apiclient

import (
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// GetCurrentUser will retrieve the current user
func (c *APIClient) GetCurrentUser() (*models.User, error) {
	var user models.User

	response, err := c.getJSON("/api/user", &user)

	// We are not authenticated
	if response != nil && response.StatusCode == http.StatusForbidden {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
