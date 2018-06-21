// Package router provides the main router. It translates http requests to App
// actions and creates http responses.
package router

import (
	"net/http"

	"github.com/gorilla/mux"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/api"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/aptly"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/uploadqueue"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/webpages"
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

	// ==== Web pages: General ====
	router.Path("/").Handler(webpages.IndexGetHandler(appCtx)).Methods(http.MethodGet)
	router.Path("/instructions").Handler(webpages.InstructionsGetHandler(appCtx)).Methods(http.MethodGet)
	router.Path("/uploads").Handler(webpages.UploadsGetHandler(appCtx)).Methods(http.MethodGet)
	router.Path("/uploads/{uploadID:[0-9]+}").Handler(webpages.UploadGetHandler(appCtx)).Methods(http.MethodGet)
	router.Path("/jobs").Handler(webpages.JobsGetHandler(appCtx)).Methods(http.MethodGet)
	router.Path("/jobs/{jobID:[0-9]+}").Handler(webpages.JobGetHandler(appCtx)).Methods(http.MethodGet)
	router.Path("/archive-upgrades").Handler(webpages.ArchiveUpgradesGetHandler(appCtx)).Methods(http.MethodGet)
	router.Path("/archive-upgrades/{archiveUpgradeID:[0-9]+}").Handler(webpages.ArchiveUpgradeGetHandler(appCtx)).Methods(http.MethodGet)
	router.Path("/new-archive-upgrade").Handler(webpages.NewArchiveUpgradeGetHandler(appCtx)).Methods(http.MethodGet)
	router.Path("/new-archive-upgrade").Handler(webpages.NewArchiveUpgradePostHandler(appCtx)).Methods(http.MethodPost)

	// ==== Web pages: Profile ====
	router.Path("/profile").Handler(webpages.ProfileGetHandler(appCtx)).Methods(http.MethodGet)

	router.Path("/profile/pgp-keys").Handler(webpages.ProfilePGPKeysGetHandler(appCtx)).Methods(http.MethodGet)
	router.Path("/profile/add-pgp-key").Handler(webpages.AddPGPKeyPostHandler(appCtx)).Methods(http.MethodPost)
	router.Path("/profile/remove-pgp-key").Handler(webpages.RemovePGPKeyPostHandler(appCtx)).Methods(http.MethodPost)

	router.Path("/profile/access-tokens").Handler(webpages.ProfileAccessTokensGetHandler(appCtx)).Methods(http.MethodGet)
	router.Path("/profile/create-access-token").Handler(webpages.CreateAccessTokenPostHandler(appCtx)).Methods(http.MethodPost)
	router.Path("/profile/remove-access-token").Handler(webpages.RemoveAccessTokenPostHandler(appCtx)).Methods(http.MethodPost)

	// APTLY
	router.PathPrefix("/aptly/").Handler(aptly.Handler(appCtx, "/aptly/"))

	// REST API Router
	restAPIRouter := router.PathPrefix("/api/").Subrouter()

	// ==== User ====
	restAPIRouter.Path("/user").Handler(api.UserGetHandler(appCtx)).Methods(http.MethodGet)

	// ==== Jobs API ====
	restAPIRouter.Path("/jobs").Handler(api.JobsPostHandler(appCtx)).Methods(http.MethodPost)
	restAPIRouter.Path("/jobs/next").Handler(api.JobsNextPostHandler(appCtx)).Methods(http.MethodPost)
	restAPIRouter.Path("/jobs/{jobID:[0-9]+}").Handler(api.JobGetHandler(appCtx)).Methods(http.MethodGet)
	restAPIRouter.Path("/jobs/{jobID:[0-9]+}/status/{jobStatus:[0-9]+}").Handler(api.JobStatusPostHandler(appCtx)).Methods(http.MethodPost)
	restAPIRouter.Path("/jobs/{jobID:[0-9]+}/log").Handler(api.JobLogPostHandler(appCtx)).Methods(http.MethodPost)
	restAPIRouter.Path("/jobs/{jobID:[0-9]+}/log.txt").Handler(api.JobLogTxtGetHandler(appCtx)).Methods(http.MethodGet)
	restAPIRouter.Path("/jobs/{jobID:[0-9]+}/artifacts").Handler(api.JobArtifactsGetHandler(appCtx)).Methods(http.MethodGet)
	restAPIRouter.Path("/jobs/{jobID:[0-9]+}/artifacts/{filename}").Handler(api.JobArtifactPostHandler(appCtx)).Methods(http.MethodPost)
	restAPIRouter.Path("/jobs/{jobID:[0-9]+}/artifacts/{filename}").Handler(api.JobArtifactGetHandler(appCtx)).Methods(http.MethodGet)

	// ==== Uploads API ====
	restAPIRouter.Path("/uploads/{uploadID:[0-9]+}").Handler(api.UploadGetHandler(appCtx)).Methods(http.MethodGet)

	// these routes are kept for dget compatibility. Dget requires that the URL ends with /<something>.dsc (TODO: open a bug)
	restAPIRouter.Path("/uploads/{uploadID:[0-9]+}/source.dsc").Handler(api.UploadDSCGetHandler(appCtx)).Methods(http.MethodGet)
	restAPIRouter.Path("/uploads/{uploadID:[0-9]+}/package.changes").Handler(api.UploadChangesGetHandler(appCtx)).Methods(http.MethodGet)

	restAPIRouter.Path("/uploads/{uploadID:[0-9]+}/dsc").Handler(api.UploadDSCGetHandler(appCtx)).Methods(http.MethodGet)
	restAPIRouter.Path("/uploads/{uploadID:[0-9]+}/changes").Handler(api.UploadChangesGetHandler(appCtx)).Methods(http.MethodGet)
	restAPIRouter.Path("/uploads/{uploadID:[0-9]+}/files").Handler(api.UploadFilesGetHandler(appCtx)).Methods(http.MethodGet)
	restAPIRouter.Path("/uploads/{uploadID:[0-9]+}/{filename}").Handler(api.UploadFileGetHandler(appCtx)).Methods(http.MethodGet)

	// ==== Artifacts API ====
	restAPIRouter.Path("/artifacts/{artifactID:[0-9]+}").Handler(api.ArtifactGetHandler(appCtx)).Methods(http.MethodGet)
	restAPIRouter.Path("/artifacts/{artifactID:[0-9]+}/content").Handler(api.ArtifactContentGetHandler(appCtx)).Methods(http.MethodGet)

	// ==== ArchiveUpgrades API ===
	restAPIRouter.Path("/archive-upgrades/{archiveUpgradeID:[0-9]+}").Handler(api.ArchiveUpgradeGetHandler(appCtx)).Methods(http.MethodGet)
	restAPIRouter.Path("/archive-upgrades/{archiveUpgradeID:[0-9]+}/jobs").Handler(api.ArchiveUpgradeJobsGetHandler(appCtx)).Methods(http.MethodGet)
	restAPIRouter.Path("/archive-upgrades/{archiveUpgradeID:[0-9]+}/successful-builds").Handler(api.ArchiveUpgradeSuccessfulBuildsGetHandler(appCtx)).Methods(http.MethodGet)

	return router
}
