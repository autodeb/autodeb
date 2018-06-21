package aptly

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"
)

// Handler returns a reverse proxy for the aptly api
func Handler(apiURL *url.URL, prefix string) http.Handler {

	director := func(r *http.Request) {
		r.URL.Scheme = apiURL.Scheme
		r.URL.Host = apiURL.Host

		originalPath := path.Clean(
			path.Join(
				"/",
				r.URL.Path,
			),
		)
		targetPath := path.Clean(
			path.Join(
				apiURL.Path,
				strings.TrimPrefix(originalPath, prefix),
			),
		)

		r.URL.Path = targetPath
	}

	// TODO:
	//  - implement ModifyResponse to add the prefix to redirects
	//  - implement Transport so that we can support unix sockets?

	handler := &httputil.ReverseProxy{
		Director: director,
	}

	return handler
}
