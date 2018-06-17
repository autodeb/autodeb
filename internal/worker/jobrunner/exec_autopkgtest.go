package jobrunner

import (
	"context"
	"fmt"
	"io"
	"os"
	osexec "os/exec"
	"path/filepath"
	"strconv"
	"syscall"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/exec"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func (jobRunner *JobRunner) execAutopkgtest(
	ctx context.Context,
	job *models.Job,
	workingDirectory string,
	artifactsDirectory string,
	logFile io.Writer) error {

	// The job input is a build job
	buildJobID, err := strconv.Atoi(job.Input)
	if err != nil {
		return errors.WithMessage(err, "could not convert input to int")
	}

	// Get the artifacts of the build job
	artifacts, err := jobRunner.apiClient.GetJobArtifacts(uint(buildJobID))
	if err != nil {
		return errors.WithMessage(err, "could not get job artifacts")
	}

	// Autopkgtest input files
	var inputFiles []string

	// Download all artifacts from the build
	for _, artifact := range artifacts {

		switch filepath.Ext(artifact.Filename) {
		case ".deb", ".dsc":
			inputFiles = append(inputFiles, artifact.Filename)
		default:
			// Not an input file but we still need to download it.
		}

		// Get the artifact
		artifactContent, err := jobRunner.apiClient.GetArtifactContent(artifact.ID)
		if err != nil {
			return errors.WithMessage(err, "could not get the artifact content")
		}

		// Write it on disk
		debPath := filepath.Join(workingDirectory, artifact.Filename)
		deb, err := os.Create(debPath)
		if err != nil {
			return errors.WithMessagef(err, "could not create artifact file %s", debPath)
		}
		defer deb.Close()
		if _, err := io.Copy(deb, artifactContent); err != nil {
			return errors.WithMessage(err, "could not copy artifact content to file")
		}

	}

	args := []string{
		"--no-built-binaries",
	}
	args = append(
		args,
		inputFiles...,
	)
	args = append(
		args,
		"--",
		"schroot",
		"unstable-amd64-sbuild",
	)

	// Run Autopkgtest
	if err := exec.RunCtxDirStdoutStderr(
		ctx, workingDirectory, logFile, logFile,
		"autopkgtest", args...,
	); err != nil {

		exitCode := -1

		// Attempt to find the exit code
		if exitError, ok := err.(*osexec.ExitError); ok {
			if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exitCode = waitStatus.ExitStatus()
			}
		}

		switch exitCode {
		case 8:
			fmt.Fprintf(logFile, "autopkgtest exited with exit code %d\n", exitCode)
		case -1:
			return errors.New("autopkgtest failed and we could not find the exit code")
		default:
			return errors.WithMessagef(err, "autopkgtest failed with exit code %d", exitCode)
		}

	}

	return nil
}
