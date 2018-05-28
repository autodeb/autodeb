// Package appctx implements the application's context.
package appctx

import (
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/htmltemplate"
	"salsa.debian.org/autodeb-team/autodeb/internal/http/sessions"
	"salsa.debian.org/autodeb-team/autodeb/internal/log"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/auth"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/config"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/services"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/services/jobs"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/services/pgp"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/services/tokens"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/services/uploads"
)

// Context is the application's context. It holds everityhing that is needed to
// serve a request.
type Context struct {
	config          *config.Config
	renderer        *htmltemplate.Renderer
	staticFS        http.FileSystem
	authBackend     auth.Backend
	sessionsManager *sessions.Manager
	services        *services.Services
	logger          log.Logger
}

// New create an application context
func New(
	config *config.Config,
	renderer *htmltemplate.Renderer,
	staticFS http.FileSystem,
	authBackend auth.Backend,
	sessionsManager *sessions.Manager,
	services *services.Services,
	logger log.Logger) *Context {

	context := &Context{
		config:          config,
		renderer:        renderer,
		staticFS:        staticFS,
		authBackend:     authBackend,
		sessionsManager: sessionsManager,
		services:        services,
		logger:          logger,
	}

	return context
}

// Sessions returns the sessions manager
func (ctx *Context) Sessions() *sessions.Manager {
	return ctx.sessionsManager
}

// Logger returns the logger
func (ctx *Context) Logger() log.Logger {
	return ctx.logger
}

// AuthBackend returns the authentification service
func (ctx *Context) AuthBackend() auth.Backend {
	return ctx.authBackend
}

// Config returns the context's config
func (ctx *Context) Config() *config.Config {
	return ctx.config
}

// StaticFS contains static files to be served over http
func (ctx *Context) StaticFS() http.FileSystem {
	return ctx.staticFS
}

// UploadsService returns the uploads service
func (ctx *Context) UploadsService() *uploads.Service {
	return ctx.services.Uploads()
}

// PGPService returns the PGP service
func (ctx *Context) PGPService() *pgp.Service {
	return ctx.services.PGP()
}

// JobsService returns the jobs service
func (ctx *Context) JobsService() *jobs.Service {
	return ctx.services.Jobs()
}

// TokensService returns the Tokens service
func (ctx *Context) TokensService() *tokens.Service {
	return ctx.services.Tokens()
}

// TemplatesRenderer returns the template renderer
func (ctx *Context) TemplatesRenderer() *htmltemplate.Renderer {
	return ctx.renderer
}
