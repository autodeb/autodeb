// Package cli is responsible for parsing command line arguments and creating
// a server instance.
package cli

import (
	"fmt"
	"io"

	"salsa.debian.org/autodeb-team/autodeb/internal/logo"
)

// Run reads arguments and creates an autodeb worker
func Run(args []string, writerOutput, writerError io.Writer) error {
	fmt.Fprintln(writerOutput, logo.Logo)

	//TODO: actually create a worker and return it

	return nil
}
