package apiclient_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUploadDSCURL(t *testing.T) {
	apiClientTest := setupTest(t)

	dscURL := apiClientTest.APIClient.GetUploadDSCURL(uint(1))

	assert.Equal(t, "/api/uploads/1/source.dsc", dscURL.Path)
}
