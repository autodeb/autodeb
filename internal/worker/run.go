package worker

import (
	"fmt"
	"time"
)

func (w *Worker) run() {

	for {

		fmt.Fprintln(w.writerOutput, "Waiting 10 seconds...")
		time.Sleep(10 * time.Second)

	}

}
