package routertest

import (
	"net/http"
	"net/http/httptest"

	"salsa.debian.org/autodeb-team/autodeb/internal/apiclient"
)

//TestAPIClient embeds an API client and allows to retrieve its last request
type TestAPIClient struct {
	*apiclient.APIClient
	handlerHTTPClient *handlerHTTPClient
}

func newTestAPIClient(handler http.Handler, token string) *TestAPIClient {
	httpClient := &handlerHTTPClient{
		handler: handler,
	}

	apiClient, err := apiclient.New(
		"https://localhost:0871",
		token,
		httpClient,
	)
	if err != nil {
		panic(err)
	}

	testAPIClient := &TestAPIClient{
		APIClient:         apiClient,
		handlerHTTPClient: httpClient,
	}

	return testAPIClient
}

//LastRecorder returns the last response recorder used by the HTTP client
func (c *TestAPIClient) LastRecorder() *httptest.ResponseRecorder {
	return c.handlerHTTPClient.lastRecorder()
}

//LastResponse returns the last response returned to the APIClient
func (c *TestAPIClient) LastResponse() *http.Response {
	return c.LastRecorder().Result()
}

//handlerHTTPClient is an http client that will forward all requests to the
//given handler and keep a copy of all response recorders.
type handlerHTTPClient struct {
	handler   http.Handler
	responses []*httptest.ResponseRecorder
}

func (c *handlerHTTPClient) Do(request *http.Request) (*http.Response, error) {
	//Create a response recorder
	recorder := httptest.NewRecorder()

	//Perform the request
	c.handler.ServeHTTP(recorder, request)

	//Save the response recorder
	c.responses = append(c.responses, recorder)

	return recorder.Result(), nil
}

func (c *handlerHTTPClient) lastRecorder() *httptest.ResponseRecorder {
	return c.responses[len(c.responses)-1]
}
