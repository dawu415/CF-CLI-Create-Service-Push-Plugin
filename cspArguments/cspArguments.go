package cspArguments

import (
	"fmt"
	"strings"
)

// CSPArguments holds the Processed input arguments
type CSPArguments struct {
	ServiceManifestFilename string
	DoNotCreateServices     bool
	DoNotPush               bool
	OtherCFArgs             []string // Holds other commandline arguments that isn't used by CSP. This will be passed to cf push.
}

// Process function is the entrypoint for processing Create-Service-Push Arguments
func Process(args []string) (CSPArguments, error) {
	if args[0] != "create-service-push" {
		return CSPArguments{}, fmt.Errorf("This plugin only works with create-service-push")
	}

	args = args[1:] // Remove the create-service-push tag

	// Process the input arguments specific to create-service-push
	serviceManifestFilename := "service-manifest.yml" // Default service manifest name
	doNotCreateServices := false
	doNotPush := false
	otherCFArgs := []string{}

	processedFlags := map[string]bool{}

	for argIdx := 0; argIdx < len(args); argIdx++ {
		switch args[argIdx] {
		case "--service-manifest": // Specify the service manifest to use
			if (argIdx + 1) < len(args) { // Ensure service-manifest has a filename parameter
				if processedFlags["--no-service-manifest"] {
					return CSPArguments{}, fmt.Errorf("--service-manifest cannot be used in conjunction with --no-service-manifest")
				}

				serviceManifestFilename = args[argIdx+1]

				if strings.HasPrefix(serviceManifestFilename, "--") {
					return CSPArguments{}, fmt.Errorf("--service-manifest requires a filename argument. \"%s\" was found instead", serviceManifestFilename)
				}
				processedFlags[args[argIdx]] = true
				argIdx++ // Increment the index because we want to not process the 'filename' on the next loop
			} else {
				return CSPArguments{}, fmt.Errorf("--service-manifest is missing a manifest filename argument")
			}
		case "--no-service-manifest": // Specify that no service manifests exists and so don't process any service creations, just push the app, if no-push was not specified.
			if processedFlags["--service-manifest"] {
				return CSPArguments{}, fmt.Errorf("--no-service-manifest cannot be used in conjunction with --service-manifest")
			}

			serviceManifestFilename = ""
			doNotCreateServices = true

			processedFlags[args[argIdx]] = true
		case "--no-push": // Specify that the services should be created but the app should not be pushed.
			processedFlags[args[argIdx]] = true
			doNotPush = true
			continue

		default: // The parameter was not recognized, so we'll just push it into the remaining args list
			otherCFArgs = append(otherCFArgs, args[argIdx])
		}
	}
	return CSPArguments{
		serviceManifestFilename,
		doNotCreateServices,
		doNotPush,
		otherCFArgs,
	}, nil
}
