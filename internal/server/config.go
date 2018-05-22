package server

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/log"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
)

// Config holds configuration that is needed by Server to initialize all of its
// dependencies.
type Config struct {
	DB                    *DBConfig
	HTTP                  *HTTPServerConfig
	Auth                  *AuthConfig
	ContextConfig         *appctx.Config
	DataDirectory         string
	TemplatesDirectory    string
	StaticFilesDirectory  string
	TemplatesCacheEnabled bool
	LogLevel              log.Level
}

// AuthConfig holds configuration related to the authentification backend
type AuthConfig struct {
	AuthentificationBackend string
	OAuth                   *OAuthConfig
}

// OAuthConfig holds configuration related to the OAuth provider
type OAuthConfig struct {
	Provider     string
	BaseURL      string
	ClientID     string
	ClientSecret string
}

// HTTPServerConfig holds configuration related to the HTTP server
type HTTPServerConfig struct {
	Address string
}

// DBConfig holds configuration related to the database
type DBConfig struct {
	Driver           string
	ConnectionString string
}
