package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/middleware"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/internal/auth"
)

//JobsNextPostHandler returns a handler that find the next job to run
func JobsNextPostHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		job, err := appCtx.JobsService().UnqueueNextJob()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if job == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		b, err := json.Marshal(job)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonJob := string(b)

		fmt.Fprint(w, jsonJob)
	}

	handler := auth.WithUserOr403(handlerFunc, appCtx)

	handler = middleware.JSONHeaders(handler)

	return handler
}

//JobGetHandler returns a handler that returns a job
func JobGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

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

		b, err := json.Marshal(job)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonJob := string(b)

		fmt.Fprint(w, jsonJob)
	}

	handler := http.Handler(http.HandlerFunc(handlerFunc))

	handler = middleware.JSONHeaders(handler)

	return handler
}

//JobLogTxtGetHandler returns a handler that retrieves the log of a job
func JobLogTxtGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		// Get input values
		vars := mux.Vars(r)
		jobID, err := strconv.Atoi(vars["jobID"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		file, err := appCtx.JobsService().GetJobLog(uint(jobID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if file == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		defer file.Close()
		io.Copy(w, file)
	}

	handler := http.Handler(http.HandlerFunc(handlerFunc))

	handler = middleware.TextPlainHeaders(handler)

	return handler
}

//JobArtifactsGetHandler returns a handler that prints all artifacts of a job
func JobArtifactsGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		// Get input values
		vars := mux.Vars(r)
		jobID, err := strconv.Atoi(vars["jobID"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jobArtifacts, err := appCtx.JobsService().GetAllJobArtifactsByJobID(uint(jobID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(jobArtifacts)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonJobArtifacts := string(b)

		fmt.Fprint(w, jsonJobArtifacts)
	}

	handler := http.Handler(http.HandlerFunc(handlerFunc))

	handler = middleware.JSONHeaders(handler)

	return handler
}

//JobArtifactGetHandler returns a handler that retrieves the artifact of a job
func JobArtifactGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		// Get input values
		vars := mux.Vars(r)
		jobID, err := strconv.Atoi(vars["jobID"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		filename, ok := vars["filename"]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		file, err := appCtx.JobsService().GetJobArtifact(uint(jobID), filename)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if file == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		defer file.Close()
		io.Copy(w, file)
	}

	handler := http.Handler(http.HandlerFunc(handlerFunc))

	handler = middleware.TextPlainHeaders(handler)

	return handler
}

//JobLogPostHandler returns a handler that saves a log for a job
func JobLogPostHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		// Get input values
		vars := mux.Vars(r)
		jobID, err := strconv.Atoi(vars["jobID"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Get the job
		job, err := appCtx.JobsService().GetJob(uint(jobID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if job == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Only accept logs for jobs that are currently running
		// Completed jobs should be immutable.
		switch job.Status {
		case models.JobStatusAssigned:
		default:
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Save the logs
		if err := appCtx.JobsService().SaveJobLog(uint(jobID), r.Body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

	}

	handler := auth.WithUserOr403(handlerFunc, appCtx)

	return handler
}

//JobStatusPostHandler returns a handler that sets the job status
func JobStatusPostHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		// Get input values
		vars := mux.Vars(r)
		jobID, err := strconv.Atoi(vars["jobID"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		jobStatus, err := strconv.Atoi(vars["jobStatus"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Validate the new status
		newStatus := models.JobStatus(jobStatus)
		switch newStatus {
		case models.JobStatusSuccess:
		case models.JobStatusFailed:
		case models.JobStatusQueued:
			// Allow workers to requeue jobs that they didn't complete.
		default:
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Get the job
		job, err := appCtx.JobsService().GetJob(uint(jobID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if job == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Validate the current status. We only accept status updates on
		// jobs that were assigned
		switch job.Status {
		case models.JobStatusAssigned:
		default:
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Update the job
		job.Status = newStatus
		if err := appCtx.JobsService().UpdateJob(job); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	handler := auth.WithUserOr403(handlerFunc, appCtx)

	return handler
}

//JobArtifactPostHandler returns a handler that saves job artifacts
func JobArtifactPostHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		// Get input values
		vars := mux.Vars(r)
		jobID, err := strconv.Atoi(vars["jobID"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		filename, ok := vars["filename"]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Get the job
		job, err := appCtx.JobsService().GetJob(uint(jobID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if job == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Only accept artifacts for jobs that are currently running
		// Completed jobs should be immutable.
		switch job.Status {
		case models.JobStatusAssigned:
		default:
			w.WriteHeader(http.StatusForbidden)
			return
		}

		// Save the artifact
		if err := appCtx.JobsService().SaveJobArtifact(uint(jobID), filename, r.Body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

	}

	handler := auth.WithUserOr403(handlerFunc, appCtx)

	return handler
}
