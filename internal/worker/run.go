package worker

import (
	"fmt"
	"time"
)

func (w *Worker) run() {

	for {

		fmt.Fprintln(w.writerOutput, "Waiting 10 seconds...")
		time.Sleep(10 * time.Second)

		job, err := w.apiClient.UnqueueNextJob()
		if err != nil {
			fmt.Printf("Error: could not obtain new job: %v\n", err)
		} else if job == nil {
			fmt.Println("No job available.")
		} else {
			fmt.Printf("Obtained job: %+v\n", job)

		}

	}

}
