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
func NewRouter(renderer *htmltemplate.Renderer, app *app.App) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	// Setup routes
	router.Path("/").Handler(indexHandler(renderer, app)).Methods(http.MethodGet)

	router.PathPrefix("/upload/").Handler(
		http.StripPrefix("/upload/", uploadHandler(renderer, app)),
	).Methods(http.MethodPut)

	return router
}
