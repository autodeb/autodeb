package api

import (
	"fmt"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/htmltemplate"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/api/internal/decorators"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
)

func indexHandler(renderer *htmltemplate.Renderer, app *app.App) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

		rendered, err := renderer.RenderTemplate(nil, "base.html", "index.html")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, rendered)
	}

	handler = decorators.HTMLHeaders(handler)

	return http.HandlerFunc(handler)
}
