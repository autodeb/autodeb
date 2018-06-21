package apiclient_test

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAptly(t *testing.T) {
	apiClientTest := setupTest(t)

	apiClientTest.FakeHTTPClient.queueResponse(
		&http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader("{}")),
		},
	)

	_, err := apiClientTest.APIClient.Aptly().CreateRepository(
		"",
		"",
		"",
		"",
	)
	assert.NoError(t, err)

	request := apiClientTest.FakeHTTPClient.requests[0]
	assert.Equal(
		t,
		"/aptly/repos",
		request.URL.Path,
		"/aptly should have been prefixed to the path path",
	)

	assert.Contains(
		t,
		request.Header["Authorization"][0],
		"testtoken",
		"the autodeb apiclient's token should also be appended to aptly requests",
	)
}
