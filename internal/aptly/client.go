package aptly

import (
	"io"
	"net/http"
	"net/url"
	"path"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
)

//APIClient is a client to the aptly REST API
type APIClient struct {
	httpClient HTTPClient
	apiURL     *url.URL
}

//HTTPClient is needed for the APIClient to perform requests. Typically, it
//would be an &http.Client.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// NewClient creates a new APIClient.
//  apiURL should contain "/api". For example: "http://localhost:8080/api"
func NewClient(apiURL *url.URL, httpClient HTTPClient) *APIClient {
	apiClient := &APIClient{
		httpClient: httpClient,
		apiURL:     apiURL,
	}
	return apiClient
}

// absoluteURL will append path to the api url
func (client *APIClient) absoluteURL(relativePath string) *url.URL {
	relativeURL := &url.URL{
		Path: path.Join(client.apiURL.Path, relativePath),
	}
	absoluteURL := client.apiURL.ResolveReference(relativeURL)
	return absoluteURL
}

// do makes an http request
func (client *APIClient) do(method, path, contentType string, body io.Reader) (*http.Response, error) {
	// Retrieve the absolute URL
	absoluteURL := client.absoluteURL(path)

	// Create the request
	request, err := http.NewRequest(method, absoluteURL.String(), body)
	if err != nil {
		return nil, errors.WithMessagef(err, "could not create http request to url %s", absoluteURL.String())
	}

	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	} else {
		request.Header.Set("Content-Type", "application/json")
	}

	// Send the request
	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
