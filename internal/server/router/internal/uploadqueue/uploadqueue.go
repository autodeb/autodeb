package uploadqueue

import (
	"encoding/json"
	"fmt"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/errorchecks"
	"salsa.debian.org/autodeb-team/autodeb/internal/http/decorators"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/internal/uploadqueue/uploadparametersparser"
)

//UploadHandler returns a handler that accepts http package uploads
func UploadHandler(app *app.App) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

		uploadParameters, err := uploadparametersparser.Parse(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		upload, err := app.ProcessUpload(uploadParameters, r.Body)

		if err != nil {
			if errorchecks.IsInputError(err) {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusCreated)

		if upload != nil {
			b, _ := json.Marshal(upload)
			jsonUpload := string(b)
			fmt.Fprint(w, jsonUpload)
		}

	}

	handler = decorators.JSONHeaders(handler)

	return http.HandlerFunc(handler)
}