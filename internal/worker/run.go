package worker

import (
	"fmt"
	"time"
)

func (w *Worker) run() {

	for {

		// Wait a little bit before asking for a new job
		fmt.Fprintln(w.writerOutput, "Waiting 10 seconds...")
		time.Sleep(10 * time.Second)

		// Get a new job
		job, err := w.apiClient.UnqueueNextJob()
		if err != nil {
			fmt.Fprintf(w.writerOutput, "Error: could not obtain new job: %v\n", err)
		} else if job == nil {
			fmt.Fprintf(w.writerOutput, "No job available.\n")
		} else {
			fmt.Fprintf(w.writerOutput, "Obtained job: %+v\n", job)
		}

	}

}
