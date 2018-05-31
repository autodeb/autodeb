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

//JobsGetHandler returns a handler that renders the jobs page
func JobsGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		jobs, err := appCtx.JobsService().GetAllJobs()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data := struct {
			Jobs []*models.Job
		}{
			Jobs: jobs,
		}

		renderWithBase(r, w, appCtx, user, "jobs.html", data)
	}

	handler := auth.MaybeWithUser(handlerFunc, appCtx)

	handler = middleware.HTMLHeaders(handler)

	return handler
}

// JobGetHandler returns a handler that renders the job detail page
func JobGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		vars := mux.Vars(r)
		jobID, err := strconv.Atoi(vars["jobID"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		job, err := appCtx.JobsService().GetJob(uint(jobID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if job == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		artifacts, err := appCtx.JobsService().GetAllJobArtifactsByJobID(uint(jobID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data := struct {
			Job       *models.Job
			Artifacts []*models.JobArtifact
		}{
			Job:       job,
			Artifacts: artifacts,
		}

		renderWithBase(r, w, appCtx, user, "job.html", data)
	}

	handler := auth.MaybeWithUser(handlerFunc, appCtx)

	handler = middleware.HTMLHeaders(handler)

	return handler
}
