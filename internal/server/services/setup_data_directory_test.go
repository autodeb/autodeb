package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/services/servicestest"
)

func TestSetupDataDirectory(t *testing.T) {
	servicesTest := servicestest.SetupTest(t)

	fs := servicesTest.DataFS

	_, err := fs.Stat("uploads")
	assert.NoError(t, err, "the directory should have been created")

	_, err = fs.Stat("jobs")
	assert.NoError(t, err, "the directory should have been created")

	_, err = fs.Stat("artifacts")
	assert.NoError(t, err, "the directory should have been created")
}
