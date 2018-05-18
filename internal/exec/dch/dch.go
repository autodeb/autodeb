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
		"--distribution", distribution,
		message,
	)

	if err := command.Run(); err != nil {
		return errors.Errorf("dch error: %s", err)
	}

	return nil
}
