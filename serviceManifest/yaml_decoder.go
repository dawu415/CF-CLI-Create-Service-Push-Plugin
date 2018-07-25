package serviceManifest

import yaml "gopkg.in/yaml.v2"

// DecoderInterface describes the method needed to decode a bytestream to a ServiceManifest
type DecoderInterface interface {
	DecodeManifest(bytes []byte) (*ServiceManifest, error)
}

// YmlDecoder is
type YmlDecoder struct {
}

// NewYmlDecoder initializes a new YAML Decoder
func NewYmlDecoder() *YmlDecoder {
	return &YmlDecoder{}
}

// DecodeManifest unmarshals a bytestream into a ServiceManifest struct using yaml.v2
func (yml *YmlDecoder) DecodeManifest(bytes []byte) (*ServiceManifest, error) {
	var m ServiceManifest
	err := yaml.Unmarshal(bytes, &m)
	return &m, err
}
