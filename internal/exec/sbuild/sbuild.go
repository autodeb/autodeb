package sbuild

import (
	"context"
	"io"
	"os/exec"
)

//Build a package
func Build(ctx context.Context, directory string, outputWriter, errorWriter io.Writer) error {
	command := exec.CommandContext(
		ctx,
		"sbuild",
		"--no-clean-source",
		"--nolog",
	)
	command.Dir = directory
	command.Stdout = outputWriter
	command.Stderr = errorWriter

	if err := command.Run(); err != nil {
		return err
	}

	return nil
}
