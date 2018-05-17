// Package router provides the main router. It translates http requests to App
// actions and creates http responses.
package router

import (
	"net/http"

	"github.com/gorilla/mux"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/internal/endpoints/api"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/internal/endpoints/auth"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/internal/endpoints/uploadqueue"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/internal/endpoints/webpages"
)

// NewRouter creates the main router for the application.
func NewRouter(app *app.App) http.Handler {

	router := mux.NewRouter().StrictSlash(true)

	// Upload Queue
	router.PathPrefix("/upload/").Handler(
		http.StripPrefix("/upload/", uploadqueue.UploadHandler(app)),
	).Methods(http.MethodPut, http.MethodPost)

	// Static files (for the web)
	router.PathPrefix("/static/").Handler(
		http.StripPrefix(
			"/static/",
			http.FileServer(app.StaticFS()),
		),
	).Methods(http.MethodGet)

	// Authentification
	router.Path("/auth/login").Handler(auth.LoginGetHandler(app))
	router.Path("/auth/logout").Handler(auth.LogoutGetHandler(app))
	router.Path("/auth/callback").Handler(auth.CallbackGetHandler(app))

	// Web pages
	router.Path("/").Handler(webpages.IndexGetHandler(app)).Methods(http.MethodGet)
	router.Path("/uploads").Handler(webpages.UploadsGetHandler(app)).Methods(http.MethodGet)
	router.Path("/jobs").Handler(webpages.JobsGetHandler(app)).Methods(http.MethodGet)
	router.Path("/profile").Handler(webpages.ProfileGetHandler(app)).Methods(http.MethodGet)

	// REST API Router
	restAPIRouter := router.PathPrefix("/api/").Subrouter()

	// ==== Jobs API ====
	restAPIRouter.Path("/jobs/next").Handler(api.JobsNextPostHandler(app)).Methods(http.MethodPost)
	restAPIRouter.Path("/jobs/{jobID:[0-9]+}").Handler(api.JobGetHandler(app)).Methods(http.MethodGet)
	restAPIRouter.Path("/jobs/{jobID:[0-9]+}/status/{jobStatus:[0-9]+}").Handler(api.JobStatusPostHandler(app)).Methods(http.MethodPost)
	restAPIRouter.Path("/jobs/{jobID:[0-9]+}/log").Handler(api.JobLogPostHandler(app)).Methods(http.MethodPost)
	restAPIRouter.Path("/jobs/{jobID:[0-9]+}/log.txt").Handler(api.JobLogTxtGetHandler(app)).Methods(http.MethodGet)

	// ==== Uploads API ====
	restAPIRouter.Path("/uploads/{uploadID:[0-9]+}/dsc").Handler(api.UploadDSCGetHandler(app)).Methods(http.MethodGet)
	// This route is kept for dget compatibility. Dget requires that the URL ends with /<something>.dsc (TODO: open a bug)
	restAPIRouter.Path("/uploads/{uploadID:[0-9]+}/source.dsc").Handler(api.UploadDSCGetHandler(app)).Methods(http.MethodGet)
	restAPIRouter.Path("/uploads/{uploadID:[0-9]+}/files").Handler(api.UploadFilesGetHandler(app)).Methods(http.MethodGet)
	restAPIRouter.Path("/uploads/{uploadID:[0-9]+}/{filename}").Handler(api.UploadFileGetHandler(app)).Methods(http.MethodGet)

	return router
}
