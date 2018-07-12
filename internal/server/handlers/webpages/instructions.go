package webpages

import (
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/middleware"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/middleware/auth"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/services/jobs"
)

//InstructionsGetHandler returns a handler for the instructions page
func InstructionsGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		data := struct {
			RepositoryBaseURL         string
			MainUpgradeRepositoryName string
		}{
			RepositoryBaseURL:         appCtx.Config().Aptly.RepositoryBaseURL.String(),
			MainUpgradeRepositoryName: jobs.MainUpgradeRepositoryName,
		}

		renderWithBase(r, w, appCtx, user, "instructions.html", data)
	}

	handler := auth.MaybeWithUser(handlerFunc, appCtx)
	handler = middleware.HTMLHeaders(handler)

	return handler
}
