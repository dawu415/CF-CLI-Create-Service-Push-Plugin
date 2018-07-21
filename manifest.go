package main

import (
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Service struct {
	ServiceName           string            `yaml:"name"`
	Type                  string            `yaml:"type"` //brokered, credentials, drain, route.  "blank" == brokered
	Broker                string            `yaml:"broker"`
	PlanName              string            `yaml:"plan"`
	Url                   string            `yaml:"url"`
	updateSIParamsAndTags bool              `yaml:"updateSIParamsAndTags"` // Does not update service plan. This should be done manually.
	Credentials           map[string]string `yaml:"credentials"`
	Tags                  string            `yaml:"tags"`
	JSONParameters        string            `yaml:"parameters"`
}

type Manifest struct {
	Services []Service `yaml:"create-services"`
}

func ParseManifest(src io.Reader) (Manifest, error) {
	var m Manifest
	b, err := ioutil.ReadAll(src)
	if err != nil {
		return m, err
	}

	err = yaml.Unmarshal(b, &m)
	if err != nil {
		return m, err
	}

	return m, nil
}
