// Package api provides the main router. It translates http requests to App
// actions and creates http responses.
package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
)

// NewRouter creates the main router for the application.
func NewRouter(app *app.App) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	// Setup routes
	router.Path("/").Handler(indexHandler(app)).Methods(http.MethodGet)

	router.PathPrefix("/upload/").Handler(
		http.StripPrefix("/upload/", uploadHandler(app)),
	)

	return router
}
