package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/middleware"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/internal/auth"
)

//UserGetHandler returns a handler returns the current user
func UserGetHandler(appCtx *appctx.Context) http.Handler {

	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		b, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonUser := string(b)
		fmt.Fprint(w, jsonUser)

	}

	handler := auth.WithUserOr403(handlerFunc, appCtx)

	handler = middleware.JSONHeaders(handler)

	return handler
}
