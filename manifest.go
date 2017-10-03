package main

import (
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Service struct {
	ServiceName    string `yaml:"instanceName"`
	Broker         string `yaml:"brokerName"`
	PlanName       string `yaml:"planName"`
	JSONParameters string `yaml:"JSONParameters"`
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
