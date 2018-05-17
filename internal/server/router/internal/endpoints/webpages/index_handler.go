package webpages

import (
	"fmt"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/middleware"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/internal/auth"
)

//IndexGetHandler returns a handler for the main page
func IndexGetHandler(app *app.App) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		data := struct {
			User      *models.User
			ServerURL string
		}{
			User:      user,
			ServerURL: app.Config().ServerURL,
		}

		rendered, err := app.TemplatesRenderer().RenderTemplate(data, "base.html", "index.html")
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
