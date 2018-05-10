// Package apiclient implements a client for the autodeb-server REST API
package apiclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

//APIClient is a client for the autodeb-server REST API
type APIClient struct {
	baseURL    *url.URL
	httpClient *http.Client
}

//New creates a new APIClient
func New(serverURL string) (*APIClient, error) {
	baseURL, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}

	apiClient := &APIClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}

	return apiClient, nil
}

func (c *APIClient) post(path string, body io.Reader) (*http.Response, error) {
	return c.do(http.MethodPost, path, body)
}

func (c *APIClient) postJSON(path string, body io.Reader, v interface{}) (*http.Response, error) {
	response, err := c.post(path, body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(v)

	return response, err
}

func (c *APIClient) do(method, path string, body io.Reader) (*http.Response, error) {
	relativeURL := &url.URL{Path: path}
	absoluteURL := c.baseURL.ResolveReference(relativeURL)

	fmt.Printf("%v %v\n", method, absoluteURL.String())

	request, err := http.NewRequest(method, absoluteURL.String(), body)
	if err != nil {
		return nil, err
	}

	return c.httpClient.Do(request)
}
