package webpages

import (
	"fmt"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/decorators"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
)

//IndexGetHandler returns a handler for the main page
func IndexGetHandler(app *app.App) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		data := struct {
			ServerURL string
		}{
			ServerURL: app.Config().ServerURL,
		}

		rendered, err := app.TemplatesRenderer().RenderTemplate(data, "base.html", "index.html")
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
