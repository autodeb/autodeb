// Package ftpmasterapi is a client for the REST API at
// api-ftp-master.debian.org
package ftpmasterapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
)

const (
	ftpMasterAPIUrl = "https://api.ftp-master.debian.org"
)

//DSC api object
type DSC struct {
	Component string `json:"component"`
	Filename  string `json:"filename"`
}

//GetDSCInSuite returns a list of DSCs matching pkg in distribution
func GetDSCInSuite(pkg, distribution string) ([]*DSC, error) {
	url := fmt.Sprintf(
		"%s/dsc_in_suite/%s/%s",
		ftpMasterAPIUrl,
		distribution,
		pkg,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Errorf("GetSourceFtpMasterAPI error: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected HTTP status code: got %d", resp.StatusCode)
	}

	var dscs []*DSC
	if err := json.NewDecoder(resp.Body).Decode(&dscs); err != nil {
		return nil, err
	}

	return dscs, nil

}
