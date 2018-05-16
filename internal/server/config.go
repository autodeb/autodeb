package server

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/log"
)

// Config contains configuration for Server
type Config struct {
	DB                    DBConfig
	HTTP                  HTTPServerConfig
	OAuth                 OAuthConfig
	DataDirectory         string
	TemplatesDirectory    string
	StaticFilesDirectory  string
	TemplatesCacheEnabled bool
	LogLevel              log.Level
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
