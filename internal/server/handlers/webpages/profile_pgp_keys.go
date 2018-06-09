package webpages

import (
	"net/http"
	"strconv"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/middleware"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/middleware/auth"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

//ProfilePGPKeysGetHandler returns a handler that renders the pgp keys page
func ProfilePGPKeysGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		pgpKeys, err := appCtx.PGPService().GetUserPGPKeys(user.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}

		data := struct {
			PGPKeys                 []*models.PGPKey
			ExpectedPGPKeyProofText string
		}{
			PGPKeys:                 pgpKeys,
			ExpectedPGPKeyProofText: appCtx.PGPService().ExpectedPGPKeyProofText(user.ID),
		}

		renderWithBase(r, w, appCtx, user, "profile_pgp_keys.html", data)
	}

	handler := auth.WithUserOrRedirect(handlerFunc, appCtx)

	handler = middleware.HTMLHeaders(handler)

	return handler
}

//AddPGPKeyPostHandler returns a handler that adds PGP key to the user's profile
func AddPGPKeyPostHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}

		key := r.Form.Get("key")
		proof := r.Form.Get("proof")

		if err := appCtx.PGPService().AddUserPGPKey(user.ID, key, proof); err != nil {
			appCtx.Sessions().Flash(r, w, "danger", err.Error())
		} else {
			appCtx.Sessions().Flash(r, w, "success", "PGP key added successfully")
		}

		http.Redirect(w, r, "/profile/pgp-keys", http.StatusSeeOther)
	}

	handler := auth.WithUserOrRedirect(handlerFunc, appCtx)

	handler = middleware.HTMLHeaders(handler)

	return handler
}

//RemovePGPKeyPostHandler returns a handler that removes a PGP key
func RemovePGPKeyPostHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}

		keyIDInt, err := strconv.Atoi(r.Form.Get("keyid"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		keyID := uint(keyIDInt)

		if err := appCtx.PGPService().RemovePGPKey(keyID, user.ID); err != nil {
			appCtx.Sessions().Flash(r, w, "danger", err.Error())
		} else {
			appCtx.Sessions().Flash(r, w, "success", "PGP key removed successfully")
		}

		http.Redirect(w, r, "/profile/pgp-keys", http.StatusSeeOther)
	}

	handler := auth.WithUserOrRedirect(handlerFunc, appCtx)

	handler = middleware.HTMLHeaders(handler)

	return handler
}
