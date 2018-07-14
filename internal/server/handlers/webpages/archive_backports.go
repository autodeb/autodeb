package webpages

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/middleware"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/middleware/auth"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

//ArchiveBackportsGetHandler returns a handler that renders the archive backports page
func ArchiveBackportsGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		page := 0
		limit := 30
		if pageParam := r.URL.Query().Get("page"); pageParam != "" {
			page, _ = strconv.Atoi(pageParam)
		}

		archiveBackports, err := appCtx.JobsService().GetAllArchiveBackportsPageLimit(page, limit)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}

		data := struct {
			ArchiveBackports []*models.ArchiveBackport
			PreviousPage     int
			NextPage         int
		}{
			ArchiveBackports: archiveBackports,
			PreviousPage:     page - 1,
			NextPage:         page + 1,
		}

		renderWithBase(r, w, appCtx, user, "archive_backports.html", data)
	}

	handler := auth.MaybeWithUser(handlerFunc, appCtx)
	handler = middleware.HTMLHeaders(handler)

	return handler
}

// ArchiveBackportGetHandler returns a handler that renders the ArchiveBackport details page
func ArchiveBackportGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		vars := mux.Vars(r)
		archiveBackportID, err := strconv.Atoi(vars["archiveBackportID"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}

		archiveBackport, err := appCtx.JobsService().GetArchiveBackport(uint(archiveBackportID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}
		if archiveBackport == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		jobs, err := appCtx.JobsService().GetAllJobsByArchiveBackportID(archiveBackport.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}

		data := struct {
			ArchiveBackport *models.ArchiveBackport
			Jobs            []*models.Job
			RepositoryURL   string
		}{
			ArchiveBackport: archiveBackport,
			Jobs:            jobs,
		}

		renderWithBase(r, w, appCtx, user, "archive_backport.html", data)
	}

	handler := auth.MaybeWithUser(handlerFunc, appCtx)
	handler = middleware.HTMLHeaders(handler)

	return handler
}

// NewArchiveBackportGetHandler returns a handler that renders the page to create a new archive backport
func NewArchiveBackportGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {
		renderWithBase(r, w, appCtx, user, "archive_backport_new.html", nil)
	}

	handler := auth.WithUserOrRedirect(handlerFunc, appCtx)
	handler = middleware.HTMLHeaders(handler)

	return handler
}

// NewArchiveBackportPostHandler returns a handle that creates a new archive backport
func NewArchiveBackportPostHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		if _, err := appCtx.JobsService().CreateArchiveBackport(user.ID); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
		} else {
			appCtx.Sessions().Flash(r, w, "success", "Archive backport created successfully")
		}

		http.Redirect(w, r, "/new-archive-backport", http.StatusSeeOther)

	}

	handler := auth.WithUserOrRedirect(handlerFunc, appCtx)

	return handler
}
