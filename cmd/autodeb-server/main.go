package main

import (
	"fmt"
	"os"
	"os/signal"

	"salsa.debian.org/autodeb-team/autodeb/cmd/autodeb-server/internal/cli"
)

func main() {
	// Retrieve args and Shift binary name off argument list.
	args := os.Args[1:]

	// Run the CLI, this may return a server instance
	if srv, err := cli.Run(args, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s.\n", err)
		os.Exit(1)
	} else if srv != nil {

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

}
