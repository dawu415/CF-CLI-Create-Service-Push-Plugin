package serviceCreatorMock

import (
	"fmt"

	"code.cloudfoundry.org/cli/plugin/models"
)

type MockCliConnection struct {
	CommandOutput []string

	GetServicesModels []plugin_models.GetServices_Model

	GetServiceExists                bool
	GetServiceModel                 plugin_models.GetService_Model
	CliCommandWasCalled             bool
	SimulateErrorOnGetServices      bool
	SimulateErrorOnGetServiceByName bool
	SimulateErrorOnCliCommand       bool
}

func NewMockCliConnection() *MockCliConnection {
	return &MockCliConnection{}
}

func (mc *MockCliConnection) CliCommandWithoutTerminalOutput(args ...string) ([]string, error) {
	return nil, nil
}
func (mc *MockCliConnection) CliCommand(args ...string) ([]string, error) {

	var err error
	mc.CliCommandWasCalled = true
	argArray := []string{}

	for _, argElement := range args {
		argArray = append(argArray, argElement)
	}
	mc.CommandOutput = argArray

	if mc.SimulateErrorOnCliCommand {
		err = fmt.Errorf("SimulateErrorOnCliCommand == true")
	}
	return argArray, err
}
func (mc *MockCliConnection) GetCurrentOrg() (plugin_models.Organization, error) {
	return plugin_models.Organization{}, nil
}
func (mc *MockCliConnection) GetCurrentSpace() (plugin_models.Space, error) {
	return plugin_models.Space{}, nil
}
func (mc *MockCliConnection) Username() (string, error)  { return "", nil }
func (mc *MockCliConnection) UserGuid() (string, error)  { return "", nil }
func (mc *MockCliConnection) UserEmail() (string, error) { return "", nil }
func (mc *MockCliConnection) IsLoggedIn() (bool, error)  { return false, nil }

// IsSSLDisabled returns true if and only if the user is connected to the Cloud Controller API with the
// `--skip-ssl-validation` flag set unless the CLI configuration file cannot be read, in which case it
// returns an error.
func (mc *MockCliConnection) IsSSLDisabled() (bool, error)         { return false, nil }
func (mc *MockCliConnection) HasOrganization() (bool, error)       { return false, nil }
func (mc *MockCliConnection) HasSpace() (bool, error)              { return false, nil }
func (mc *MockCliConnection) ApiEndpoint() (string, error)         { return "", nil }
func (mc *MockCliConnection) ApiVersion() (string, error)          { return "", nil }
func (mc *MockCliConnection) HasAPIEndpoint() (bool, error)        { return false, nil }
func (mc *MockCliConnection) LoggregatorEndpoint() (string, error) { return "", nil }
func (mc *MockCliConnection) DopplerEndpoint() (string, error)     { return "", nil }
func (mc *MockCliConnection) AccessToken() (string, error)         { return "", nil }
func (mc *MockCliConnection) GetApp(string) (plugin_models.GetAppModel, error) {
	return plugin_models.GetAppModel{}, nil
}
func (mc *MockCliConnection) GetApps() ([]plugin_models.GetAppsModel, error) {
	appModels := []plugin_models.GetAppsModel{}

	return append(appModels, plugin_models.GetAppsModel{}), nil
}
func (mc *MockCliConnection) GetOrgs() ([]plugin_models.GetOrgs_Model, error) {
	orgModels := []plugin_models.GetOrgs_Model{}

	return append(orgModels, plugin_models.GetOrgs_Model{}), nil
}
func (mc *MockCliConnection) GetSpaces() ([]plugin_models.GetSpaces_Model, error) {
	spaceModels := []plugin_models.GetSpaces_Model{}

	return append(spaceModels, plugin_models.GetSpaces_Model{}), nil
}
func (mc *MockCliConnection) GetOrgUsers(string, ...string) ([]plugin_models.GetOrgUsers_Model, error) {
	orgUsersModels := []plugin_models.GetOrgUsers_Model{}

	return append(orgUsersModels, plugin_models.GetOrgUsers_Model{}), nil
}
func (mc *MockCliConnection) GetSpaceUsers(string, string) ([]plugin_models.GetSpaceUsers_Model, error) {
	spaceUsersModels := []plugin_models.GetSpaceUsers_Model{}

	return append(spaceUsersModels, plugin_models.GetSpaceUsers_Model{}), nil
}
func (mc *MockCliConnection) GetServices() ([]plugin_models.GetServices_Model, error) {
	var err error
	if mc.SimulateErrorOnGetServices {
		err = fmt.Errorf("SimulateErrorOnGetServices = true")
	}

	return mc.GetServicesModels, err
}
func (mc *MockCliConnection) GetService(string) (plugin_models.GetService_Model, error) {

	var err error
	serviceModel := plugin_models.GetService_Model{}
	if mc.GetServiceExists {
		serviceModel = mc.GetServiceModel
	}

	if mc.SimulateErrorOnGetServiceByName {
		err = fmt.Errorf("SimulateErrorOnGetServiceByName")
	}
	return serviceModel, err
}
func (mc *MockCliConnection) GetOrg(string) (plugin_models.GetOrg_Model, error) {
	return plugin_models.GetOrg_Model{}, nil
}
func (mc *MockCliConnection) GetSpace(string) (plugin_models.GetSpace_Model, error) {
	return plugin_models.GetSpace_Model{}, nil
}
