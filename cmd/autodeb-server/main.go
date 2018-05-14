// autodeb-server is the main server
package main

import (
	"context"
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
	cfg, err := cli.Parse(args, os.Stdout)
	if err != nil {
		printErrorAndExit(err)
	}
	if cfg == nil {
		os.Exit(0)
	}

	fmt.Fprintln(os.Stdout, logo.Logo)
	fmt.Fprintf(os.Stdout, "Starting autodeb API on %v:%d.\n", cfg.HTTP.Address, cfg.HTTP.Port)

	// Start the server
	srv, err := server.New(cfg, os.Stderr)
	if err != nil {
		printErrorAndExit(err)
	}

	// Wait for SIGINT
	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan

	// SIGINT received, shutdown...
	fmt.Println("\nShutting down the server. Send SIGINT again to force quit.")

	// Create a context, cancel it if we receive SIGINT again
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()
	go func() {
		sigchan := make(chan os.Signal)
		signal.Notify(sigchan, os.Interrupt)
		select {
		case <-sigchan:
			fmt.Println("\nForcing quit...")
			cancelCtx()
		}
	}()

	// Shutdown the server, this blocks until the shutdown is complete
	// or until we cancel the context
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println(err)
	}
}

func printErrorAndExit(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s.\n", err)
	os.Exit(1)
}
