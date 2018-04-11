package server

import (
	"salsa.debian.org/aviau/autopkgupdate/internal/http"
	"salsa.debian.org/aviau/autopkgupdate/internal/server/app"
	"salsa.debian.org/aviau/autopkgupdate/internal/server/database"
)

// Config contains configuration for Server
type Config struct {
	App      *app.Config
	HTTP     *http.ServerConfig
	Database *database.Config
}
