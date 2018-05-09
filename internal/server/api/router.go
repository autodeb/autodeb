// Package api provides the main router. It translates http requests to App
// actions and creates http responses.
package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"salsa.debian.org/autodeb-team/autodeb/internal/htmltemplate"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/api/webpages"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
)

// NewRouter creates the main router for the application.
func NewRouter(renderer *htmltemplate.Renderer, staticFS http.FileSystem, app *app.App) http.Handler {

	router := mux.NewRouter().StrictSlash(true)

	// Upload API
	router.PathPrefix("/upload/").Handler(
		http.StripPrefix("/upload/", uploadHandler(app)),
	).Methods(http.MethodPut)

	// Static files (for the web)
	router.PathPrefix("/static/").Handler(
		http.StripPrefix(
			"/static/",
			http.FileServer(staticFS),
		),
	).Methods(http.MethodGet)

	// Web pages
	router.Path("/").Handler(webpages.IndexGetHandler(renderer, app)).Methods(http.MethodGet)
	router.Path("/uploads").Handler(webpages.UploadsGetHandler(renderer, app)).Methods(http.MethodGet)
	router.Path("/jobs").Handler(webpages.JobsGetHandler(renderer, app)).Methods(http.MethodGet)

	// REST API Router
	restAPIRouter := router.PathPrefix("/api/").Subrouter()
	restAPIRouter.Path("/jobs/next").Handler(jobsNextPostHandler(app)).Methods(http.MethodPost)

	return router
}
