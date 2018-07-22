package dch

import (
	"os/exec"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
)

//NewVersion creates a new version
func NewVersion(changelogPath, version, distribution, message string) error {
	command := exec.Command(
		"dch",
		"--changelog", changelogPath,
		"--newversion", version,
		"--force-distribution",
		"--distribution", distribution,
		message,
	)

	if output, err := command.CombinedOutput(); err != nil {
		return errors.WithMessagef(err, "dch failed: %s", output)
	}

	return nil
}
