// Package api provides the main router. It translates http requests to App
// actions and creates http responses.
package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"salsa.debian.org/autodeb-team/autodeb/internal/htmltemplate"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
)

// NewRouter creates the main router for the application.
func NewRouter(renderer *htmltemplate.Renderer, staticFS http.FileSystem, app *app.App) http.Handler {

	router := mux.NewRouter().StrictSlash(true)

	// Upload API
	router.PathPrefix("/upload/").Handler(
		http.StripPrefix("/upload/", uploadHandler(renderer, app)),
	).Methods(http.MethodPut)

	// Static files (for the web)
	router.PathPrefix("/static/").Handler(
		http.StripPrefix(
			"/static/",
			http.FileServer(staticFS),
		),
	).Methods(http.MethodGet)

	// Web pages
	router.Path("/").Handler(indexGetHandler(renderer, app)).Methods(http.MethodGet)
	router.Path("/uploads").Handler(uploadsGetHandler(renderer, app)).Methods(http.MethodGet)
	router.Path("/jobs").Handler(jobsGetHandler(renderer, app)).Methods(http.MethodGet)

	return router
}
