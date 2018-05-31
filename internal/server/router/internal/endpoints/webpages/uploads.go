package webpages

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/middleware"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/internal/auth"
)

//UploadsGetHandler returns a handler that renders the uploads page
func UploadsGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		uploads, err := appCtx.UploadsService().GetAllUploads()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data := struct {
			Uploads []*models.Upload
		}{
			Uploads: uploads,
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
			return
		}

		upload, err := appCtx.UploadsService().GetUpload(uint(uploadID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jobs, err := appCtx.JobsService().GetAllJobsByUploadID(uint(uploadID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
