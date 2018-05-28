// Package apiclient implements a client for the autodeb-server REST API
package apiclient

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
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

func (c *APIClient) post(path string, body io.Reader) (*http.Response, []byte, error) {
	return c.do(http.MethodPost, path, body)
}

func (c *APIClient) get(path string) (*http.Response, []byte, error) {
	return c.do(http.MethodGet, path, nil)
}

func (c *APIClient) postJSON(path string, body io.Reader, v interface{}) (*http.Response, error) {
	response, responseBody, err := c.post(path, body)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(
		bytes.NewReader(responseBody),
	).Decode(v)

	return response, err
}

func (c *APIClient) getJSON(path string, v interface{}) (*http.Response, error) {
	response, body, err := c.get(path)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(
		bytes.NewReader(body),
	).Decode(v)

	return response, err
}

func (c *APIClient) do(method, path string, body io.Reader) (*http.Response, []byte, error) {
	absoluteURL := c.url(path)

	request, err := http.NewRequest(method, absoluteURL.String(), body)
	if err != nil {
		return nil, nil, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}

	return response, b, nil
}
