package dch

import (
	"fmt"
	"os/exec"
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
		return fmt.Errorf("dch error: %s", err)
	}

	return nil
}
