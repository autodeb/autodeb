package webpages

import (
	"fmt"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/decorators"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

//UploadsGetHandler returns a handler that renders the uploads page
func UploadsGetHandler(app *app.App) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		uploads, err := app.GetAllUploads()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data := struct {
			Uploads []*models.Upload
		}{
			Uploads: uploads,
		}

		rendered, err := app.TemplatesRenderer().RenderTemplate(data, "base.html", "uploads.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, rendered)
	}

	handler := http.Handler(http.HandlerFunc(handlerFunc))

	handler = decorators.HTMLHeaders(handler)

	return handler
}
