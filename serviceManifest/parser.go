package serviceManifest

import (
	"fmt"
	"io"
	"io/ioutil"
)

// ParserInterface is an interface describing the default methods used to decode a manifest file
type ParserInterface interface {
	Parse(varsFilePaths []string, vars map[string]string) (*ServiceManifest, error)
	CreateParser(filename string) (*ParseData, error)
}

// ParseData holds the Parser reader and the interface that will provide the methods to process the
// input data
type ParseData struct {
	Parser  ParserInterface
	Reader  io.Reader
	Decoder DecoderInterface
	FileIO  FileIOInterface
}

// NewParser returns a ParseData structure with the default interfaces described in its struct
func NewParser() *ParseData {
	return &ParseData{
		Decoder: NewYmlDecoder(),
		FileIO:  NewFileIO(),
	}
}

// CreateParser returns a Parser struct with a reader
func (p *ParseData) CreateParser(filename string) (*ParseData, error) {
	var reader io.Reader
	var err error
	if _, err = p.FileIO.Stat(filename); !p.FileIO.IsNotExist(err) {
		fmt.Printf("Found Service Manifest File: %s\n", filename)
		reader, err = p.FileIO.OpenReadOnly(filename)
		if err != nil {
			err = fmt.Errorf("Unable to open %s because %s", filename, err)
		}
	} else {
		err = fmt.Errorf("The file %s was not found", filename)
	}

	p.Parser = p
	p.Reader = reader
	return p, err
}

// Parse parses a manifest from a reader
func (p *ParseData) Parse(varsFilePaths []string, vars map[string]string) (*ServiceManifest, error) {
	bytes, err := ioutil.ReadAll(p.Reader)
	if err != nil {
		return nil, err
	}

	return p.Decoder.DecodeManifest(bytes, varsFilePaths, vars)
}
