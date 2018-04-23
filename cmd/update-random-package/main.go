package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"salsa.debian.org/aviau/autodeb/internal/apt"
	"salsa.debian.org/aviau/autodeb/internal/dch"
	"salsa.debian.org/aviau/autodeb/internal/sbuild"
	"salsa.debian.org/aviau/autodeb/internal/udd"
	"salsa.debian.org/aviau/autodeb/internal/uscan"
)

func main() {
	pkg := getRandomPackage()
	fmt.Printf("Selected package to update: %s\n", pkg.Package)

	workDir, err := ioutil.TempDir("", "update-random-package-")
	if err != nil {
		log.Fatal(err)
	}

	packageDir := filepath.Join(workDir, pkg.Package)
	debianDir := filepath.Join(packageDir, "debian")
	changelogPath := filepath.Join(debianDir, "changelog")

	if err := os.Mkdir(packageDir, 0700); err != nil {
		log.Fatal(err)
	}

	if err := apt.GetLatestDebianDirectory(pkg.Package, packageDir); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Debian directory unpacked in %s\n", debianDir)

	uscanResult, err := uscan.Uscan(workDir)
	if err != nil {
		log.Fatal(err)
	}
	if uscanResult.Status != uscan.ResultStatusNewerPackageAvailable {
		log.Fatal("uscan did not find a new upstream version to download\n")
	}

	fmt.Printf(
		"Current version is %s, we have downloaded %s\n",
		uscanResult.DebianUVersion,
		uscanResult.UpstreamVersion,
	)

	if err := dch.NewVersion(
		changelogPath,
		uscanResult.UpstreamVersion+"-1",
		"unstable",
		"automatic package update",
	); err != nil {
		log.Fatal(err)
	}

	if err := sbuild.Build(
		packageDir,
		os.Stdout,
		os.Stderr,
	); err != nil {
		log.Fatal(err)
	}

	outputPath := "update-random-package-output-" + pkg.Package
	if err := os.Rename(workDir, outputPath); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Build output is available at %s\n", outputPath)

}

func getRandomPackage() *udd.Package {
	packages, err := udd.PackagesWithNewerUpstreamVersions()
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().Unix())

	pkg := packages[rand.Intn(len(packages))]

	return pkg
}
