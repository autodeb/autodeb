package apiclient

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// GetUpload will retrieve an upload by its id
func (c *APIClient) GetUpload(jobID uint) (*models.Upload, error) {
	var upload models.Upload

	response, err := c.getJSON(
		fmt.Sprintf("/api/uploads/%d", jobID),
		&upload,
	)

	if response != nil && response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &upload, nil
}

// GetUploadDSCURL returns the .dsc URL for a given upload
func (c *APIClient) GetUploadDSCURL(uploadID uint) *url.URL {
	dscURL := c.url(
		fmt.Sprintf("/api/uploads/%d/source.dsc", uploadID),
	)
	return dscURL
}

// GetUploadChangesURL returns the .changes URL for a given upload
func (c *APIClient) GetUploadChangesURL(uploadID uint) *url.URL {
	dscURL := c.url(
		fmt.Sprintf("/api/uploads/%d/package.changes", uploadID),
	)
	return dscURL
}

// GetUploadDSC returns the .dsc of an upload
func (c *APIClient) GetUploadDSC(uploadID uint) (io.Reader, error) {
	response, body, err := c.get(
		c.GetUploadDSCURL(uploadID).EscapedPath(),
	)

	if response != nil && response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(body), nil
}

// GetUploadChanges returns the .changes of an upload
func (c *APIClient) GetUploadChanges(uploadID uint) (io.Reader, error) {
	response, body, err := c.get(
		c.GetUploadChangesURL(uploadID).EscapedPath(),
	)

	if response != nil && response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(body), nil
}

// GetUploadFiles returns all files for an upload
func (c *APIClient) GetUploadFiles(uploadID uint) ([]*models.FileUpload, error) {
	var fileUploads []*models.FileUpload

	response, err := c.getJSON(
		fmt.Sprintf("/api/uploads/%d/files", uploadID),
		&fileUploads,
	)

	if response != nil && response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return fileUploads, nil
}

// GetUploadFile returns an upload's file
func (c *APIClient) GetUploadFile(uploadID uint, filename string) (io.Reader, error) {
	response, body, err := c.get(
		fmt.Sprintf("/api/uploads/%d/%s", uploadID, filename),
	)

	if response != nil && response.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return bytes.NewReader(body), nil
}
