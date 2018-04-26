package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/app"
)

func uploadHandler(app *app.App) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {

		uploadParameters, err := getUploadParameters(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := app.ProcessUpload(uploadParameters, r.Body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	}

	return http.HandlerFunc(handler)
}

func getUploadParameters(r *http.Request) (*app.UploadParameters, error) {
	// Upload parameters can be set in two ways
	//
	// 1. /<filename>?param1=value1&param2=value2
	// 2. /<param1>/<value1>/<param2>/<value2>/<filename>
	//
	// Parameters set with method #1 will override those
	// set with method #2

	uploadParameters := app.UploadParameters{}

	splitPath := strings.Split(r.URL.Path, "/")
	if len(splitPath) < 1 {
		return nil, fmt.Errorf("upload parameters should atleast contain the filename")
	}

	// The file name is the last element of the path
	uploadParameters.Filename = splitPath[len(splitPath)-1]

	// Pop the filename
	splitPath = splitPath[0 : len(splitPath)-1]

	// Get parameters from method #2
	queryParams := make(map[string][]string)
	for len(splitPath) >= 2 {
		param := splitPath[0]
		value := splitPath[1]

		splitPath = splitPath[2:]

		if existingValue, ok := queryParams[param]; ok {
			queryParams[param] = append(existingValue, value)
		} else {
			queryParams[param] = []string{value}
		}
	}

	// Override with parameters from method #1
	for param, value := range r.URL.Query() {
		queryParams[param] = value
	}

	// Set the values in uploadParameters
	for param, value := range queryParams {

		switch param {
		case "forward_upload":
			if forwardUpload, err := strconv.ParseBool(value[0]); err == nil {
				uploadParameters.ForwardUpload = forwardUpload
			} else {
				return nil, fmt.Errorf("invalid value for forward_upload: %s", value[0])
			}
		default:
			return nil, fmt.Errorf("unrecognized upload parameter: %s", param)
		}

	}

	return &uploadParameters, nil
}
