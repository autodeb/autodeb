// Package router provides the main router. It translates http requests to App
// actions and creates http responses.
package router

import (
	"net/http"

	"github.com/gorilla/mux"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/internal/endpoints/api"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/internal/endpoints/uploadqueue"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/internal/endpoints/webpages"
)

// NewRouter creates the main router.
func NewRouter(appCtx *appctx.Context) http.Handler {

	router := mux.NewRouter().StrictSlash(true)

	// Upload Queue
	router.PathPrefix("/upload/").Handler(
		http.StripPrefix("/upload/", uploadqueue.UploadHandler(appCtx)),
	).Methods(http.MethodPut, http.MethodPost)

	// Static files (for the web)
	router.PathPrefix("/static/").Handler(
		http.StripPrefix(
			"/static/",
			http.FileServer(appCtx.StaticFS()),
		),
	).Methods(http.MethodGet)

	// Authentification
	router.Path("/login").Handler(appCtx.AuthBackend().LoginHandler()).Methods(http.MethodGet)
	router.Path("/logout").Handler(appCtx.AuthBackend().LogoutHandler()).Methods(http.MethodGet)
	router.PathPrefix("/auth/").Handler(http.StripPrefix("/auth/", appCtx.AuthBackend().AuthHandler()))

	// Web pages
	router.Path("/").Handler(webpages.IndexGetHandler(appCtx)).Methods(http.MethodGet)
	router.Path("/uploads").Handler(webpages.UploadsGetHandler(appCtx)).Methods(http.MethodGet)
	router.Path("/jobs").Handler(webpages.JobsGetHandler(appCtx)).Methods(http.MethodGet)
	router.Path("/profile").Handler(webpages.ProfileGetHandler(appCtx)).Methods(http.MethodGet)
	router.Path("/profile/add-pgp-key").Handler(webpages.AddPGPKeyPostHandler(appCtx)).Methods(http.MethodPost)

	// REST API Router
	restAPIRouter := router.PathPrefix("/api/").Subrouter()

	// ==== Jobs API ====
	restAPIRouter.Path("/jobs/next").Handler(api.JobsNextPostHandler(appCtx)).Methods(http.MethodPost)
	restAPIRouter.Path("/jobs/{jobID:[0-9]+}").Handler(api.JobGetHandler(appCtx)).Methods(http.MethodGet)
	restAPIRouter.Path("/jobs/{jobID:[0-9]+}/status/{jobStatus:[0-9]+}").Handler(api.JobStatusPostHandler(appCtx)).Methods(http.MethodPost)
	restAPIRouter.Path("/jobs/{jobID:[0-9]+}/log").Handler(api.JobLogPostHandler(appCtx)).Methods(http.MethodPost)
	restAPIRouter.Path("/jobs/{jobID:[0-9]+}/log.txt").Handler(api.JobLogTxtGetHandler(appCtx)).Methods(http.MethodGet)

	// ==== Uploads API ====
	restAPIRouter.Path("/uploads/{uploadID:[0-9]+}/dsc").Handler(api.UploadDSCGetHandler(appCtx)).Methods(http.MethodGet)
	// This route is kept for dget compatibility. Dget requires that the URL ends with /<something>.dsc (TODO: open a bug)
	restAPIRouter.Path("/uploads/{uploadID:[0-9]+}/source.dsc").Handler(api.UploadDSCGetHandler(appCtx)).Methods(http.MethodGet)
	restAPIRouter.Path("/uploads/{uploadID:[0-9]+}/files").Handler(api.UploadFilesGetHandler(appCtx)).Methods(http.MethodGet)
	restAPIRouter.Path("/uploads/{uploadID:[0-9]+}/{filename}").Handler(api.UploadFileGetHandler(appCtx)).Methods(http.MethodGet)

	return router
}
