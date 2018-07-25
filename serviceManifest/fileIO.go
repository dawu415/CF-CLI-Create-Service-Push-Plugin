package serviceManifest

import (
	"io"
	"os"
)

// FileIOInterface interface
type FileIOInterface interface {
	Stat(name string) (os.FileInfo, error)
	IsNotExist(err error) bool
	OpenReadOnly(filename string) (io.Reader, error)
}

// FileIO struct
type FileIO struct {
	FileIOInterface
}

// NewFileIO initializes a new mock decoder
func NewFileIO() *FileIO {
	return &FileIO{}
}

// Stat returns a FileInfo object for the given filename, name.
func (fio *FileIO) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// IsNotExist returns a boolean as to weather a file exists
func (fio *FileIO) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

// OpenReadOnly performs an os.Open on a given filename
func (fio *FileIO) OpenReadOnly(filename string) (io.Reader, error) {
	return os.Open(filename)
}
