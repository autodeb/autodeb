package jobrunner

import (
	"fmt"
	"os"
	"path/filepath"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

// setupJob will create a working directory and a log file for a job. It is the
// caller's responsibility to delete the directory
func (jobRunner *JobRunner) setupJob(job *models.Job) (string, *os.File, error) {
	workingDirectory := filepath.Join(
		jobRunner.workingDirectory,
		fmt.Sprint(job.ID),
	)

	if err := os.Mkdir(workingDirectory, 0755); err != nil {
		return "", nil, err
	}

	logFilePath := filepath.Join(workingDirectory, "log.txt")
	logFile, err := os.Create(logFilePath)
	if err != nil {
		defer os.RemoveAll(workingDirectory)
		return "", nil, err
	}

	return workingDirectory, logFile, nil
}
