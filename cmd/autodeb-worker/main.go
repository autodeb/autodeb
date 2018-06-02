// autodeb-worker retrieves jobs from the main server and executes them
package main

import (
	"fmt"
	"golang.org/x/sys/unix" // "syscall" is deprecated
	"os"
	"os/signal"

	"salsa.debian.org/autodeb-team/autodeb/cmd/autodeb-worker/internal/cli"
	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/logo"
	"salsa.debian.org/autodeb-team/autodeb/internal/worker"
)

func main() {
	// Retrieve args and Shift binary name off argument list.
	args := os.Args[1:]

	// Parse the command-line args
	cfg, err := cli.Parse(args, filesystem.NewOsFS(), os.Stdout)
	if err != nil {
		printErrorAndExit(err)
	}
	if cfg == nil {
		os.Exit(0)
	}

	fmt.Fprintln(os.Stdout, logo.Logo)
	fmt.Fprintln(os.Stdout, "Starting autodeb worker.")

	// Start the server
	worker, err := worker.New(cfg, os.Stderr)
	if err != nil {
		printErrorAndExit(err)
	}

	// Wait for SIGINT/SIGTERM
	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, unix.SIGINT, unix.SIGTERM)
	<-sigchan

	fmt.Println("\nShutting down the worker.")

	if err := worker.Shutdown(); err != nil {
		fmt.Println(err)
	}
}

func printErrorAndExit(err error) {
	fmt.Fprintf(os.Stderr, "Error: %+v.\n", err)
	os.Exit(1)
}
