package createServicePush

import (
	"fmt"
	"os"
	"strings"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/dawu415/CF-CLI-Create-Service-Push-Plugin/cspArguments"
	"github.com/dawu415/CF-CLI-Create-Service-Push-Plugin/serviceCreator"
	"github.com/dawu415/CF-CLI-Create-Service-Push-Plugin/serviceManifest"
)

// CreateServicePush is the struct implementing the interface defined by the core CLI. It can
// be found at  "code.cloudfoundry.org/cli/plugin/plugin.go"
type CreateServicePush struct {
}

// Create instantiates a new CreateServicePush struct and returns it as a pointer
func Create() *CreateServicePush {
	return &CreateServicePush{}
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

	// Process the input arguments
	CSPArguments, err := cspArguments.Process(args)

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}

	// If we are specified to process a service manifest (by default), then
	// read in the service manifest and instantiate the services from that
	if !CSPArguments.DoNotCreateServices {
		parser, err := serviceManifest.NewParser(CSPArguments.ServiceManifestFilename)

		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}

		manifest, err := parser.Parse()

		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}

		serviceCreator.CreateServices(manifest, cliConnection)
	}

	// If no-push was specified, don't push the application. Otherwise, push the application
	// to CF.
	if CSPArguments.DoNotPush {
		fmt.Printf("--no-push applied: Your application will not be pushed to CF ...\n")
	} else {
		fmt.Printf("Performing a CF Push with arguments [ %s ] ...\n", strings.Join(args, " "))

		// Perform the cf push
		output, err := cliConnection.CliCommand(append([]string{"push"}, args...)...)
		fmt.Printf("%s\n", output)

		if err != nil {
			fmt.Printf("ERROR while pushing: %s\n", err)
		}
	}
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
