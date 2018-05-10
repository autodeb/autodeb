package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"

	"salsa.debian.org/autodeb-team/autodeb/internal/http/decorators"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
)

//UploadDSCGetHandler returns handler the DSC of the upload
func UploadDSCGetHandler(app *app.App) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		uploadID, err := strconv.Atoi(vars["uploadID"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		dsc, err := app.GetUploadDSC(uint(uploadID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if dsc == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		defer dsc.Close()
		io.Copy(w, dsc)
	}

	handler = decorators.TextPlainHeaders(handler)

	return http.HandlerFunc(handler)
}

//UploadFilesGetHandler returns a handler that lists all files for an upload
func UploadFilesGetHandler(app *app.App) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		uploadID, err := strconv.Atoi(vars["uploadID"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fileUploads, err := app.GetAllFileUploadsByUploadID(uint(uploadID))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(fileUploads)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jsonFileUploads := string(b)

		fmt.Fprint(w, jsonFileUploads)
	}

	handler = decorators.JSONHeaders(handler)

	return http.HandlerFunc(handler)
}

//UploadFileGetHandler returns a handler that returns upload files
func UploadFileGetHandler(app *app.App) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)

		uploadID, err := strconv.Atoi(vars["uploadID"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		filename, ok := vars["filename"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Clean the file name, keeping only the file name if it is a path
		_, filename = filepath.Split(filename)

		file, err := app.GetUploadFile(uint(uploadID), filename)
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

	return http.HandlerFunc(handler)
}
