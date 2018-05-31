package jobrunner

import (
	"context"
	"io"
	"io/ioutil"
	"path/filepath"

	"pault.ag/go/debian/control"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/exec"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func (jobRunner *JobRunner) execBuild(
	ctx context.Context,
	job *models.Job,
	workingDirectory string,
	artifactsDirectory string,
	logFile io.Writer) error {

	// Get the .dsc URL
	dscURL := jobRunner.apiClient.GetUploadDSCURL(job.UploadID)

	// Download the source
	if err := exec.RunCtxDirStdoutStderr(
		ctx, workingDirectory, logFile, logFile,
		"dget", "--allow-unauthenticated", dscURL.String(),
	); err != nil {
		return errors.WithMessage(err, "dget failed")
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
		return errors.WithMessage(err, "sbuild failed")
	}

	// Find .changes file
	changes, err := getFirstChangesInDirectory(workingDirectory)
	if err != nil {
		return errors.WithMessage(err, "couldn't get changes file in output directory")
	}
	if changes == nil {
		return errors.New("no changes file found in output directory")
	}

	// Move .changes and referenced files to artifacts directory
	if err := changes.Move(artifactsDirectory); err != nil {
		return errors.WithMessage(err, "couldn't move build output to artifacts directory")
	}

	return nil
}

// getFirstChangesInDirectory returns the first .changes file found in a directory
func getFirstChangesInDirectory(directory string) (*control.Changes, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {

		if filepath.Ext(file.Name()) != ".changes" {
			continue
		}

		changes, err := control.ParseChangesFile(
			filepath.Join(directory, file.Name()),
		)
		return changes, err

	}

	return nil, nil
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
