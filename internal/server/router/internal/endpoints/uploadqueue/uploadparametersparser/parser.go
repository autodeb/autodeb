package uploadparametersparser

import (
	"net/http"
	"strconv"
	"strings"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/services/uploads"
)

//Parse an http request for upload parameters
func Parse(r *http.Request) (*uploads.UploadParameters, error) {
	// Upload parameters can be set in two ways
	//
	// 1. /<filename>?param1=value1&param2=value2
	// 2. /<param1>/<value1>/<param2>/<value2>/<filename>
	//
	// Parameters set with method #1 will override those
	// set with method #2

	filename, splitPath, err := splitURLPath(r.URL.Path)
	if err != nil {
		return nil, err
	}

	uploadParameters := uploads.UploadParameters{
		Filename: filename,
	}

	// Get parameters from method #2
	parameters := getURLPathParameters(splitPath)

	// Override with parameters from method #1
	for param, value := range r.URL.Query() {
		parameters[param] = value
	}

	// Set the values in uploadParameters
	for param, value := range parameters {

		switch param {
		case "forward_upload":
			if forwardUpload, err := strconv.ParseBool(value[0]); err == nil {
				uploadParameters.ForwardUpload = forwardUpload
			} else {
				return nil, errors.Errorf("invalid value for forward_upload: %s", value[0])
			}
		default:
			return nil, errors.Errorf("unrecognized upload parameter: %s", param)
		}

	}

	return &uploadParameters, nil
}

func splitURLPath(path string) (string, []string, error) {
	splitPath := strings.Split(
		strings.Trim(path, "/"),
		"/",
	)

	if splitPath[0] == "" {
		return "", nil, errors.Errorf("upload parameters should atleast contain the filename")
	}

	// The file name is the last element of the path
	filename := splitPath[len(splitPath)-1]

	// Pop the filename
	splitPath = splitPath[0 : len(splitPath)-1]

	return filename, splitPath, nil
}

func getURLPathParameters(splitPath []string) map[string][]string {
	urlPathParameters := make(map[string][]string)

	for len(splitPath) >= 2 {
		param := splitPath[0]
		value := splitPath[1]

		splitPath = splitPath[2:]

		if existingValue, ok := urlPathParameters[param]; ok {
			urlPathParameters[param] = append(existingValue, value)
		} else {
			urlPathParameters[param] = []string{value}
		}
	}

	return urlPathParameters
}
