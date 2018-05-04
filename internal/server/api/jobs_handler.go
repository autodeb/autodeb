package api

import (
	"fmt"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/htmltemplate"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/api/internal/decorators"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func jobsGetHandler(renderer *htmltemplate.Renderer, app *app.App) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

		jobs, err := app.GetAllJobs()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data := struct {
			Jobs []*models.Job
		}{
			Jobs: jobs,
		}

		rendered, err := renderer.RenderTemplate(data, "base.html", "jobs.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, rendered)
	}

	handler = decorators.HTMLHeaders(handler)

	return http.HandlerFunc(handler)
}
