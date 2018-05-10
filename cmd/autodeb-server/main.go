// autodeb-server is the main server
package main

import (
	"fmt"
	"os"
	"os/signal"

	"salsa.debian.org/autodeb-team/autodeb/cmd/autodeb-server/internal/cli"
	"salsa.debian.org/autodeb-team/autodeb/internal/logo"
	"salsa.debian.org/autodeb-team/autodeb/internal/server"
)

func main() {
	// Retrieve args and Shift binary name off argument list.
	args := os.Args[1:]

	// Parse the command-line args
	cfg, err := cli.Parse(args, os.Stdout, os.Stderr)
	if err != nil {
		printErrorAndExit(err)
	}
	if cfg == nil {
		os.Exit(0)
	}

	fmt.Fprintln(os.Stdout, logo.Logo)
	fmt.Fprintf(os.Stdout, "Starting autodeb API on %v:%d.\n", cfg.HTTP.Address, cfg.HTTP.Port)

	// Start the server
	srv, err := server.New(cfg)
	if err != nil {
		printErrorAndExit(err)
	}

	// Handle SIGINT
	go func() {
		sigchan := make(chan os.Signal, 10)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		fmt.Println("\nStopping server...")
		srv.Close()
		os.Exit(0)
	}()

	// Wait for SIGINT
	select {}
}

func printErrorAndExit(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s.\n", err)
	os.Exit(1)
}
