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
	httpClient HTTPClient
}

//HTTPClient is needed for the APIClient to perform requests. Typically, it
//would be an &http.Client.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

//New creates a new APIClient
func New(serverURL string, httpClient HTTPClient) (*APIClient, error) {
	baseURL, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}

	apiClient := &APIClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}

	return apiClient, nil
}

func (c *APIClient) url(path string) *url.URL {
	relativeURL := &url.URL{Path: path}
	absoluteURL := c.baseURL.ResolveReference(relativeURL)
	return absoluteURL
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
	absoluteURL := c.url(path)

	fmt.Printf("%v %v\n", method, absoluteURL.String())

	request, err := http.NewRequest(method, absoluteURL.String(), body)
	if err != nil {
		return nil, err
	}

	return c.httpClient.Do(request)
}
