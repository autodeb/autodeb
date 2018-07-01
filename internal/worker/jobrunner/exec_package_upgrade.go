package jobrunner

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/exec"
	"salsa.debian.org/autodeb-team/autodeb/internal/exec/apt"
	"salsa.debian.org/autodeb-team/autodeb/internal/exec/dch"
	"salsa.debian.org/autodeb-team/autodeb/internal/exec/uscan"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
)

func (jobRunner *JobRunner) execPackageUpgrade(
	ctx context.Context,
	job *models.Job,
	workingDirectory string,
	artifactsDirectory string,
	logFile io.Writer) error {

	pkg := job.Input
	fmt.Fprintf(logFile, "Upgrading package %s...", pkg)

	packageDir := filepath.Join(workingDirectory, pkg)

	if err := os.Mkdir(packageDir, 0700); err != nil {
		return errors.WithMessage(err, "could not create package directory")
	}

	if err := apt.GetLatestDebianDirectory(pkg, packageDir); err != nil {
		return errors.WithMessage(err, "could not get latest debian directory")
	}

	uscanResult, err := uscan.Uscan(ctx, packageDir)
	if err != nil {
		return errors.WithMessage(err, "uscan failed")
	} else if uscanResult.Status != uscan.ResultStatusNewerPackageAvailable {
		return errors.New("uscan did not find a new upstream version to download")
	}

	fmt.Fprintf(
		logFile,
		"Current version is %s, we have downloaded %s\n",
		uscanResult.DebianUVersion,
		uscanResult.UpstreamVersion,
	)

	changelogPath := filepath.Join(packageDir, "debian", "changelog")

	if err := dch.NewVersion(
		changelogPath,
		uscanResult.UpstreamVersion+"-1",
		"unstable",
		"Update to new upstream version by autodeb.",
	); err != nil {
		return errors.WithMessage(err, "dch failed")
	}

	if err := exec.RunCtxDirStdoutStderr(
		ctx, packageDir, logFile, logFile,
		"bash",
		"-c",
		"while QUILT_PATCHES=\"debian/patches\" quilt push -f; do quilt refresh; done && quilt pop -a",
	); err != nil {
		// Exit code 2 is not an error:
		// > An exit status of 2 denotes that quilt did not
		// > do anything to complete the command.
		exitCode, err := exec.ExitCodeFromError(err)
		if err != nil {
			return errors.WithMessage(err, "could not obtain exit code")
		} else if exitCode != 2 {
			return errors.Errorf("could not refresh patches")
		}
	}

	// Run sbuild
	if err := exec.RunCtxDirStdoutStderr(
		ctx, packageDir, logFile, logFile,
		"sbuild",
		"--no-clean-source",
		"--nolog",
		"--arch-all",
		"--source",
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
