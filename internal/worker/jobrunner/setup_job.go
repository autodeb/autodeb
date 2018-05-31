package jobrunner

import (
	"fmt"
	"os"
	"path/filepath"

	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

type jobDirectory struct {
	parentDirectory    string
	workingDirectory   string
	artifactsDirectory string
	logFile            *os.File
}

func (jobDirectory *jobDirectory) Close() {
	defer os.RemoveAll(jobDirectory.parentDirectory)
	defer jobDirectory.logFile.Close()
}

// setupJob will create a job directory with the following layout:
//   /log.txt
//   /working-directory
//   /artifacts
// it is the caller's responsibility to call Close when done with the jobDirectory
func (jobRunner *JobRunner) setupJobDirectory(job *models.Job) (*jobDirectory, error) {
	parentDirectory := filepath.Join(jobRunner.workingDirectory, fmt.Sprint(job.ID))
	workingDirectory := filepath.Join(parentDirectory, "working-directory")
	artifactsDirectory := filepath.Join(parentDirectory, "artifacts")

	// Create the directories
	if err := os.Mkdir(parentDirectory, 0755); err != nil {
		return nil, err
	}
	if err := os.Mkdir(workingDirectory, 0755); err != nil {
		defer os.RemoveAll(parentDirectory)
		return nil, err
	}
	if err := os.Mkdir(artifactsDirectory, 0755); err != nil {
		defer os.RemoveAll(parentDirectory)
		return nil, err
	}

	logFilePath := filepath.Join(parentDirectory, "log.txt")

	// Create the log file
	logFile, err := os.Create(logFilePath)
	if err != nil {
		defer os.RemoveAll(workingDirectory)
		return nil, err
	}

	jobDirectory := &jobDirectory{
		parentDirectory:    parentDirectory,
		workingDirectory:   workingDirectory,
		artifactsDirectory: artifactsDirectory,
		logFile:            logFile,
	}

	return jobDirectory, nil
}
