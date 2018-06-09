package webpages

import (
	"fmt"
	"net/http"
	"strconv"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/middleware"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/middleware/auth"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

//ProfileAccessTokensGetHandler returns a handler that renders the access tokens page
func ProfileAccessTokensGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {
		accessTokens, err := appCtx.TokensService().GetUserTokens(user.ID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}

		data := struct {
			AccessTokens []*models.AccessToken
		}{
			AccessTokens: accessTokens,
		}
		renderWithBase(r, w, appCtx, user, "profile_access_tokens.html", data)
	}

	handler := auth.WithUserOrRedirect(handlerFunc, appCtx)

	handler = middleware.HTMLHeaders(handler)

	return handler
}

//CreateAccessTokenPostHandler returns a handler that creates a new access token for the user
func CreateAccessTokenPostHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}

		name := r.Form.Get("name")

		if token, err := appCtx.TokensService().CreateToken(user.ID, name); err != nil {
			appCtx.Sessions().Flash(r, w, "danger", err.Error())
		} else {
			msg := fmt.Sprintf("Your access token is %s", token.Token)
			appCtx.Sessions().Flash(r, w, "success", msg)
			appCtx.Sessions().Flash(r, w, "success", "This is the last time that your token will be displayed to you.")
		}

		http.Redirect(w, r, "/profile/access-tokens", http.StatusSeeOther)
	}

	handler := auth.WithUserOrRedirect(handlerFunc, appCtx)

	handler = middleware.HTMLHeaders(handler)

	return handler
}

//RemoveAccessTokenPostHandler returns a handler that removes an access token
func RemoveAccessTokenPostHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}

		tokenIDInt, err := strconv.Atoi(r.Form.Get("tokenid"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tokenID := uint(tokenIDInt)

		if err := appCtx.TokensService().RemoveToken(tokenID, user.ID); err != nil {
			appCtx.Sessions().Flash(r, w, "danger", err.Error())
		} else {
			appCtx.Sessions().Flash(r, w, "success", "Access token removed successfully")
		}

		http.Redirect(w, r, "/profile/access-tokens", http.StatusSeeOther)
	}

	handler := auth.WithUserOrRedirect(handlerFunc, appCtx)

	handler = middleware.HTMLHeaders(handler)

	return handler
}
