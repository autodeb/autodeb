package tar

import (
	"os/exec"

	"salsa.debian.org/autodeb-team/autodeb/internal/errors"
)

//Untar the file in the directory
func Untar(filename, directory string) error {
	command := exec.Command(
		"tar",
		"xf",
		filename,
	)
	command.Dir = directory

	if err := command.Run(); err != nil {
		return errors.Errorf("tar error: %s", err)
	}

	return nil
}
