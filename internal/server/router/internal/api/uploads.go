package api

import (
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
)

//UploadGetHandler returns handler that renders an upload
func UploadGetHandler(app *app.App) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

	}

	return http.HandlerFunc(handler)
}
