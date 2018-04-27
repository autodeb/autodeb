package app

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"salsa.debian.org/autodeb-team/autodeb/internal/crypto/sha256"
	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
)

// UploadParameters defines upload behaviour
type UploadParameters struct {
	Filename      string
	ForwardUpload bool
}

// ProcessUpload receives uploaded files
func (app *App) ProcessUpload(uploadParameters *UploadParameters, uploadContent io.Reader) error {
	// Clean the file name, ensure that it contains only a file name and
	// that it isn't something shady like ../../filename.txt
	_, uploadFileName := filepath.Split(uploadParameters.Filename)

	// Check if this is a .changes upload
	isChanges := strings.HasSuffix(uploadFileName, ".changes")

	// Save the upload to a temp file on the os's filesystem so that we can
	// calculate the shasum
	tmpfileName, err := writeToTempfile(uploadContent)
	if err != nil {
		return err
	}
	defer os.Remove(tmpfileName)

	// Find out the destination directory
	var destDir string
	if isChanges {

		if upload, err := app.dataStore.CreateUpload(); err == nil {
			destDir = filepath.Join(app.UploadsDirectory(), fmt.Sprint(upload.ID))
		} else {
			return err
		}

	} else {

		if shasum, err := sha256.Sum256HexFile(tmpfileName); err == nil {
			destDir = filepath.Join(app.UploadedFilesDirectory(), shasum)
		} else {
			return err
		}

	}

	// Open the temporary file
	tmpFile, err := os.Open(tmpfileName)
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	// Write the upload in the fileStorage
	err = writeDataToDestInFS(
		tmpFile,
		destDir,
		uploadFileName,
		app.dataFS,
	)

	return err
}

func writeDataToDestInFS(data io.Reader, destDir, destFileName string, fs filesystem.FS) error {
	// Create the destination directory
	if err := fs.MkdirAll(destDir, 0744); err != nil {
		return err
	}

	// Create the destination file
	destFile, err := fs.Create(filepath.Join(destDir, destFileName))
	if err != nil {
		fs.RemoveAll(destDir)
		return err
	}
	defer destFile.Close()

	// Write data
	if _, err := io.Copy(destFile, data); err != nil {
		destFile.Close()
		fs.RemoveAll(destDir)
	}

	return nil
}

func writeToTempfile(data io.Reader) (string, error) {
	// Create temp file
	tmpfile, err := ioutil.TempFile("", "autodeb-upload")
	if err != nil {
		return "", err
	}
	defer tmpfile.Close()

	filename := tmpfile.Name()

	// Write data
	_, err = io.Copy(tmpfile, data)
	if err != nil {
		tmpfile.Close()
		os.Remove(filename)
		return "", err
	}

	return filename, nil
}
