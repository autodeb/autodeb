package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/middleware"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
)

//ArtifactGetHandler returns a handler that returns an artifact
func ArtifactGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		artifactID, err := strconv.Atoi(vars["artifactID"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		artifact, err := appCtx.ArtifactsService().GetArtifact(uint(artifactID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if artifact == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		b, err := json.Marshal(artifact)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonArtifact := string(b)

		fmt.Fprint(w, jsonArtifact)
	}

	handler := http.Handler(http.HandlerFunc(handlerFunc))

	handler = middleware.JSONHeaders(handler)

	return handler
}

//ArtifactContentGetHandler returns the content of a job's artifact
func ArtifactContentGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		// Get input values
		vars := mux.Vars(r)
		artifactID, err := strconv.Atoi(vars["artifactID"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		file, err := appCtx.ArtifactsService().GetArtifactContent(uint(artifactID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if file == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		defer file.Close()
		io.Copy(w, file)
	}

	handler := http.Handler(http.HandlerFunc(handlerFunc))

	return handler
}
