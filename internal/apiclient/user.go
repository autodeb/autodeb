package apiclient

import (
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// GetCurrentUser will retrieve the current user
func (c *APIClient) GetCurrentUser() (*models.User, error) {
	var user models.User

	response, err := c.getJSON("/api/user", &user)

	// We are redirected to an auth page
	if response != nil && response.StatusCode == http.StatusSeeOther {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}
