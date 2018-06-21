package apiclient

import (
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/aptly"
)

//aptlyHTTPClient will forward the aptly API client's do calls to the autodeb
//client's do method.
type aptlyHTTPClient struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func (c *aptlyHTTPClient) Do(r *http.Request) (*http.Response, error) {
	return c.DoFunc(r)
}

//Aptly returns a client to the aptly api
func (c *APIClient) Aptly() *aptly.APIClient {

	aptlyAPIClient := aptly.NewClient(
		c.absoluteURL("/aptly"),
		&aptlyHTTPClient{
			DoFunc: c.do,
		},
	)

	return aptlyAPIClient
}
