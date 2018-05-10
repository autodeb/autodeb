package cli_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"salsa.debian.org/autodeb-team/autodeb/cmd/autodeb-server/internal/cli"
	"salsa.debian.org/autodeb-team/autodeb/internal/server"
)

type cliTest struct {
	outputWriter bytes.Buffer
	errorWriter  bytes.Buffer
}

func (cliTest *cliTest) Parse(args ...string) (*server.Config, error) {
	return cli.Parse(args, &cliTest.outputWriter, &cliTest.errorWriter)
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
