package jobrunner

import (
	"context"
	"io"
	"io/ioutil"
	"path/filepath"

	"salsa.debian.org/autodeb-team/autodeb/internal/exec/dget"
	"salsa.debian.org/autodeb-team/autodeb/internal/exec/sbuild"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func (jobRunner *JobRunner) execBuild(ctx context.Context, job *models.Job, workingDirectory string, logFile io.Writer) error {
	// Get the .dsc URL
	dscURL := jobRunner.apiClient.GetUploadDSCURL(job.UploadID)

	// Download the source
	if err := dget.Dget(dscURL.String(), workingDirectory); err != nil {
		return err
	}

	// Find the source directory
	dirs, err := getDirectories(workingDirectory)
	if err != nil {
		return err
	}
	if numDirs := len(dirs); numDirs != 1 {
		return err
	}
	sourceDirectory := filepath.Join(workingDirectory, dirs[0])

	// Run sbuild
	if err := sbuild.Build(ctx, sourceDirectory, logFile, logFile); err != nil {
		return err
	}

	return nil
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
