// Package server contains the core of the autodeb server. It creates all
// dependencies injects them at the right place.
package server

import (
	"context"
	"crypto/rand"
	"io"
	"net/url"

	"github.com/gorilla/sessions"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/htmltemplate"
	"salsa.debian.org/autodeb-team/autodeb/internal/http"
	"salsa.debian.org/autodeb-team/autodeb/internal/log"
	"salsa.debian.org/autodeb-team/autodeb/internal/oauth"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router"
)

// Server is the main server. It has dput-compatible interface
// and provides REST API.
type Server struct {
	httpServer *http.Server
}

// New creates a Server
func New(cfg *Config, loggingOutput io.Writer) (*Server, error) {
	db, err := database.NewDatabase(cfg.DB.Driver, cfg.DB.ConnectionString)
	if err != nil {
		return nil, err
	}

	dataFS, err := filesystem.NewFS(cfg.DataDirectory)
	if err != nil {
		return nil, err
	}

	oauthProvider, err := getOAuthProvider(cfg)
	if err != nil {
		return nil, err
	}

	templatesFS, err := filesystem.NewFS(cfg.TemplatesDirectory)
	if err != nil {
		return nil, err
	}

	renderer := htmltemplate.NewRenderer(templatesFS, cfg.TemplatesCacheEnabled)

	staticFilesFS, err := filesystem.NewFS(cfg.StaticFilesDirectory)
	if err != nil {
		return nil, err
	}

	sessionsStore, err := getSessionStore()
	if err != nil {
		return nil, err
	}

	logger := log.New(loggingOutput)
	logger.SetLevel(cfg.LogLevel)

	app, err := app.NewApp(
		cfg.AppConfig,
		db,
		dataFS,
		oauthProvider,
		renderer,
		filesystem.NewHTTPFS(staticFilesFS),
		sessionsStore,
		logger,
	)
	if err != nil {
		return nil, err
	}

	router := router.NewRouter(app)

	httpServer, err := http.NewHTTPServer(cfg.HTTP.Address, router, logger)
	if err != nil {
		return nil, err
	}

	server := Server{
		httpServer: httpServer,
	}

	return &server, nil
}

func getSessionStore() (sessions.Store, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	// TODO: Ask for the session secret in the CLI
	// instead of generating a random one
	store := sessions.NewCookieStore(b)

	return store, nil
}

func getOAuthProvider(cfg *Config) (oauth.Provider, error) {
	baseURL, err := url.Parse(cfg.OAuth.BaseURL)
	if err != nil {
		return nil, err
	}

	provider, err := oauth.NewProvider(
		cfg.OAuth.Provider,
		baseURL,
		cfg.OAuth.ClientID,
		cfg.OAuth.ClientSecret,
	)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// Shutdown will gracefully stop the server
func (srv *Server) Shutdown(ctx context.Context) error {
	return srv.httpServer.Shutdown(ctx)
}
