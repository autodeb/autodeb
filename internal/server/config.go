package server

import (
	"salsa.debian.org/aviau/autodeb/internal/http"
	"salsa.debian.org/aviau/autodeb/internal/server/app"
	"salsa.debian.org/aviau/autodeb/internal/server/database"
)

// Config contains configuration for Server
type Config struct {
	App      *app.Config
	HTTP     *http.ServerConfig
	Database *database.Config
}
