// Package filesystem provides an abstraction for a file system.
// It can be extended to suppot multiple backends, such as remote backends.
package filesystem

import (
	"github.com/spf13/afero"
)

// FS is an abstraction for a file system
type FS = afero.Fs

// File represents a file on the file system
type File = afero.File

//NewFS creates a new file system
func NewFS(basepath string) (FS, error) {
	// TODO: Fow now, we hardcode BasePathFs. In the future,
	// the function could accept some sort of ConnectionString
	// that can be used to select a specific backend.

	basePathFs := afero.NewBasePathFs(afero.NewOsFs(), basepath)

	return basePathFs, nil
}
