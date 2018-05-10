package cli_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"salsa.debian.org/autodeb-team/autodeb/cmd/autodeb-worker/internal/cli"
	"salsa.debian.org/autodeb-team/autodeb/internal/worker"
)

type cliTest struct {
	outputWriter bytes.Buffer
	errorWriter  bytes.Buffer
}

func (cliTest *cliTest) Parse(args ...string) (*worker.Config, error) {
	return cli.Parse(args, &cliTest.outputWriter, &cliTest.errorWriter)
}

func testSetup() *cliTest {
	cliTest := cliTest{}
	return &cliTest
}

func TestMissingServerAddress(t *testing.T) {
	cliTest := testSetup()

	cfg, err := cliTest.Parse(
		"-server-port", "1",
	)

	assert.Nil(t, cfg)
	assert.EqualError(t, err, "missing argument: server-address")
}

func TestMissingServerPort(t *testing.T) {
	cliTest := testSetup()

	cfg, err := cliTest.Parse(
		"-server-address", "test.com",
	)

	assert.Nil(t, cfg)
	assert.EqualError(t, err, "missing argument: server-port")
}

func TestEmptyArguments(t *testing.T) {
	cliTest := testSetup()

	cfg, err := cliTest.Parse(
		"-server-address", "test.com",
		"-server-port", "1",
	)

	assert.NotNil(t, cfg)
	assert.NoError(t, err)
}
