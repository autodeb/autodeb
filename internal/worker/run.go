package worker

import (
	"fmt"
	"time"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
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
			continue
		}
		if job == nil {
			fmt.Fprintf(w.writerOutput, "No job available.\n")
			continue
		}
		fmt.Fprintf(w.writerOutput, "Obtained job: %+v\n", job)

		// Execute the job
		switch job.Type {
		case models.JobTypeBuild:
			if err := w.execBuild(job); err != nil {
				fmt.Fprintf(w.writerOutput, "Job execution error: %v\n", err)
			}
		default:
			fmt.Fprintf(w.writerOutput, "Unknown job type: %s\n", job.Type)
		}
	}

}
