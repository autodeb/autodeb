package webpages

import (
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/middleware"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/router/internal/auth"
)

//IndexGetHandler returns a handler for the main page
func IndexGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request, user *models.User) {

		serverURL := appCtx.Config().ServerURL

		serverScheme := "http"
		if serverURL.Scheme != "" {
			serverScheme = serverURL.Scheme
		}

		serverHostnamePort := serverURL.Hostname()
		if port := serverURL.Port(); port != "" {
			serverHostnamePort = serverHostnamePort + ":" + port
		}

		data := struct {
			ServerURL          string
			ServerHostnamePort string
			ServerScheme       string
		}{
			ServerURL:          serverURL.String(),
			ServerHostnamePort: serverHostnamePort,
			ServerScheme:       serverScheme,
		}

		renderWithBase(r, w, appCtx, user, "index.html", data)
	}

	handler := auth.MaybeWithUser(handlerFunc, appCtx)

	handler = middleware.HTMLHeaders(handler)

	return handler
}
