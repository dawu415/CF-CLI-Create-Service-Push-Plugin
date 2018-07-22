package serviceManifest

import (
	"fmt"
	"io"
	"os"
)

// Parser holds a reader that will provide the input data to be Parsed
type Parser struct {
	Reader io.Reader
}

// NewParser returns a Parser struct with a reader
func NewParser(filename string) (*Parser, error) {

	var reader io.Reader
	var err error
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		fmt.Printf("Found Service Manifest File: %s\n", filename)
		reader, err = os.Open(filename)
		if err != nil {
			err = fmt.Errorf("Unable to open %s because %s", filename, err)
		}
	} else {
		err = fmt.Errorf("The file %s was not found", filename)
	}

	return &Parser{
		Reader: reader,
	}, err
}

// Parse outputs a Service Manifest struct
func (p *Parser) Parse() (ServiceManifest, error) {
	return ParseManifest(p.Reader)
}
