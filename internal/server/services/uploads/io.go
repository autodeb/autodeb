package uploads

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"salsa.debian.org/autodeb-team/autodeb/internal/filesystem"
)

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
