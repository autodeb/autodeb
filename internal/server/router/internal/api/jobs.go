package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/decorators"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
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
