// Package filesystem provides an abstraction for a file system.
// It can be extended to suppot multiple backends, such as remote backends.
package filesystem

import (
	"net/http"

	"github.com/spf13/afero"
)

// FS is an abstraction for a file system
type FS = afero.Fs

// File represents a file on the file system
type File = afero.File

//NewOsFS returns a file system that uses the os
func NewOsFS() FS {
	return afero.NewOsFs()
}

//NewBasePathFS returns a file system based on the given path
func NewBasePathFS(fs FS, basepath string) FS {
	return afero.NewBasePathFs(fs, basepath)
}

//NewMemMapFS creates a memory-based file system
func NewMemMapFS() FS {
	return afero.NewMemMapFs()
}

//NewHTTPFS creates a net/http compatible filesystem from the source
func NewHTTPFS(source FS) http.FileSystem {
	return afero.NewHttpFs(source)
}
