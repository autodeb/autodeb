package sbuild

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

//Build a package
func Build(directory string, outputWriter, errorWriter io.Writer) error {
	command := exec.Command(
		"sbuild",
		"--no-clean-source",
	)
	command.Dir = directory
	command.Stdout = outputWriter
	command.Stderr = errorWriter
	command.Stdin = os.Stdin // TODO: remove this

	if err := command.Run(); err != nil {
		return fmt.Errorf("sbuild error: %s", err)
	}

	return nil
}
