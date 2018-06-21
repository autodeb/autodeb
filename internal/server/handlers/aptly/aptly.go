package aptly

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/handlers/middleware/auth"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// Handler returns a reverse proxy for the aptly api
func Handler(appCtx *appctx.Context, prefix string) http.Handler {
	reverseProxy := getReverseProxy(prefix, &appCtx.Config().Aptly.APIURL.URL)

	userHandlerFunc := func(w http.ResponseWriter, r *http.Request, _ *models.User) {
		reverseProxy.ServeHTTP(w, r)
	}

	handler := auth.WithUserOr403(userHandlerFunc, appCtx)

	return handler
}

func getReverseProxy(prefix string, url *url.URL) http.Handler {
	director := func(r *http.Request) {
		r.URL.Scheme = url.Scheme
		r.URL.Host = url.Host

		originalPath := path.Clean(
			path.Join(
				"/",
				r.URL.Path,
			),
		)
		targetPath := path.Clean(
			path.Join(
				url.Path,
				strings.TrimPrefix(originalPath, prefix),
			),
		)

		r.URL.Path = targetPath
	}

	// TODO:
	//  - implement ModifyResponse to add the prefix to redirects
	//  - implement Transport so that we can support unix sockets?

	reverseProxy := &httputil.ReverseProxy{
		Director: director,
	}

	return reverseProxy
}
