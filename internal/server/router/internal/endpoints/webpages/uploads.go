package webpages

import (
	"fmt"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/middleware"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/internal/auth"
)

//UploadsGetHandler returns a handler that renders the uploads page
func UploadsGetHandler(app *app.App) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		uploads, err := app.GetAllUploads()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data := struct {
			User    *models.User
			Uploads []*models.Upload
		}{
			User:    user,
			Uploads: uploads,
		}

		rendered, err := app.TemplatesRenderer().RenderTemplate(data, "base.html", "uploads.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, rendered)
	}

	handler := auth.MaybeWithUser(handlerFunc, app)

	handler = middleware.HTMLHeaders(handler)

	return handler
}
