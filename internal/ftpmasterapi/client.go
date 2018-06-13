package ftpmasterapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
)

const (
	ftpMasterAPIUrl = "https://api.ftp-master.debian.org"
)

// Client for the ftpmasters api
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new Client
func NewClient(httpClient *http.Client) *Client {
	client := &Client{
		httpClient: httpClient,
	}
	return client
}

//DSC api object as returned by dsc_in_suite
type DSC struct {
	Component string `json:"component"`
	Filename  string `json:"filename"`
}

//GetDSCSInSuite returns a list of DSCs matching pkg in distribution
func (client *Client) GetDSCSInSuite(pkg, distribution string) ([]*DSC, error) {
	url := fmt.Sprintf(
		"%s/dsc_in_suite/%s/%s",
		ftpMasterAPIUrl,
		distribution,
		pkg,
	)

	resp, err := client.httpClient.Get(url)
	if err != nil {
		return nil, errors.WithMessage(err, "get failed")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected HTTP status code: got %d", resp.StatusCode)
	}

	var dscs []*DSC
	if err := json.NewDecoder(resp.Body).Decode(&dscs); err != nil {
		return nil, err
	}

	return dscs, nil

}

// Source api object as returned by source_in_suite
type Source struct {
	Source  string `json:"source"`
	Version string `json:"version"`
}

// GetSourcesInSuite returns all source packages in a suite
func (client *Client) GetSourcesInSuite(distribution string) ([]*Source, error) {
	url := fmt.Sprintf(
		"%s/sources_in_suite/%s",
		ftpMasterAPIUrl,
		distribution,
	)

	resp, err := client.httpClient.Get(url)
	if err != nil {
		return nil, errors.WithMessage(err, "get failed")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected HTTP status code: got %d", resp.StatusCode)
	}

	var sources []*Source
	if err := json.NewDecoder(resp.Body).Decode(&sources); err != nil {
		return nil, err
	}

	return sources, nil
}
