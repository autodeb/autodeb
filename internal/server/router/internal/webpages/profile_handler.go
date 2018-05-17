package webpages

import (
	"fmt"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/middleware"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

//ProfileGetHandler returns a handler that renders the profile page
func ProfileGetHandler(app *app.App) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		user, err := app.AuthService().GetUser(r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
			return
		}

		data := struct {
			User *models.User
		}{
			User: user,
		}

		rendered, err := app.TemplatesRenderer().RenderTemplate(data, "base.html", "profile.html")
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, rendered)
	}

	handler := http.Handler(http.HandlerFunc(handlerFunc))

	handler = middleware.HTMLHeaders(handler)

	return handler
}
