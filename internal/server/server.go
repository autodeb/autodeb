// Package server contains the core of the autodeb server. It creates all
// dependencies injects them at the right place.
package server

import (
	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/htmltemplate"
	"salsa.debian.org/autodeb-team/autodeb/internal/http"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/api"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
)

// Server is the main server. It has dput-compatible interface
// and provides REST API.
type Server struct {
	httpServer *http.Server
}

// NewServer creates a Server
func NewServer(cfg *Config) (*Server, error) {
	db, err := database.NewDatabase(cfg.Database)
	if err != nil {
		return nil, err
	}

	dataFS, err := filesystem.NewFS(cfg.DataDirectory)
	if err != nil {
		return nil, err
	}

	app, err := app.NewApp(db, dataFS)
	if err != nil {
		return nil, err
	}

	templatesFS, err := filesystem.NewFS(cfg.TemplatesDirectory)
	if err != nil {
		return nil, err
	}

	renderer := htmltemplate.NewRenderer(templatesFS)

	staticFilesFS, err := filesystem.NewFS(cfg.StaticFilesDirectory)
	if err != nil {
		return nil, err
	}

	router := api.NewRouter(
		renderer,
		filesystem.NewHTTPFS(staticFilesFS),
		app,
	)

	httpServer, err := http.NewHTTPServer(router, cfg.HTTP)
	if err != nil {
		return nil, err
	}

	server := Server{
		httpServer: httpServer,
	}

	return &server, nil
}

// Close will shutdown the server
func (srv *Server) Close() error {
	return srv.httpServer.Close()
}
