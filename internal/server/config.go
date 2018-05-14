package server

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/log"
)

// Config contains configuration for Server
type Config struct {
	DB                    DBConfig
	HTTP                  HTTPServerConfig
	DataDirectory         string
	TemplatesDirectory    string
	StaticFilesDirectory  string
	TemplatesCacheEnabled bool
	LogLevel              log.Level
}

// HTTPServerConfig holds configuration related to the HTTP server
type HTTPServerConfig struct {
	Address string
	Port    int
}

// DBConfig holds configuration related to the database
type DBConfig struct {
	Driver           string
	ConnectionString string
}
