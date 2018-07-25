package serviceManifest

// Service describes a CF service that will be instantiated
type Service struct {
	ServiceName    string            `yaml:"name"`
	Type           string            `yaml:"type"` //brokered, credentials, drain, route.  "blank" == brokered
	Broker         string            `yaml:"broker"`
	PlanName       string            `yaml:"plan"`
	URL            string            `yaml:"url"`
	UpdateService  bool              `yaml:"updateService"` // Does not update service plan. This should be done manually.
	Credentials    map[string]string `yaml:"credentials"`
	Tags           string            `yaml:"tags"`
	JSONParameters string            `yaml:"parameters"`
}

// ServiceManifest describes a service Manifest as an array of services
type ServiceManifest struct {
	Services []Service `yaml:"create-services"`
}
