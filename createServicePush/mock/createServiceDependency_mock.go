package createService_mock

import (
	"fmt"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/dawu415/CF-CLI-Create-Service-Push-Plugin/cspArguments"
	"github.com/dawu415/CF-CLI-Create-Service-Push-Plugin/serviceManifest"
)

type MockCreateService struct {
	ArgumentHasError      bool
	CreateParserHasError  bool
	ParseHasError         bool
	CreateServiceHasError bool
	ServicesCreated       bool
	DoNotCreateServices   bool
	DoNotPush             bool
	PlugIsUninstalling    bool
}

func NewMockCreateService() *MockCreateService {
	return &MockCreateService{}
}

func (mcsp *MockCreateService) Process(args []string) (*cspArguments.CSPArguments, error) {
	var err error
	if mcsp.ArgumentHasError {
		err = fmt.Errorf("ArgumentHasError = true")
	}
	return &cspArguments.CSPArguments{
		DoNotCreateServices:  mcsp.DoNotCreateServices,
		DoNotPush:            mcsp.DoNotPush,
		IsUninstallingPlugin: mcsp.PlugIsUninstalling,
	}, err
}

func (mcsp *MockCreateService) GetUsage() string {
	return ""
}

func (mcsp *MockCreateService) GetArgumentsDescription() map[string]string {
	return map[string]string{}
}

func (mcsp *MockCreateService) CreateServices(manifest *serviceManifest.ServiceManifest, cf plugin.CliConnection) error {

	var err error
	if mcsp.CreateServiceHasError {
		err = fmt.Errorf("CreateServiceHasError = true")
	} else {
		mcsp.ServicesCreated = true
	}
	return err
}

// Parse parses a manifest from a reader
func (mcsp *MockCreateService) Parse([]string, map[string]string) (*serviceManifest.ServiceManifest, error) {

	var err error
	if mcsp.ParseHasError {
		err = fmt.Errorf("ParseHasError = true")
	}

	return &serviceManifest.ServiceManifest{}, err
}

func (mcsp *MockCreateService) CreateParser(filename string) (*serviceManifest.ParseData, error) {

	var err error
	if mcsp.CreateParserHasError {
		err = fmt.Errorf("CreateParserHasError = true")
	}

	return &serviceManifest.ParseData{
		Parser:  mcsp,
		Reader:  nil,
		Decoder: nil,
		FileIO:  nil,
	}, err
}
