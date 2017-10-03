package main

import (
	"fmt"
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

	fmt.Printf("Arguments: %s\n", strings.Join(args, "**"))

	// 1. Find an argument of -f in the list.  This will tell us the manifest file
	var manifestFilename = "manifest.yml"

	for i, arg := range args {
		if arg == "-f" {
			manifestFilename = args[i+1]
			break
		} else if arg == "--no-manifest" {
			manifestFilename = ""
			break
		}
	}

	fmt.Printf("Found ManifestFile: %s\n", manifestFilename)

	// 2. Whatever the manifest file is, check to make sure it exists!
	if len(manifestFilename) > 0 {
		if _, err := os.Stat(manifestFilename); !os.IsNotExist(err) {
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

				createServicesobject.createServices()
			} else {
				fmt.Printf("ERROR: %s\n", err)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("Performing a CF Push with arguments %s\n", strings.Join(args[1:], " "))

	newArgs := append([]string{"push"}, args[1:]...)
	// 3. Perform the cf push
	output, err := cliConnection.CliCommand(newArgs...)
	fmt.Printf("%s\n", output)

	if err != nil {
		fmt.Printf("ERROR while pushing: %s\n", err)
	}
}

func (c *CreateServicePush) createServices() error {

	for _, serviceObject := range c.manifest.Services {
		if err := c.createService(serviceObject.ServiceName, serviceObject.Broker, serviceObject.PlanName, serviceObject.JSONParameters); err != nil {
			fmt.Printf("Create Service Error: %+v ", err)
		}
	}

	return nil
}

func (c *CreateServicePush) run(args ...string) error {
	if os.Getenv("DEBUG") != "" {
		fmt.Printf(">> %s\n", strings.Join(args, " "))
	}
	if os.Getenv("DRYRUN") != "" {
		return nil
	}

	fmt.Printf("Now Running CLI Command: %s\n", strings.Join(args, " "))
	_, err := c.cf.CliCommand(args...)
	return err
}

func (c *CreateServicePush) createService(name, broker, plan, JSONParam string) error {
	s, err := c.cf.GetServices()
	if err != nil {
		return err
	}

	fmt.Printf("Checking existence of services\n")

	for _, svc := range s {
		if svc.Name == name {

			fmt.Printf("%s already exists. Stopping service creation\n", name)
			/* FIXME: check configuration */
			return nil
		}
	}

	fmt.Printf("Creating Service: %s\n", name)

	var result error
	if JSONParam == "" {
		result = c.run("create-service", broker, plan, name)
	} else {
		result = c.run("create-service", broker, plan, name, "-c", JSONParam)
	}

	return result
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
		Name: "MyCreateServicePush",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
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
				HelpText: "Works in the same manner as cf push, except that it will create services defined in a service manifest",

				// UsageDetails is optional
				// It is used to show help of usage of each command
				UsageDetails: plugin.Usage{
					Usage: "create-service-push\n   cf create-service-push",
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
