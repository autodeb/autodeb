package ftpmasterapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
)

const (
	ftpMasterAPIUrl = "https://api.ftp-master.debian.org"
	mirrorURL       = "https://deb.debian.org/debian"
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
		return nil, errors.WithMessage(err, "cannot parse json output")
	}

	return dscs, nil
}

// SHA256SumInArchive is an element of the sha256sum_in_archive query
type SHA256SumInArchive struct {
	SHA256Sum string `json:"sha256sum"`
	Filename  string `json:"filename"`
}

// GetSHA256SumInArchive returns a list of files with matching shasums in the archive
func (client *Client) GetSHA256SumInArchive(sha256sum string) ([]*SHA256SumInArchive, error) {
	url := fmt.Sprintf(
		"%s/sha256sum_in_archive/%s",
		ftpMasterAPIUrl,
		sha256sum,
	)

	resp, err := client.httpClient.Get(url)
	if err != nil {
		return nil, errors.WithMessage(err, "get failed")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected HTTP status code: got %d for url %s", resp.StatusCode, url)
	}

	var files []*SHA256SumInArchive
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, errors.WithMessage(err, "cannot parse json output")
	}

	return files, nil
}

// GetFileBySHA256Sum returns a file in the archive by sha256sum
func (client *Client) GetFileBySHA256Sum(sha256sum string) (io.ReadCloser, error) {
	files, err := client.GetSHA256SumInArchive(sha256sum)
	if err != nil {
		return nil, err
	} else if len(files) < 1 {
		return nil, errors.Errorf("could not find file for sha256sum %s", sha256sum)
	}

	url := fmt.Sprintf(
		"%s/pool/main/%s",
		mirrorURL,
		files[0].Filename,
	)

	resp, err := client.httpClient.Get(url)
	if err != nil {
		return nil, errors.WithMessage(err, "get failed")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected HTTP status code: got %d for url %s", resp.StatusCode, url)
	}

	return resp.Body, nil
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
