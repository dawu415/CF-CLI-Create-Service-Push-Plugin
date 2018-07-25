package parserMock

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// MockFileIO struct
type MockFileIO struct {
	StatError    error
	FileNotExist bool
	FileCanOpen  bool
}

// NewMockFileIO initializes a new mock decoder
func NewMockFileIO() *MockFileIO {
	return &MockFileIO{
		StatError:    nil,
		FileNotExist: false,
		FileCanOpen:  true,
	}
}

// Stat returns a FileInfo object for the given filename, name.
func (fio *MockFileIO) Stat(name string) (os.FileInfo, error) {
	return nil, fio.StatError
}

// IsNotExist returns a boolean as to weather a file exists
func (fio *MockFileIO) IsNotExist(err error) bool {
	return fio.FileNotExist
}

// OpenReadOnly Mock here gets the file name returns its prepended with an Opened_ string
func (fio *MockFileIO) OpenReadOnly(filename string) (io.Reader, error) {
	var err error
	if !fio.FileCanOpen {
		err = fmt.Errorf("File Set to Not Openable")
	}
	outputString := string("Opened_" + filename)
	return bytes.NewBufferString(outputString), err
}
