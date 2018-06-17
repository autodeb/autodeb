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

//UploadsGetHandler returns a handler that renders the uploads page
func UploadsGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		page := 0
		limit := 30
		if pageParam := r.URL.Query().Get("page"); pageParam != "" {
			page, _ = strconv.Atoi(pageParam)
		}

		var uploads []*models.Upload
		var err error

		userIDParam := r.URL.Query().Get("user_id")
		if userIDParam != "" {
			userID, err := strconv.Atoi(userIDParam)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				appCtx.RequestLogger().Error(r, err)
				return
			}
			uploads, err = appCtx.UploadsService().GetAllUploadsByUserIDPageLimit(uint(userID), page, limit)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				appCtx.RequestLogger().Error(r, err)
				return
			}
		} else {
			uploads, err = appCtx.UploadsService().GetAllUploadsPageLimit(page, limit)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				appCtx.RequestLogger().Error(r, err)
				return
			}
		}

		data := struct {
			Uploads      []*models.Upload
			PreviousPage int
			NextPage     int
			UserID       string
		}{
			Uploads:      uploads,
			PreviousPage: page - 1,
			NextPage:     page + 1,
			UserID:       userIDParam,
		}

		renderWithBase(r, w, appCtx, user, "uploads.html", data)
	}

	handler := auth.MaybeWithUser(handlerFunc, appCtx)

	handler = middleware.HTMLHeaders(handler)

	return handler
}

// UploadGetHandler returns a handler that renders the upload detail page
func UploadGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		vars := mux.Vars(r)
		uploadID, err := strconv.Atoi(vars["uploadID"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}

		upload, err := appCtx.UploadsService().GetUpload(uint(uploadID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}
		if upload == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		jobs, err := appCtx.JobsService().GetAllJobsByUploadID(uint(uploadID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}

		data := struct {
			Upload *models.Upload
			Jobs   []*models.Job
		}{
			Upload: upload,
			Jobs:   jobs,
		}

		renderWithBase(r, w, appCtx, user, "upload.html", data)
	}

	handler := auth.MaybeWithUser(handlerFunc, appCtx)

	handler = middleware.HTMLHeaders(handler)

	return handler
}
