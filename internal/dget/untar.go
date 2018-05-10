package dget

import (
	"fmt"
	"os/exec"
)

//Dget the URL in the directory
func Dget(url, directory string) error {
	command := exec.Command(
		"dget",
		url,
	)
	command.Dir = directory

	if err := command.Run(); err != nil {
		return fmt.Errorf("dget error: %s", err)
	}

	return nil
}
