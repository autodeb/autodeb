// autodeb-server is the main server
package main

import (
	"context"
	"fmt"
	"golang.org/x/sys/unix" // "syscall" is deprecated
	"os"
	"os/signal"

	"salsa.debian.org/autodeb-team/autodeb/cmd/autodeb-server/internal/cli"
	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
	"salsa.debian.org/autodeb-team/autodeb/internal/logo"
	"salsa.debian.org/autodeb-team/autodeb/internal/server"
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
	fmt.Fprintf(os.Stdout, "Starting autodeb API on %s.\n", cfg.HTTP.Address)

	// Start the server
	srv, err := server.New(cfg, os.Stderr)
	if err != nil {
		printErrorAndExit(err)
	}

	// Wait for SIGINT/SIGTERM
	sigchan := make(chan os.Signal)
	signal.Notify(sigchan, unix.SIGINT, unix.SIGTERM)
	<-sigchan

	// SIGINT/SIGTERM received, shutdown...
	fmt.Println("\nShutting down the server. Send SIGINT/SIGTERM again to force quit.")

	// Create a context, cancel it if we receive SIGINT again
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()
	go func() {
		sigchan := make(chan os.Signal)
		signal.Notify(sigchan, unix.SIGINT, unix.SIGTERM)
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
	fmt.Fprintf(os.Stderr, "Error: %+v.\n", err)
	os.Exit(1)
}
