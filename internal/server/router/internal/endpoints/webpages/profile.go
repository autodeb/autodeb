package webpages

import (
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/middleware"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/internal/auth"
)

//ProfileGetHandler returns a handler that renders the profile page
func ProfileGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {
		renderWithBase(r, w, appCtx, user, "profile.html", nil)
	}

	handler := auth.WithUserOrRedirect(handlerFunc, appCtx)

	handler = middleware.HTMLHeaders(handler)

	return handler
}
