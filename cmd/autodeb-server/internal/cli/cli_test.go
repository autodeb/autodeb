package cli_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"salsa.debian.org/autodeb-team/autodeb/cmd/autodeb-server/internal/cli"
	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/config"
)

type cliTest struct {
	outputWriter bytes.Buffer
	fs           filesystem.FS
}

func (cliTest *cliTest) Parse(args ...string) (*config.Config, error) {
	return cli.Parse(args, cliTest.fs, &cliTest.outputWriter)
}

func testSetup() *cliTest {
	fs := filesystem.NewMemMapFS()

	cliTest := &cliTest{
		fs: fs,
	}

	return cliTest
}

func TestNoConfigFound(t *testing.T) {
	cliTest := testSetup()

	cfg, err := cliTest.Parse("")

	assert.Nil(t, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not open configuration file")
}

func TestEmptyConfig(t *testing.T) {
	cliTest := testSetup()

	f, err := cliTest.fs.Create("test.cfg")
	assert.NoError(t, err)
	f.Close()

	cfg, err := cliTest.Parse(
		"-config=test.cfg",
	)

	assert.NotNil(t, cfg)
	assert.NoError(t, err)
	assert.Equal(t, ":8071", cfg.HTTP.Address, "the config should use the defaults")
}

func TestUnrecognizedArgument(t *testing.T) {
	cliTest := testSetup()

	f, err := cliTest.fs.Create("test.cfg")
	assert.NoError(t, err)
	f.Close()

	cfg, err := cliTest.Parse(
		"-config=test.cfg",
		"test",
	)

	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "unrecognized argument")
}
