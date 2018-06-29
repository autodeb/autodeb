package config

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/log"
	"salsa.debian.org/autodeb-team/autodeb/internal/net/url"
)

// Config holds configuration that is needed by Server to initialize all of its
// dependencies.
type Config struct {
	DB                    *DBConfig         `toml:"database"`
	HTTP                  *HTTPServerConfig `toml:"http"`
	Auth                  *AuthConfig       `toml:"authentication"`
	ServerURL             *url.URL          `toml:"server_url"`
	Aptly                 *Aptly            `toml:"aptly"`
	DataDirectory         string            `toml:"data_directory"`
	TemplatesDirectory    string            `toml:"templates_directory"`
	StaticFilesDirectory  string            `toml:"static_files_directory"`
	TemplatesCacheEnabled bool              `toml:"templates_cache_enabled"`
	LogLevel              log.Level         `toml:"log_level"`
}

// Aptly holds configuration related to aptly
type Aptly struct {
	APIURL            *url.URL `toml:"api_url"`
	RepositoryBaseURL *url.URL `toml:"repository_base_url"`
}

// AuthConfig holds configuration related to the authentification backend
type AuthConfig struct {
	AuthentificationBackend string       `toml:"backend"`
	OAuth                   *OAuthConfig `toml:"oauth"`
}

// OAuthConfig holds configuration related to the OAuth provider
type OAuthConfig struct {
	Provider     string `toml:"provider"`
	BaseURL      string `toml:"base_url"`
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
}

// HTTPServerConfig holds configuration related to the HTTP server
type HTTPServerConfig struct {
	Address string `toml:"address"`
}

// DBConfig holds configuration related to the database
type DBConfig struct {
	Driver           string `toml:"driver"`
	ConnectionString string `toml:"connection_string"`
}
