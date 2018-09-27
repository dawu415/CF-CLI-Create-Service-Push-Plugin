package serviceManifest

import (
	"fmt"
	"io/ioutil"

	"github.com/cloudfoundry/bosh-cli/director/template"
	yaml "gopkg.in/yaml.v2"
)

// DecoderInterface describes the method needed to decode a bytestream to a ServiceManifest
type DecoderInterface interface {
	DecodeManifest(bytes []byte, varsFilePaths []string, vars map[string]string) (*ServiceManifest, error)
}

// YmlDecoder is
type YmlDecoder struct {
}

// NewYmlDecoder initializes a new YAML Decoder
func NewYmlDecoder() *YmlDecoder {
	return &YmlDecoder{}
}

// DecodeManifest unmarshals a bytestream into a ServiceManifest struct using yaml.v2
// In addition, it will also evaluate any templated variables that are specified in the input service manifest yaml
func (yml *YmlDecoder) DecodeManifest(bytes []byte, varsFilePaths []string, vars map[string]string) (*ServiceManifest, error) {
	var m ServiceManifest
	var err error

	tpl := template.NewTemplate(bytes)
	yamlVars := template.StaticVariables{}

	for _, path := range varsFilePaths {
		rawVarsFile, ioerr := ioutil.ReadFile(path)
		if ioerr != nil {
			return nil, ioerr
		}

		var sv template.StaticVariables

		err = yaml.Unmarshal(rawVarsFile, &sv)
		if err != nil {
			return nil, fmt.Errorf("Invalid vars file %s: %s", path, err)
		}

		for k, v := range sv {
			yamlVars[k] = v
		}
	}

	for key, value := range vars {
		yamlVars[key] = value
	}

	bytes, err = tpl.Evaluate(yamlVars, nil, template.EvaluateOpts{ExpectAllKeys: true})
	if err != nil {
		return nil, fmt.Errorf("Error while trying to evaluate vars in service manifest: %s", err)
	}

	err = yaml.Unmarshal(bytes, &m)
	if err != nil {
		return nil, err
	}

	return &m, err
}
