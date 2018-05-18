//Package exec contains methods that call os commands
package exec

import (
	"context"
	"io"
	"os/exec"
)

// RunCtxDirStdoutStderr is a wrapper around exec.CommandContext and command.Run with additional arguments
func RunCtxDirStdoutStderr(ctx context.Context, directory string, stdout, stderr io.Writer, command string, args ...string) error {
	cmd := exec.CommandContext(
		ctx,
		command,
		args...,
	)
	cmd.Dir = directory
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}
