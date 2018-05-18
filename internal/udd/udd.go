//Package udd is a client for the API at udd.debian.org
package udd

import (
	"encoding/json"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
)

const (
	upstreamStatusJSONURL = "https://udd.debian.org/cgi-bin/upstream-status.json.cgi"
)

//Package contains information about a Debian source package
type Package struct {
	Package               string `json:"package"`
	DebianMangledUversion string `json:"debian-mangled-uversion"`
	DebianUversion        string `json:"debian-uversion"`
	Status                string `json:"status"`
	UpstreamURL           string `json:"upstream-url"`
	UpstreamVersion       string `json:"upstream-version"`
	Warnings              string `json:"warnings"`
	Errors                string `json:"errors"`
}

//PackagesWithNewerUpstreamVersions returns a list of souce packages that have newer upstream versions available
func PackagesWithNewerUpstreamVersions() ([]*Package, error) {
	resp, err := http.Get(upstreamStatusJSONURL)
	if err != nil {
		return nil, errors.Errorf("getting %q: %v", upstreamStatusJSONURL, err)
	}

	if got, want := resp.StatusCode, http.StatusOK; got != want {
		return nil, errors.Errorf("unexpected HTTP status code: got %d, want %d", got, want)
	}

	var pkgs []*Package
	if err := json.NewDecoder(resp.Body).Decode(&pkgs); err != nil {
		return nil, err
	}

	var pkgsToUpdate []*Package
	for _, pkg := range pkgs {
		if pkg.Status == "newer package available" {
			pkgsToUpdate = append(pkgsToUpdate, pkg)
		}
	}

	return pkgsToUpdate, nil
}
