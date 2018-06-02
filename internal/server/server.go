// Package server contains the core of the autodeb server. It creates all
// dependencies injects them at the right place.
package server

import (
	"context"
	"crypto/rand"
	"io"
	"net/url"

	gorillaSessions "github.com/gorilla/sessions"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/htmltemplate"
	"salsa.debian.org/autodeb-team/autodeb/internal/http"
	"salsa.debian.org/autodeb-team/autodeb/internal/http/sessions"
	"salsa.debian.org/autodeb-team/autodeb/internal/log"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/auth"
	authDisabled "salsa.debian.org/autodeb-team/autodeb/internal/server/auth/disabled"
	authOAuth "salsa.debian.org/autodeb-team/autodeb/internal/server/auth/oauth"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/config"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/database"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/services"
)

// Server is the main server. It has dput-compatible interface
// and provides REST API.
type Server struct {
	httpServer *http.Server
}

// New creates a Server
func New(cfg *config.Config, loggingOutput io.Writer) (*Server, error) {
	db, err := database.NewDatabase(cfg.DB.Driver, cfg.DB.ConnectionString)
	if err != nil {
		return nil, err
	}

	dataFS := filesystem.NewBasePathFS(filesystem.NewOsFS(), cfg.DataDirectory)
	staticFilesFS := filesystem.NewBasePathFS(filesystem.NewOsFS(), cfg.StaticFilesDirectory)
	templatesFS := filesystem.NewBasePathFS(filesystem.NewOsFS(), cfg.TemplatesDirectory)

	renderer := htmltemplate.NewRenderer(templatesFS, cfg.TemplatesCacheEnabled)

	sessionsManager, err := getSessionManager()
	if err != nil {
		return nil, err
	}

	authBackend, err := getAuthBackend(cfg, db, sessionsManager)
	if err != nil {
		return nil, err
	}

	logger := log.New(loggingOutput)
	logger.SetLevel(cfg.LogLevel)

	services, err := services.New(db, dataFS, &cfg.ServerURL.URL)
	if err != nil {
		return nil, err
	}

	appCtx := appctx.New(
		cfg,
		renderer,
		filesystem.NewHTTPFS(staticFilesFS),
		authBackend,
		sessionsManager,
		services,
		logger,
	)

	router := router.NewRouter(appCtx)

	httpServer, err := http.NewHTTPServer(cfg.HTTP.Address, router, logger)
	if err != nil {
		return nil, err
	}

	server := Server{
		httpServer: httpServer,
	}

	return &server, nil
}

func getAuthBackend(cfg *config.Config, db *database.Database, sessionsManager *sessions.Manager) (auth.Backend, error) {
	switch cfg.Auth.AuthentificationBackend {
	case "oauth":
		return getOAuthBackend(cfg, db, sessionsManager)
	case "disabled":
		return authDisabled.NewBackend(), nil
	default:
		return nil, errors.Errorf("unrecognized authentification backend: %s (use oauth or disabled)", cfg.Auth.AuthentificationBackend)
	}
}

func getOAuthBackend(cfg *config.Config, db *database.Database, sessionsManager *sessions.Manager) (auth.Backend, error) {
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
		sessionsManager,
		oauthProvider,
		&cfg.ServerURL.URL,
	)

	return authBackend, nil
}

func getSessionManager() (*sessions.Manager, error) {
	// TODO: Ask for the session secret in the CLI
	// instead of generating a random one
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	sessionStore := gorillaSessions.NewCookieStore(b)

	sessionsManager := sessions.NewManager(sessionStore, "autodeb")

	return sessionsManager, nil
}

// Shutdown will gracefully stop the server
func (srv *Server) Shutdown(ctx context.Context) error {
	return srv.httpServer.Shutdown(ctx)
}
