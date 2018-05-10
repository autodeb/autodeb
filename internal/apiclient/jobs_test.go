package apiclient_test

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnqueueNextJobNoJob(t *testing.T) {
	apiClientTest := setupTest(t)

	apiClientTest.FakeHTTPClient.queueResponse(
		&http.Response{
			StatusCode: http.StatusNoContent,
			Body:       ioutil.NopCloser(strings.NewReader("")),
		},
	)

	job, err := apiClientTest.APIClient.UnqueueNextJob()

	assert.Nil(t, job)
	assert.NoError(t, err)
}
