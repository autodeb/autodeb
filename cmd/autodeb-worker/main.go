// autodeb-worker retrieves jobs from the main server and executes them
package main

import (
	"fmt"
	"os"
	"os/signal"

	"salsa.debian.org/autodeb-team/autodeb/cmd/autodeb-worker/internal/cli"
)

func main() {
	// Retrieve args and Shift binary name off argument list.
	args := os.Args[1:]

	// Run the CLI, this may return a worker instance
	if err := cli.Run(args, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s.\n", err)
		os.Exit(1)
	} else {
		// Handle SIGINT
		go func() {
			sigchan := make(chan os.Signal, 10)
			signal.Notify(sigchan, os.Interrupt)
			<-sigchan
			fmt.Println("\nStopping worker...")
			os.Exit(0)
		}()

		fmt.Println("\nDoing nothing. This is a dummy program for now.")

		// Wait for SIGINT
		select {}
	}

}
