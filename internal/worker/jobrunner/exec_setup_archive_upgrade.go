package jobrunner

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"time"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/server/models"
	"salsa.debian.org/autodeb-team/autodeb/internal/udd"
)

func (jobRunner *JobRunner) execSetupArchiveUpgrade(
	ctx context.Context,
	job *models.Job,
	workingDirectory string,
	artifactsDirectory string,
	logFile io.Writer) error {

	if job.ParentType != models.JobParentTypeArchiveUpgrade {
		return errors.Errorf("unsupported parent type %s", job.ParentType)
	}

	// Get the archive upgrade
	archiveUpgrade, err := jobRunner.apiClient.GetArchiveUpgrade(job.ParentID)
	if err != nil {
		return errors.WithMessage(err, "could not get archive upgrade")
	}

	// Get all packages that need upgrading
	sourcePackages, err := udd.PackagesWithNewerUpstreamVersions()
	if err != nil {
		return errors.WithMessagef(err, "could not get source packages to update")
	}

	rand.Seed(time.Now().UnixNano())

	fmt.Fprintln(logFile, "Creating upgrade jobs...")

	for i := uint(0); i < archiveUpgrade.PackageCount; i++ {
		// select an index
		index := rand.Intn(len(sourcePackages))

		// grab it from the package list
		pkg := sourcePackages[index]

		// remove it from the package list
		sourcePackages = append(sourcePackages[:index], sourcePackages[index+1:]...)

		// Create an upgrade job
		fmt.Fprintf(
			logFile,
			"Creating job: %s\t%s => %s\n",
			pkg.Package,
			pkg.DebianUversion,
			pkg.UpstreamVersion,
		)

	}

	return nil
}
