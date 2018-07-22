package jobrunner

import (
	"context"
	"io"
	"net/http"
	"path/filepath"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/exec"
	"salsa.debian.org/autodeb-team/autodeb/internal/ftpmasterapi"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func (jobRunner *JobRunner) execBackport(
	ctx context.Context,
	job *models.Job,
	workingDirectory string,
	artifactsDirectory string,
	logFile io.Writer) error {

	sourcePackage := job.Input

	ftpmasterapiClient := ftpmasterapi.NewClient(http.DefaultClient)

	// Retrieve the testing .dsc
	dscs, err := ftpmasterapiClient.GetDSCSInSuite(sourcePackage, "testing")
	if err != nil {
		return errors.WithMessagef(err, "could not retrieve dscs for source package %s in testing", sourcePackage)
	}
	if len(dscs) < 1 {
		return errors.WithMessagef(err, "no matching dscs for package %s in testing", sourcePackage)
	}

	// Get the dsc's URL
	dscURL := ftpmasterapiClient.DSCURL(dscs[0])

	// Download the source
	if err := exec.RunCtxDirStdoutStderr(
		ctx, workingDirectory, logFile, logFile,
		"dget", "--allow-unauthenticated", dscURL,
	); err != nil {
		return errors.WithMessage(err, "dget failed")
	}

	// Obtain the unpacked source directory
	dirs, err := getDirectories(workingDirectory)
	if err != nil {
		return errors.WithMessagef(err, "could not obtain the list of directories in %s", workingDirectory)
	} else if len(dirs) > 1 {
		return errors.Errorf("too many directories, cannot guess which one contains the source: %s", dirs)
	}
	sourceDirectory := filepath.Join(workingDirectory, dirs[0])

	// Run dch
	if err := exec.RunCtxDirStdoutStderr(
		ctx, sourceDirectory, logFile, logFile,
		"dch",
		"--force-distribution",
		"--distribution=autodeb",
		"Automatic backport by Autodeb.",
	); err != nil {
		return errors.WithMessage(err, "dch failed")
	}

	// Run sbuild
	if err := exec.RunCtxDirStdoutStderr(
		ctx, sourceDirectory, logFile, logFile,
		"sbuild",
		"--no-clean-source",
		"--nolog",
		"--arch-all",
		"--source",
		"--dist=stable-backports",
	); err != nil {
		return errors.WithMessage(err, "sbuild failed")
	}

	// Find .changes file and copy all referenced files to the artifacts directory
	if changes, err := getFirstChangesInDirectory(workingDirectory); err != nil {
		return errors.WithMessage(err, "couldn't get changes file in output directory")
	} else if changes == nil {
		return errors.New("no changes file found in output directory")
	} else if err := changes.Copy(artifactsDirectory); err != nil {
		return errors.WithMessage(err, "couldn't copy changes and referenced files to the artifacts directory")
	}

	return nil
}
