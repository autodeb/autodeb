package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/middleware"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
)

//ArchiveBackportGetHandler returns a handler that returns an ArchiveBackport
func ArchiveBackportGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		archiveBackportID, err := strconv.Atoi(vars["archiveBackportID"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}

		archiveBackport, err := appCtx.JobsService().GetArchiveBackport(uint(archiveBackportID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}
		if archiveBackport == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(archiveBackport); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}

	}

	handler := http.Handler(http.HandlerFunc(handlerFunc))
	handler = middleware.JSONHeaders(handler)

	return handler
}
