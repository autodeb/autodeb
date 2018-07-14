package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/middleware"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/appctx"
)

//ArchiveUpgradeGetHandler returns a handler that returns an archive upgrade
func ArchiveUpgradeGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		archiveUpgradeID, err := strconv.Atoi(vars["archiveUpgradeID"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}

		archiveUpgrade, err := appCtx.JobsService().GetArchiveUpgrade(uint(archiveUpgradeID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}
		if archiveUpgrade == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(archiveUpgrade); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}

	}

	handler := http.Handler(http.HandlerFunc(handlerFunc))
	handler = middleware.JSONHeaders(handler)

	return handler
}

//ArchiveUpgradeJobsGetHandler returns a handler that lists all jobs for an ArchiveUpgrade
func ArchiveUpgradeJobsGetHandler(appCtx *appctx.Context) http.Handler {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		archiveUpgradeID, err := strconv.Atoi(vars["archiveUpgradeID"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}

		if archiveUpgrade, err := appCtx.JobsService().GetArchiveUpgrade(uint(archiveUpgradeID)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		} else if archiveUpgrade == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		jobs, err := appCtx.JobsService().GetAllJobsByArchiveUpgradeID(uint(archiveUpgradeID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
		}

		if err := json.NewEncoder(w).Encode(jobs); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			appCtx.RequestLogger().Error(r, err)
			return
		}

	}

	handler := http.Handler(http.HandlerFunc(handlerFunc))
	handler = middleware.JSONHeaders(handler)

	return handler
}
