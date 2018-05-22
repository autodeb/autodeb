package webpages

import (
	"fmt"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/middleware"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/internal/auth"
)

//JobsGetHandler returns a handler that renders the jobs page
func JobsGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		jobs, err := appCtx.JobsService().GetAllJobs()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data := struct {
			User *models.User
			Jobs []*models.Job
		}{
			User: user,
			Jobs: jobs,
		}

		rendered, err := appCtx.TemplatesRenderer().RenderTemplate(data, "base.html", "jobs.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Fprint(w, rendered)
	}

	handler := auth.MaybeWithUser(handlerFunc, appCtx)

	handler = middleware.HTMLHeaders(handler)

	return handler
}
