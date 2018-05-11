package apiclient_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"salsa.debian.org/autodeb-team/autodeb/internal/apiclient"
)

type apiClientTest struct {
	APIClient      *apiclient.APIClient
	FakeHTTPClient *FakeHTTPClient
}

func setupTest(t *testing.T) *apiClientTest {
	fakeHTTPClient := &FakeHTTPClient{}

	apiClient, err := apiclient.New(
		"https://auto.debian.net:8080",
		fakeHTTPClient,
	)
	require.NoError(t, err)

	apiClientTest := &apiClientTest{
		APIClient:      apiClient,
		FakeHTTPClient: fakeHTTPClient,
	}

	return apiClientTest
}

type FakeHTTPClient struct {
	requests  []*http.Request
	responses []*http.Response
}

func (c *FakeHTTPClient) queueResponse(response *http.Response) {
	c.responses = append(c.responses, response)
}

func (c *FakeHTTPClient) Do(request *http.Request) (*http.Response, error) {
	// record the request
	c.requests = append(c.requests, request)

	// pop the next response
	response := c.responses[0]
	c.responses = c.responses[1:]

	return response, nil
}
