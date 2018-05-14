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
}

func (cliTest *cliTest) Parse(args ...string) (*worker.Config, error) {
	return cli.Parse(args, &cliTest.outputWriter)
}

func testSetup() *cliTest {
	cliTest := cliTest{}
	return &cliTest
}

func TestMissingServerURL(t *testing.T) {
	cliTest := testSetup()

	cfg, err := cliTest.Parse()

	assert.Nil(t, cfg)
	assert.EqualError(t, err, "missing argument: server-url")
}

func TestEmptyArguments(t *testing.T) {
	cliTest := testSetup()

	cfg, err := cliTest.Parse(
		"-server-url", "hello",
		"",
	)

	assert.NotNil(t, cfg)
	assert.NoError(t, err)
}

func TestUnrecognizedLogLevel(t *testing.T) {
	cliTest := testSetup()

	cfg, err := cliTest.Parse(
		"-server-url", "hello",
		"-log-level=potato",
	)

	assert.Nil(t, cfg)
	assert.Error(t, err)
	assert.EqualError(t, err, "unrecognized log level: potato")
}
