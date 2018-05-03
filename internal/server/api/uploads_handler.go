package api

import (
	"fmt"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/htmltemplate"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/api/internal/decorators"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func uploadsGetHandler(renderer *htmltemplate.Renderer, app *app.App) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

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

		rendered, err := renderer.RenderTemplate(data, "base.html", "uploads.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, rendered)
	}

	handler = decorators.HTMLHeaders(handler)

	return http.HandlerFunc(handler)
}
