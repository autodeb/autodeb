package worker

import (
	"fmt"
	"os"
	"path/filepath"

	"salsa.debian.org/autodeb-team/autodeb/internal/exec/dget"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func (w *Worker) execBuild(job *models.Job) error {
	workingDirectory := filepath.Join(w.workingDirectory, fmt.Sprint(job.ID))

	// Create the job directory
	if err := os.Mkdir(workingDirectory, 0755); err != nil {
		return err
	}

	// Get the .dsc URL
	dscURL := w.apiClient.GetUploadDSCURL(job.UploadID)

	// Download the source
	if err := dget.Dget(dscURL.String(), workingDirectory); err != nil {
		return err
	}

	return nil
}
