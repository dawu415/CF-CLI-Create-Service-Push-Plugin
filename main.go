package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"github.com/dawu415/CF-CLI-Create-Service-Push-Plugin/createServicePush"
)

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
	plugin.Start(createServicePush.Create())
	// Plugin code should be written in the Run method,
	// ensuring the plugin environment is bootstrapped.
}
