package apiclient

import (
	"fmt"
	"net/url"
)

// GetUploadDSCURL returns the .dsc URL for a given upload
func (c *APIClient) GetUploadDSCURL(uploadID uint) *url.URL {
	dscURL := c.url(
		fmt.Sprintf("/api/uploads/%v/source.dsc", uploadID),
	)
	return dscURL
}
