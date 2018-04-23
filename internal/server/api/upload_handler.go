package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"salsa.debian.org/aviau/autodeb/internal/server/app"
)

func uploadHandler(app *app.App) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		filename := vars["filename"]

		if err := app.ProcessUpload(filename, r.Body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

	}

	return http.HandlerFunc(handler)
}
