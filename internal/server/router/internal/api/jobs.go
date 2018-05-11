package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/decorators"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

//JobsNextPostHandler returns a handler that find the next job to run
func JobsNextPostHandler(app *app.App) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

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

	handler = decorators.JSONHeaders(handler)

	return http.HandlerFunc(handler)
}

//JobGetHandler returns a handler that returns a job
func JobGetHandler(app *app.App) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

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

	handler = decorators.JSONHeaders(handler)

	return http.HandlerFunc(handler)
}

//JobStatusPostHandler returns a handler that sets the job status
func JobStatusPostHandler(app *app.App) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

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
			break
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

		// Update the job
		job.Status = newStatus
		if err := app.UpdateJob(job); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	return http.HandlerFunc(handler)
}
