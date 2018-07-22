package serviceManifest

import (
	"fmt"
	"io/ioutil"
	"os"
)

// ParseFromFilename takes in a filename and outputs a Manifest
func ParseFromFilename(filename string) (ServiceManifest, error) {

	var manifest ServiceManifest

	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		fmt.Printf("Found Service Manifest File: %s\n", filename)
		filePointer, err := os.Open(filename)
		if err == nil {
			bytes, err := ioutil.ReadAll(filePointer)
			if err != nil {
				return manifest, err
			}

			manifest, err := ParseManifest(bytes)
			if err != nil {
				return manifest, err
			}
		} else {
			return manifest, fmt.Errorf("Unable to open %s", filename)
		}
	} else {
		return manifest, fmt.Errorf("The file %s was not found", filename)
	}

	return manifest, nil
}
