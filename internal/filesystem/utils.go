package filesystem

import (
	"github.com/spf13/afero"
)

//ReadFile is the equivalent of ioutil.ReadFile
func ReadFile(fs FS, filename string) ([]byte, error) {
	return afero.ReadFile(fs, filename)
}
