package api

import (
	"fmt"
	"net/http"

	"salsa.debian.org/aviau/autodeb/internal/server/api/internal/decorators"
	"salsa.debian.org/aviau/autodeb/internal/server/app"
)

func indexHandler(app *app.App) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

		rendered, err := app.RenderTemplate("index.html", nil)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, rendered)
	}

	handler = decorators.HTMLHeaders(handler)

	return http.HandlerFunc(handler)
}
