package config_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/log"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/config"
)

func TestNotFound(t *testing.T) {
	fs := filesystem.NewMemMapFS()

	cfg, err := config.ParseConfig("test.cfg", fs)
	assert.Nil(t, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file does not exist")
}

func TestDefaultConfig(t *testing.T) {
	fs := filesystem.NewMemMapFS()

	f, err := fs.Create("empty.cfg")
	assert.NoError(t, err)
	f.Close()

	cfg, err := config.ParseConfig("empty.cfg", fs)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, ":8071", cfg.HTTP.Address, "the config should contain defaults")
}

func TestParseConfigUnknownKey(t *testing.T) {
	fs := filesystem.NewMemMapFS()

	f, err := fs.Create("server.cfg")
	assert.NoError(t, err)

	fmt.Fprintln(f, "unknownkey=11")

	f.Close()

	cfg, err := config.ParseConfig("server.cfg", fs)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unrecognized configuration key: unknownkey")
	assert.Nil(t, cfg)
}

func TestParseConfig(t *testing.T) {
	fs := filesystem.NewMemMapFS()

	f, err := fs.Create("server.cfg")
	assert.NoError(t, err)

	var configText = `
server_url = "https://test-server-url:1234"
data_directory = "test-data-directory"
templates_directory = "test-templates-directory"
static_files_directory = "test-static-files-directory"
templates_cache_enabled = false
log_level = "warning"

[database]
driver = "test-driver"
connection_string = "test-connection-string"

[http]
address = ":8071"

[authentication]
backend = "test-backend"

    [authentication.oauth]
	provider = "test-provider"
	base_url = "test-base-url"
	client_id = "test-client-id"
	client_secret = "test-client-secret"

`

	fmt.Fprintln(f, configText)

	f.Close()

	cfg, err := config.ParseConfig("server.cfg", fs)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Base config
	assert.Equal(t, "https://test-server-url:1234", cfg.ServerURL.String())
	assert.Equal(t, "test-data-directory", cfg.DataDirectory)
	assert.Equal(t, "test-static-files-directory", cfg.StaticFilesDirectory)
	assert.Equal(t, false, cfg.TemplatesCacheEnabled)
	assert.Equal(t, log.WarningLevel, cfg.LogLevel)

	// Database
	assert.Equal(t, "test-driver", cfg.DB.Driver)
	assert.Equal(t, "test-connection-string", cfg.DB.ConnectionString)

	// HTTP
	assert.Equal(t, ":8071", cfg.HTTP.Address)

	// Auth
	assert.Equal(t, "test-backend", cfg.Auth.AuthentificationBackend)
	assert.Equal(t, "test-provider", cfg.Auth.OAuth.Provider)
	assert.Equal(t, "test-base-url", cfg.Auth.OAuth.BaseURL)
	assert.Equal(t, "test-client-id", cfg.Auth.OAuth.ClientID)
	assert.Equal(t, "test-client-secret", cfg.Auth.OAuth.ClientSecret)
}
