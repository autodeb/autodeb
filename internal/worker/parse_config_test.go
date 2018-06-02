package worker_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/log"
	"salsa.debian.org/autodeb-team/autodeb/internal/worker"
)

func TestNotFound(t *testing.T) {
	fs := filesystem.NewMemMapFS()

	cfg, err := worker.ParseConfig("test.cfg", fs)
	assert.Nil(t, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file does not exist")
}

func TestParseConfigUnknownKey(t *testing.T) {
	fs := filesystem.NewMemMapFS()

	f, err := fs.Create("server.cfg")
	assert.NoError(t, err)

	fmt.Fprintln(f, "unknownkey=11")

	f.Close()

	cfg, err := worker.ParseConfig("server.cfg", fs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unrecognized configuration key: unknownkey")
	assert.Nil(t, cfg)
}

func TestDefaultConfig(t *testing.T) {
	fs := filesystem.NewMemMapFS()

	f, err := fs.Create("empty.cfg")
	assert.NoError(t, err)
	f.Close()

	cfg, err := worker.ParseConfig("empty.cfg", fs)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "jobs", cfg.WorkingDirectory, "the config should contain defaults")
}

func TestParseConfig(t *testing.T) {
	fs := filesystem.NewMemMapFS()

	f, err := fs.Create("worker.cfg")
	assert.NoError(t, err)

	fmt.Fprintln(f, `server_url="test-server-url"`)
	fmt.Fprintln(f, `access_token="test-access-token"`)
	fmt.Fprintln(f, `working_directory="test-working-directory"`)
	fmt.Fprintln(f, `log_level="error"`)
	fmt.Fprintln(f, `runner_count=42`)

	f.Close()

	cfg, err := worker.ParseConfig("worker.cfg", fs)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "test-server-url", cfg.ServerURL)
	assert.Equal(t, "test-access-token", cfg.AccessToken)
	assert.Equal(t, "test-working-directory", cfg.WorkingDirectory)
	assert.Equal(t, log.ErrorLevel, cfg.LogLevel)
	assert.Equal(t, 42, cfg.RunnerCount)
}
