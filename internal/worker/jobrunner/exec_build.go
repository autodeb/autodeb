package jobrunner

import (
	"context"
	"io"
	"io/ioutil"
	"path/filepath"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/exec"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func (jobRunner *JobRunner) execBuild(ctx context.Context, job *models.Job, workingDirectory string, logFile io.Writer) error {
	// Get the .dsc URL
	dscURL := jobRunner.apiClient.GetUploadDSCURL(job.UploadID)

	// Download the source
	if err := exec.RunCtxDirStdoutStderr(
		ctx, workingDirectory, logFile, logFile,
		"dget", "--allow-unauthenticated", dscURL.String(),
	); err != nil {
		return errors.Errorf("dget error: %s", err)
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
	if err := exec.RunCtxDirStdoutStderr(
		ctx, sourceDirectory, logFile, logFile,
		"sbuild", "--no-clean-source", "--nolog",
	); err != nil {
		return errors.Errorf("sbuild error: %s", err)
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
