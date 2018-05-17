package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/middleware"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

//JobsNextPostHandler returns a handler that find the next job to run
func JobsNextPostHandler(app *app.App) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		job, err := app.UnqueueNextJob()
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

	handler := http.Handler(http.HandlerFunc(handlerFunc))

	handler = middleware.JSONHeaders(handler)

	return handler
}

//JobGetHandler returns a handler that returns a job
func JobGetHandler(app *app.App) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		jobID, err := strconv.Atoi(vars["jobID"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		job, err := app.GetJob(uint(jobID))
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
func JobLogTxtGetHandler(app *app.App) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		// Get input values
		vars := mux.Vars(r)
		jobID, err := strconv.Atoi(vars["jobID"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		file, err := app.GetJobLog(uint(jobID))
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
func JobLogPostHandler(app *app.App) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		// Get input values
		vars := mux.Vars(r)
		jobID, err := strconv.Atoi(vars["jobID"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Get the job
		job, err := app.GetJob(uint(jobID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if job == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Only accept logs on jobs that are completed
		switch job.Status {
		case models.JobStatusSuccess:
		case models.JobStatusFailed:
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Save the logs
		if err := app.SaveJobLog(uint(jobID), r.Body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	return http.HandlerFunc(handlerFunc)
}

//JobStatusPostHandler returns a handler that sets the job status
func JobStatusPostHandler(app *app.App) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

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
			// The job was a success.
		case models.JobStatusFailed:
			// The job failed.
		case models.JobStatusQueued:
			// Allow workers to requeue jobs that they didn't complete.
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Get the job
		job, err := app.GetJob(uint(jobID))
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
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Update the job
		job.Status = newStatus
		if err := app.UpdateJob(job); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	return http.HandlerFunc(handlerFunc)
}
