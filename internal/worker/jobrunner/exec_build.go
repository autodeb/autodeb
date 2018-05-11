package jobrunner

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"salsa.debian.org/autodeb-team/autodeb/internal/exec/dget"
	"salsa.debian.org/autodeb-team/autodeb/internal/exec/sbuild"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func (jobRunner *JobRunner) execBuild(job *models.Job) {
	workingDirectory := filepath.Join(jobRunner.workingDirectory, fmt.Sprint(job.ID))

	// Create the job directory
	if err := os.Mkdir(workingDirectory, 0755); err != nil {
		jobRunner.submitFailure(job, err)
		return
	}

	// Get the .dsc URL
	dscURL := jobRunner.apiClient.GetUploadDSCURL(job.UploadID)

	// Download the source
	if err := dget.Dget(dscURL.String(), workingDirectory); err != nil {
		jobRunner.submitFailure(job, err)
		return
	}

	// Find the source directory
	dirs, err := getDirectories(workingDirectory)
	if err != nil {
		jobRunner.submitFailure(job, err)
		return
	}
	if numDirs := len(dirs); numDirs != 1 {
		jobRunner.submitFailure(job, err)
		return
	}
	sourceDirectory := filepath.Join(workingDirectory, dirs[0])

	// Run sbuild
	if err := sbuild.Build(sourceDirectory, os.Stdout, os.Stderr); err != nil {
		jobRunner.submitFailure(job, err)
		return
	}

	jobRunner.submitSuccess(job)
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
