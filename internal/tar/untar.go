package tar

import (
	"fmt"
	"os/exec"
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
		return fmt.Errorf("tar error: %s", err)
	}

	return nil
}
