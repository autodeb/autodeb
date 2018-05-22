// Package server contains the core of the autodeb server. It creates all
// dependencies injects them at the right place.
package server

import (
	"context"
	"crypto/rand"
	"io"
	"net/url"

	"github.com/gorilla/sessions"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/htmltemplate"
	"salsa.debian.org/autodeb-team/autodeb/internal/http"
	"salsa.debian.org/autodeb-team/autodeb/internal/log"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/auth"
	authDisabled "salsa.debian.org/autodeb-team/autodeb/internal/server/auth/disabled"
	authOAuth "salsa.debian.org/autodeb-team/autodeb/internal/server/auth/oauth"
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

	dataFS := filesystem.NewBasePathFS(filesystem.NewOsFS(), cfg.DataDirectory)
	staticFilesFS := filesystem.NewBasePathFS(filesystem.NewOsFS(), cfg.StaticFilesDirectory)
	templatesFS := filesystem.NewBasePathFS(filesystem.NewOsFS(), cfg.TemplatesDirectory)

	renderer := htmltemplate.NewRenderer(templatesFS, cfg.TemplatesCacheEnabled)

	authBackend, err := getAuthBackend(cfg, db)
	if err != nil {
		return nil, err
	}

	logger := log.New(loggingOutput)
	logger.SetLevel(cfg.LogLevel)

	app, err := app.NewApp(
		cfg.AppConfig,
		db,
		dataFS,
		renderer,
		filesystem.NewHTTPFS(staticFilesFS),
		authBackend,
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

func getAuthBackend(cfg *Config, db *database.Database) (auth.Backend, error) {
	switch cfg.Auth.AuthentificationBackend {
	case "oauth":
		return getOAuthBackend(cfg, db)
	case "disabled":
		return authDisabled.NewBackend(), nil
	default:
		return nil, errors.Errorf("unrecognized authentification backend: %s (use oauth or disabled)", cfg.Auth.AuthentificationBackend)
	}
}

func getOAuthBackend(cfg *Config, db *database.Database) (auth.Backend, error) {
	sessionStore, err := getSessionStore()
	if err != nil {
		return nil, err
	}

	baseURL, err := url.Parse(cfg.Auth.OAuth.BaseURL)
	if err != nil {
		return nil, err
	}

	oauthProvider, err := authOAuth.NewProvider(
		cfg.Auth.OAuth.Provider,
		baseURL,
		cfg.Auth.OAuth.ClientID,
		cfg.Auth.OAuth.ClientSecret,
	)
	if err != nil {
		return nil, err
	}

	authBackend := authOAuth.NewBackend(
		db,
		sessionStore,
		oauthProvider,
		cfg.AppConfig.ServerURL,
	)

	return authBackend, nil
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

// Shutdown will gracefully stop the server
func (srv *Server) Shutdown(ctx context.Context) error {
	return srv.httpServer.Shutdown(ctx)
}
