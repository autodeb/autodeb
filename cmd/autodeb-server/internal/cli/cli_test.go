package cli_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"salsa.debian.org/autodeb-team/autodeb/cmd/autodeb-server/internal/cli"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/config"
)

type cliTest struct {
	outputWriter bytes.Buffer
}

func (cliTest *cliTest) Parse(args ...string) (*config.Config, error) {
	return cli.Parse(args, &cliTest.outputWriter)
}

func testSetup() *cliTest {
	cliTest := cliTest{}
	return &cliTest
}

func TestEmptyArgsNoError(t *testing.T) {
	cliTest := testSetup()

	cfg, err := cliTest.Parse("")

	assert.NotNil(t, cfg)
	assert.NoError(t, err)
}

func TestUnrecognizedLogLevel(t *testing.T) {
	cliTest := testSetup()

	cfg, err := cliTest.Parse(
		"-log-level=potato",
	)

	assert.Nil(t, cfg)
	assert.Error(t, err)
	assert.EqualError(t, err, "unrecognized log level: potato")
}
