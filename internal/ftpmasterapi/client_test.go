package ftpmasterapi_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"salsa.debian.org/autodeb-team/autodeb/internal/ftpmasterapi"
)

func setupTest(t *testing.T) *ftpmasterapi.Client {
	client := ftpmasterapi.NewClient(http.DefaultClient)
	return client
}

func TestDSCURL(t *testing.T) {
	client := setupTest(t)

	dsc := &ftpmasterapi.DSC{
		Component: "main",
		Filename:  "i/influxdb/influxdb_1.1.1+dfsg1-4.dsc",
	}

	assert.Equal(
		t,
		"https://deb.debian.org/debian/pool/main/i/influxdb/influxdb_1.1.1+dfsg1-4.dsc",
		client.DSCURL(dsc),
	)
}
