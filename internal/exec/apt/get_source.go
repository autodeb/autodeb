package apt

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
	"salsa.debian.org/autodeb-team/autodeb/internal/exec/tar"
	"salsa.debian.org/autodeb-team/autodeb/internal/ftpmasterapi"
)

const (
	debianMirrorURL = "https://deb.debian.org/debian"
)

//getSourceFtpMasterAPI tries to mimic apt-get source but by using the
//ftpmasterapi and dget.
func getSourceFtpMasterAPI(pkg, distribution, directory string) error {
	dscs, err := ftpmasterapi.GetDSCInSuite(pkg, distribution)
	if err != nil {
		return errors.Errorf("ftpmasterapi error: %s", err)
	}

	numberOfDSC := len(dscs)
	if numberOfDSC < 1 {
		return errors.Errorf("expected at least one dsc, got none")
	}

	dsc := dscs[numberOfDSC-1]

	dscURL := fmt.Sprintf(
		"%s/pool/%s/%s",
		debianMirrorURL,
		dsc.Component,
		dsc.Filename,
	)

	command := exec.Command("dget", dscURL)
	command.Dir = directory

	if err := command.Run(); err != nil {
		return errors.Errorf("dget error: %s", err)
	}

	return nil
}

//GetLatestDebianDirectory will download and unpack the latest debian
//directory available for the specified package in the destination folder.
func GetLatestDebianDirectory(pkg, destDir string) error {
	//Create a directory to download the source
	aptGetSourceDir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	defer os.RemoveAll(aptGetSourceDir)

	//Download the latest source from unstable
	if err := getSourceFtpMasterAPI(pkg, "unstable", aptGetSourceDir); err != nil {
		return err
	}

	//Find the debian tarball and unpack it
	debianTarballPath, err := findDebianTarballInSourceDirectory(aptGetSourceDir)
	if err != nil {
		return err
	}

	err = tar.Untar(debianTarballPath, destDir)

	return err
}

//findDebianTarballInSourceDirectory returns the path of the debian tarball
//where apt-get source was run.
func findDebianTarballInSourceDirectory(directory string) (string, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		filename := file.Name()
		if strings.Contains(filename, ".debian.") {
			fullPath := filepath.Join(directory, filename)
			return fullPath, nil
		}
	}

	return "", errors.Errorf("could not find debian tarball in %s", directory)
}
