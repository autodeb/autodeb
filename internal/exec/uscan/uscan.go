package uscan

import (
	"context"
	"encoding/xml"
	"os/exec"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
)

//Uscan runs uscan from the given directory
func Uscan(ctx context.Context, directory string) (*Result, error) {
	command := exec.CommandContext(
		ctx,
		"uscan",
		"--dehs",
	)
	command.Dir = directory

	output, err := command.Output()
	if err != nil {
		combinedOutput, _ := command.CombinedOutput()
		return nil, errors.WithMessagef(err, "uscan error: %s: %s", err, combinedOutput)
	}

	result := &Result{}

	if err := xml.Unmarshal(output, result); err != nil {
		return nil, errors.WithMessagef(err, "cannot parse uscan dehs xml output: %s", output)
	}

	return result, nil
}
