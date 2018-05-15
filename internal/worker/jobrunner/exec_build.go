package jobrunner

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"salsa.debian.org/autodeb-team/autodeb/internal/exec/dget"
	"salsa.debian.org/autodeb-team/autodeb/internal/exec/sbuild"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func (jobRunner *JobRunner) execBuild(ctx context.Context, job *models.Job) {
	workingDirectory := filepath.Join(jobRunner.workingDirectory, fmt.Sprint(job.ID))
	defer os.RemoveAll(workingDirectory)

	// Create the job directory
	if err := os.Mkdir(workingDirectory, 0755); err != nil {
		jobRunner.submitFailure(ctx, job, err)
		return
	}

	// Create log file
	logFilePath := filepath.Join(workingDirectory, "build.log")
	logFile, err := os.Create(logFilePath)
	if err != nil {
		jobRunner.submitFailure(ctx, job, err)
		return
	}
	defer logFile.Close()

	// Get the .dsc URL
	dscURL := jobRunner.apiClient.GetUploadDSCURL(job.UploadID)

	// Download the source
	if err := dget.Dget(dscURL.String(), workingDirectory); err != nil {
		jobRunner.submitFailure(ctx, job, err)
		return
	}

	// Find the source directory
	dirs, err := getDirectories(workingDirectory)
	if err != nil {
		jobRunner.submitFailure(ctx, job, err)
		return
	}
	if numDirs := len(dirs); numDirs != 1 {
		jobRunner.submitFailure(ctx, job, err)
		return
	}
	sourceDirectory := filepath.Join(workingDirectory, dirs[0])

	// Run sbuild
	if err := sbuild.Build(ctx, sourceDirectory, logFile, logFile); err != nil {
		jobRunner.submitFailure(ctx, job, err)
		return
	}

	jobRunner.submitSuccess(ctx, job)
}

//getDirectories returns a list of all directories in a directory
func getDirectories(dir string) ([]string, error) {
	var directories []string

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.Mode().IsDir() {
			directories = append(directories, file.Name())
		}
	}

	return directories, nil
}
