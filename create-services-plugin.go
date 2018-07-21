package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
)

// CreateServicePush is the struct implementing the interface defined by the core CLI. It can
// be found at  "code.cloudfoundry.org/cli/plugin/plugin.go"
type CreateServicePush struct {
	manifest *Manifest
	cf       plugin.CliConnection
}

// Run must be implemented by any plugin because it is part of the
// plugin interface defined by the core CLI.
//
// Run(....) is the entry point when the core CLI is invoking a command defined
// by a plugin. The first parameter, plugin.CliConnection, is a struct that can
// be used to invoke cli commands. The second paramter, args, is a slice of
// strings. args[0] will be the name of the command, and will be followed by
// any additional arguments a cli user typed in.
//
// Any error handling should be handled with the plugin itself (this means printing
// user facing errors). The CLI will exit 0 if the plugin exits 0 and will exit
// 1 should the plugin exits nonzero.
func (c *CreateServicePush) Run(cliConnection plugin.CliConnection, args []string) {

	if args[0] != "create-service-push" {
		return
	}

	args = args[1:] // Remove the create-service-push

	// 1. Process the input arguments specific to create-service-push
	var manifestFilename = "services-manifest.yml"
	processedFlags := map[string]bool{}
	remainingArgs := []string{}

	for argIdx := 0; argIdx < len(args); argIdx++ {
		switch args[argIdx] {
		case "--service-manifest": // Specify the service manifest to use
			if (argIdx + 1) < len(args) { // Ensure service-manifest has a filename parameter
				if processedFlags["--no-service-manifest"] {
					fmt.Printf("--service-manifest cannot be used in conjunction with --no-service-manifest\n")
					os.Exit(1)
				}

				manifestFilename = args[argIdx+1]

				if strings.HasPrefix(manifestFilename, "--") {
					fmt.Printf("--service-manifest requires a filename argument. \"%s\" was found instead\n", manifestFilename)
					os.Exit(1)
				}
				processedFlags[args[argIdx]] = true
				argIdx++ // Increment the index because we want to not process the 'filename' on the next loop
			} else {
				fmt.Printf("--service-manifest is missing a manifest filename argument\n")
				os.Exit(1)
			}
		case "--no-service-manifest": // Specify that no service manifests exists and so don't process any service creations, just push the app, if no-push was not specified.
			if processedFlags["--service-manifest"] {
				fmt.Printf("--no-service-manifest cannot be used in conjunction with --service-manifest\n")
				os.Exit(1)
			}

			manifestFilename = ""
			processedFlags[args[argIdx]] = true
		case "--no-push": // Specify that the services should be created but the app should not be pushed.
			processedFlags[args[argIdx]] = true
			continue

		default: // The parameter was not recognized, so we'll just push it into the remaining args list
			remainingArgs = append(remainingArgs, args[argIdx])
		}
	}

	// args should now only contain arguments that are not specific to create-service-push
	args = remainingArgs
	// End of arg processing

	// 2. Whatever the manifest file is, check to make sure it exists!
	if len(manifestFilename) > 0 {
		if _, err := os.Stat(manifestFilename); !os.IsNotExist(err) {
			fmt.Printf("Found Service Manifest File: %s\n", manifestFilename)
			filePointer, err := os.Open(manifestFilename)
			if err == nil {
				manifest, err := ParseManifest(filePointer)
				if err != nil {
					fmt.Printf("ERROR: %s\n", err)
					os.Exit(1)
				}

				createServicesobject := &CreateServicePush{
					manifest: &manifest,
					cf:       cliConnection,
				}
				if err := createServicesobject.createServices(); err != nil {
					os.Exit(1)
				}
			} else {
				fmt.Printf("ERROR: Unable to open %s.\n", manifestFilename)
				os.Exit(1)
			}
		} else {
			fmt.Printf("ERROR: The file %s was not found.\n", manifestFilename)
			os.Exit(1)
		}
	}

	// If no-push was specified, don't push the application. Otherwise, push the application
	// to CF.
	if _, ok := processedFlags["--no-push"]; ok {
		fmt.Printf("--no-push applied: Your application will not be pushed to CF ...\n")
	} else {
		fmt.Printf("Performing a CF Push with arguments [ %s ] ...\n", strings.Join(args, " "))

		newArgs := append([]string{"push"}, args...)
		// 3. Perform the cf push
		output, err := cliConnection.CliCommand(newArgs...)
		fmt.Printf("%s\n", output)

		if err != nil {
			fmt.Printf("ERROR while pushing: %s\n", err)
		}
	}
}

