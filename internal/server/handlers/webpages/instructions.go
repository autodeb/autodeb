package webpages

import (
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/middleware"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/middleware/auth"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

//InstructionsGetHandler returns a handler for the instructions page
func InstructionsGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {
		renderWithBase(r, w, appCtx, user, "instructions.html", nil)
	}

	handler := auth.MaybeWithUser(handlerFunc, appCtx)
	handler = middleware.HTMLHeaders(handler)

	return handler
}
