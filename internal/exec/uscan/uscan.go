package uscan

import (
	"encoding/xml"
	"fmt"
	"os/exec"
)

//Uscan runs uscan
func Uscan(directory string) (*Result, error) {
	command := exec.Command(
		"uscan",
		"--dehs",
	)
	command.Dir = directory

	output, err := command.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("uscan error: %s", err)
	}

	result := &Result{}

	if err := xml.Unmarshal(output, result); err != nil {
		return nil, err
	}

	return result, nil
}
