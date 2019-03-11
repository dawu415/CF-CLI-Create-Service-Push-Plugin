package serviceCreator

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/dawu415/CF-CLI-Create-Service-Push-Plugin/serviceManifest"
)

// CreatorInterface shows the set of methods that describes the serviceCreator
type CreatorInterface interface {
	CreateServices(manifest *serviceManifest.ServiceManifest, cf plugin.CliConnection) error
}

// ServiceCreator describes the components required for service creation
type ServiceCreator struct {
	manifest         *serviceManifest.ServiceManifest
	cf               plugin.CliConnection
	progressReporter *ProgressReporter
}

// NewServiceCreator creates a service creator with the default progress reporter
func NewServiceCreator() *ServiceCreator {
	return &ServiceCreator{progressReporter: NewProgressReporter()}
}

// CreateServices creates the services specified by manifest via a cliConnection
func (c *ServiceCreator) CreateServices(manifest *serviceManifest.ServiceManifest, cf plugin.CliConnection) error {

	createServicesobject := &ServiceCreator{
		manifest:         manifest,
		cf:               cf,
		progressReporter: NewProgressReporter(),
	}

	return createServicesobject.createServices()
}

func (c *ServiceCreator) createServices() error {
	var err error
	// Detect the type of service and then go and create them.
	// credentials: User provided credentials service
	// drain: User provided log drain service
	// route: User provided route service
	// brokered: Brokered service.  The type field can be blank to specify this as well.
	for _, serviceObject := range c.manifest.Services {
		if serviceObject.Type == "credentials" {
			err = c.createUserProvidedCredentialsService(
				serviceObject.ServiceName,
				serviceObject.Credentials,
				serviceObject.Tags,
				serviceObject.UpdateService)
		} else if serviceObject.Type == "drain" {
			err = c.createUserProvidedLogDrainService(
				serviceObject.ServiceName,
				serviceObject.URL,
				serviceObject.Tags,
				serviceObject.UpdateService)
		} else if serviceObject.Type == "route" {
			err = c.createUserProvidedRouteService(
				serviceObject.ServiceName,
				serviceObject.URL,
				serviceObject.Tags,
				serviceObject.UpdateService)
		} else {
			if serviceObject.Type == "brokered" || serviceObject.Type == "" {
				err = c.createService(serviceObject.ServiceName,
					serviceObject.Broker,
					serviceObject.PlanName,
					serviceObject.JSONParameters,
					serviceObject.Tags,
					serviceObject.UpdateService)
			} else {
				err = fmt.Errorf("Service Type: %s unsupported", serviceObject.Type)
			}
		}

		// If we encounter any errors, quit immediately, so errors are caught early.
		if err != nil {
			fmt.Printf("Create Service Error: %+v \n", err)
			break
		}
	}

	return err
}

func (c *ServiceCreator) run(args ...string) error {
	fmt.Printf("Now Running CLI Command: %s\n", strings.Join(args, " "))
	_, err := c.cf.CliCommand(args...)
	return err
}

func (c *ServiceCreator) createUserProvidedCredentialsService(name string, credentials map[string]string, tags string, updateService bool) error {
	fmt.Printf("%s - ", name)
	var shouldUpdateService bool
	s, err := c.cf.GetServices()
	if err != nil {
		return err
	}

	for _, svc := range s {
		if svc.Name == name {
			if !updateService {
				fmt.Print("already exists...skipping creation\n")
				return nil
			}
			shouldUpdateService = true
		}
	}

	credentialsJSON, _ := json.Marshal(credentials)

	if shouldUpdateService {
		fmt.Print("user provided credential service will now be updated.\n")
		err = c.run("uups", name, "-p", string(credentialsJSON))
	} else {
		fmt.Print("will now be created as a user provided credential service.\n")
		err = c.run("cups", name, "-p", string(credentialsJSON))
	}

	return err
}

func (c *ServiceCreator) createUserProvidedRouteService(name, urlString, tags string, updateService bool) error {
	fmt.Printf("%s - ", name)
	var shouldUpdateService bool
	s, err := c.cf.GetServices()
	if err != nil {
		return err
	}

	for _, svc := range s {
		if svc.Name == name {
			if !updateService {
				fmt.Print("already exists...skipping creation\n")
				return nil
			}
			shouldUpdateService = true
		}
	}

	// Check to ensure that the url begins with HTTPS because that is the only scheme supported for now.
	urlStruct, err := url.Parse(urlString)

	if err != nil {
		return err
	}

	if strings.ToLower(urlStruct.Scheme) != "https" {
		return fmt.Errorf("route scheme not specified or unsupported. User provided route service only supports https")
	}

	if shouldUpdateService {
		fmt.Print("user provided route service will now be updated.\n")
		err = c.run("uups", name, "-r", urlString)
	} else {
		fmt.Print("will now be created as a user provided route service.\n")
		err = c.run("cups", name, "-r", urlString)
	}

	return err
}

func (c *ServiceCreator) createUserProvidedLogDrainService(name, urlString, tags string, updateService bool) error {
	fmt.Printf("%s - ", name)
	var shouldUpdateService bool
	s, err := c.cf.GetServices()
	if err != nil {
		return err
	}

	for _, svc := range s {
		if svc.Name == name {
			if !updateService {
				fmt.Print("already exists...skipping creation\n")
				return nil
			}
			shouldUpdateService = true
		}
	}

	if shouldUpdateService {
		fmt.Print("user provided log drain service will now be updated.\n")
		err = c.run("uups", name, "-l", urlString)
	} else {
		fmt.Print("will now be created as a user provided log drain service.\n")
		err = c.run("cups", name, "-l", urlString)
	}

	return err
}

func (c *ServiceCreator) createService(name, broker, plan, JSONParam, tags string, updateService bool) error {
	fmt.Printf("%s - ", name)
	var shouldUpdateService bool
	s, err := c.cf.GetServices()
	if err != nil {
		return err
	}

	for _, svc := range s {
		if svc.Name == name {
			if !updateService {
				fmt.Print("already exists...skipping creation\n")
				return nil
			}
			shouldUpdateService = true
		}
	}

	// Collect the parameters
	optionalArgs := []string{}
	if tags != "" {
		optionalArgs = append(optionalArgs, "-t")
		optionalArgs = append(optionalArgs, fmt.Sprintf("\"%s\"", tags))
	}

	if JSONParam != "" {
		optionalArgs = append(optionalArgs, "-c")
		optionalArgs = append(optionalArgs, fmt.Sprintf("%s", JSONParam))
	}

	if shouldUpdateService {
		fmt.Printf("broker service will now be updated.\n")
		err = c.run(append([]string{"update-service", name}, optionalArgs...)...)
	} else {
		fmt.Printf("will now be created as a brokered service.\n")
		err = c.run(append([]string{"create-service", broker, plan, name}, optionalArgs...)...)
	}

	if err != nil {
		return err
	}

	// Now wait for the service creation to complete.
	// We wait 'infinitely' here because we don't know how long the service
	// will take to complete. There exists some service brokers where
	// provisioning requires user ticket approval input to complete.
	for {
		service, err := c.cf.GetService(name)
		if err != nil {
			return err
		}

		c.progressReporter.Step(service.LastOperation.Description)

		if service.LastOperation.State == "succeeded" {
			break
		} else if service.LastOperation.State == "failed" {
			return fmt.Errorf(
				"error %s [status: %s]",
				service.LastOperation.Description,
				service.LastOperation.State,
			)
		}
	}

	return nil
}
