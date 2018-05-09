// Package router provides the main router. It translates http requests to App
// actions and creates http responses.
package router

import (
	"net/http"

	"github.com/gorilla/mux"

	"salsa.debian.org/autodeb-team/autodeb/internal/htmltemplate"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/internal/api"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/internal/uploads"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/internal/webpages"
)

// NewRouter creates the main router for the application.
func NewRouter(renderer *htmltemplate.Renderer, staticFS http.FileSystem, app *app.App) http.Handler {

	router := mux.NewRouter().StrictSlash(true)

	// Upload API
	router.PathPrefix("/upload/").Handler(
		http.StripPrefix("/upload/", uploads.UploadHandler(app)),
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
	//    Jobs
	restAPIRouter.Path("/jobs/next").Handler(api.JobsNextPostHandler(app)).Methods(http.MethodPost)
	//    Upload
	//restAPIRouter.Path("/uploads/{id:[0-9]+}")).Handler()

	return router
}
