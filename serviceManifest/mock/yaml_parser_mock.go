package parserMock

import (
	"github.com/dawu415/CF-CLI-Create-Service-Push-Plugin/serviceManifest"
)

// MockDecoder struct
type MockDecoder struct {
}

// NewMockDecoder initializes a new mock decoder
func NewMockDecoder() *MockDecoder {
	return &MockDecoder{}
}

// DecodeManifest Performs the mocked decoding
func (mock *MockDecoder) DecodeManifest(bytes []byte) (*serviceManifest.ServiceManifest, error) {
	return &serviceManifest.ServiceManifest{
		Services: []serviceManifest.Service{
			serviceManifest.Service{
				ServiceName: string(bytes[:]),
			},
		},
	}, nil
}
