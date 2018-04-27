package server

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/http"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
)

// Config contains configuration for Server
type Config struct {
	HTTP               *http.ServerConfig
	Database           *database.Config
	DataDirectory      string
	TemplatesDirectory string
}