func (c *CreateServicePush) createServices() error {
	var err error
	// Detect the type of service and then go and create them.
	// credentials: User provided credentials service
	// drain: User provided log drain service
	// route: User provided route service
	// brokered: Brokered service.  The type field can be blank to specify this as well.
	for _, serviceObject := range c.manifest.Services {
		if serviceObject.Type == "credentials" {
			err = c.createUserProvidedCredentialsService(serviceObject.ServiceName, serviceObject.Credentials)
		} else if serviceObject.Type == "drain" {
			err = c.createUserProvidedLogDrainService(serviceObject.ServiceName, serviceObject.Url)
		} else if serviceObject.Type == "route" {
			err = c.createUserProvidedRouteService(serviceObject.ServiceName, serviceObject.Url)
		} else {
			if serviceObject.Type == "brokered" || serviceObject.Type == "" {
				err = c.createService(serviceObject.ServiceName,
					serviceObject.Broker,
					serviceObject.PlanName,
					serviceObject.JSONParameters,
					serviceObject.Tags,
					serviceObject.updateSIParamsAndTags)
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

func (c *CreateServicePush) run(args ...string) error {
	if os.Getenv("DEBUG") != "" {
		fmt.Printf(">> %s\n", strings.Join(args, " "))
	}

	fmt.Printf("Now Running CLI Command: %s\n", strings.Join(args, " "))
	_, err := c.cf.CliCommand(args...)
	return err
}

func (c *CreateServicePush) createUserProvidedCredentialsService(name string, credentials map[string]string) error {
	fmt.Printf("%s - ", name)
	s, err := c.cf.GetServices()
	if err != nil {
		return err
	}

	for _, svc := range s {
		if svc.Name == name {
			fmt.Print("already exists...skipping creation\n")
			return nil
		}
	}

	fmt.Print("will now be created as a user provided credential service.\n")
	credentialsJSON, _ := json.Marshal(credentials)

	return c.run("cups", name, "-p", string(credentialsJSON))
}

func (c *CreateServicePush) createUserProvidedRouteService(name, urlString string) error {
	fmt.Printf("%s - ", name)
	s, err := c.cf.GetServices()
	if err != nil {
		return err
	}

	for _, svc := range s {
		if svc.Name == name {
			fmt.Print("already exists...skipping creation\n")
			return nil
		}
	}

	// Check to ensure that the url begins with HTTPS because that is the only scheme supported.
	urlStruct, err := url.Parse(urlString)

	if err != nil {
		return err
	}

	if strings.ToLower(urlStruct.Scheme) != "https" {
		return fmt.Errorf("route scheme not specified or unsupported. User provided route service only supports https")
	}

	fmt.Print("will now be created as a user provided route service.\n")

	return c.run("cups", name, "-r", urlString)
}

func (c *CreateServicePush) createUserProvidedLogDrainService(name, urlString string) error {
	fmt.Printf("%s - ", name)
	s, err := c.cf.GetServices()
	if err != nil {
		return err
	}

	for _, svc := range s {
		if svc.Name == name {
			fmt.Print("already exists...skipping creation\n")
			return nil
		}
	}

	fmt.Print("will now be created as a user provided log drain service.\n")

	return c.run("cups", name, "-l", urlString)
}

func (c *CreateServicePush) createService(name, broker, plan, JSONParam, tags string, updateSIParamsAndTags bool) error {
	fmt.Printf("%s - ", name)
	s, err := c.cf.GetServices()
	if err != nil {
		return err
	}

	for _, svc := range s {
		if svc.Name == name {
			if !updateSIParamsAndTags {
				fmt.Print("already exists...skipping creation\n")
				return nil
			}
		}
	}

	fmt.Printf("will now be created as a brokered service.\n")

	// Process the parameters
	optionalArgs := []string{}
	if tags != "" {
		optionalArgs = append(optionalArgs, "-t")
		optionalArgs = append(optionalArgs, tags)
	}

	if JSONParam != "" {
		optionalArgs = append(optionalArgs, "-c")
		optionalArgs = append(optionalArgs, JSONParam)
	}

	if updateSIParamsAndTags {
		err = c.run(append([]string{"update-service", name}, optionalArgs...)...)
	} else {
		err = c.run(append([]string{"create-service", broker, plan, name}, optionalArgs...)...)
	}

	if err != nil {
		return err
	}

	pb := NewProgressSpinner(os.Stdout)
	for {
		service, err := c.cf.GetService(name)
		if err != nil {
			return err
		}

		pb.Next(service.LastOperation.Description)

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

// GetMetadata must be implemented as part of the plugin interface
// defined by the core CLI.
//
// GetMetadata() returns a PluginMetadata struct. The first field, Name,
// determines the name of the plugin which should generally be without spaces.
// If there are spaces in the name a user will need to properly quote the name
// during uninstall otherwise the name will be treated as seperate arguments.
// The second value is a slice of Command structs. Our slice only contains one
// Command Struct, but could contain any number of them. The first field Name
// defines the command `cf basic-plugin-command` once installed into the CLI. The
// second field, HelpText, is used by the core CLI to display help information
// to the user in the core commands `cf help`, `cf`, or `cf -h`.
func (c *CreateServicePush) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "Create-Service-Push",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 2,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "create-service-push",
				Alias:    "csp",
				HelpText: "Works in the same manner as cf push, except that it will create services defined in a services-manifest.yml file first before performing a cf push.",

				// UsageDetails is optional
				// It is used to show help of usage of each command
				UsageDetails: plugin.Usage{
					Usage: "create-service-push\n   cf create-service-push",
					Options: map[string]string{
						"--service-manifest <MANIFEST_FILE>": "Specify the fullpath and filename of the services creation manifest.  Defaults to services-manifest.yml.",
						"--no-service-manifest":              "Specifies that there is no service creation manifest",
						"--no-push":                          "Create the services but do not push the application",
					},
				},
			},
		},
	}
}

// Unlike most Go programs, the `Main()` function will not be used to run all of the
// commands provided in your plugin. Main will be used to initialize the plugin
// process, as well as any dependencies you might require for your
// plugin.
func main() {
	// Any initialization for your plugin can be handled here
	//
	// Note: to run the plugin.Start method, we pass in a pointer to the struct
	// implementing the interface defined at "code.cloudfoundry.org/cli/plugin/plugin.go"
	//
	// Note: The plugin's main() method is invoked at install time to collect
	// metadata. The plugin will exit 0 and the Run([]string) method will not be
	// invoked.
	plugin.Start(new(CreateServicePush))
	// Plugin code should be written in the Run([]string) method,
	// ensuring the plugin environment is bootstrapped.
}
