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

		pgpKeys, err := appCtx.PGPService().GetUserPGPKeys(user.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		data := struct {
			PGPKeys                 []*models.PGPKey
			ExpectedPGPKeyProofText string
		}{
			PGPKeys:                 pgpKeys,
			ExpectedPGPKeyProofText: appCtx.PGPService().ExpectedPGPKeyProofText(user.ID),
		}

		renderWithBase(r, w, appCtx, user, "profile.html", data)
	}

	handler := auth.WithUser(handlerFunc, appCtx)

	handler = middleware.HTMLHeaders(handler)

	return handler
}

//AddPGPKeyPostHandler returns a handler that adds PGP key to the user's profile
func AddPGPKeyPostHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		key := r.Form.Get("key")
		proof := r.Form.Get("proof")

		if err := appCtx.PGPService().AddUserPGPKey(user.ID, key, proof); err != nil {
			appCtx.Sessions().Flash(r, w, "danger", err.Error())
		} else {
			appCtx.Sessions().Flash(r, w, "success", "PGP key added successfully")
		}

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}

	handler := auth.WithUser(handlerFunc, appCtx)

	handler = middleware.HTMLHeaders(handler)

	return handler
}
