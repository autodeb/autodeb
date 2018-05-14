package sbuild

import (
	"context"
	"io"
	"os"
	"os/exec"
)

//Build a package
func Build(ctx context.Context, directory string, outputWriter, errorWriter io.Writer) error {
	command := exec.CommandContext(
		ctx,
		"sbuild",
		"--no-clean-source",
	)
	command.Dir = directory
	command.Stdout = outputWriter
	command.Stderr = errorWriter
	command.Stdin = os.Stdin // TODO: remove this

	if err := command.Run(); err != nil {
		return err
	}

	return nil
}
