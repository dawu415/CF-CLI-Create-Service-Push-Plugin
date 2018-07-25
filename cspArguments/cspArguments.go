package cspArguments

import (
	"fmt"
	"strings"
)

// Interface describes the interface to process input commandline arguments
type Interface interface {
	Process(args []string) (*CSPArguments, error)
}

// CSPArguments holds the Processed input arguments
type CSPArguments struct {
	ServiceManifestFilename string
	DoNotCreateServices     bool
	DoNotPush               bool
	OtherCFArgs             []string // Holds other commandline arguments that isn't used by CSP. This will be passed to cf push.
}

// NewCSPArguments returns an initialized CSPArguments struct
func NewCSPArguments() *CSPArguments {
	return &CSPArguments{
		ServiceManifestFilename: "service-manifest.yml",
		DoNotCreateServices:     false,
		DoNotPush:               false,
		OtherCFArgs:             []string{},
	}
}

// Process function is the entrypoint for processing Create-Service-Push Arguments
func (csp *CSPArguments) Process(args []string) (*CSPArguments, error) {
	if args[0] != "create-service-push" {
		return csp, fmt.Errorf("This plugin only works with create-service-push")
	}

	args = args[1:] // Remove the create-service-push tag

	// Process the input arguments specific to create-service-push
	processedFlags := map[string]bool{}

	for argIdx := 0; argIdx < len(args); argIdx++ {
		switch args[argIdx] {
		case "--service-manifest": // Specify the service manifest to use
			if (argIdx + 1) < len(args) { // Ensure service-manifest has a filename parameter
				if processedFlags["--no-service-manifest"] {
					return csp, fmt.Errorf("--service-manifest cannot be used in conjunction with --no-service-manifest")
				}

				csp.ServiceManifestFilename = args[argIdx+1]

				if strings.HasPrefix(csp.ServiceManifestFilename, "--") {
					return csp, fmt.Errorf(
						"--service-manifest requires a filename argument. \"%s\" was found instead",
						csp.ServiceManifestFilename)
				}
				processedFlags[args[argIdx]] = true
				argIdx++ // Increment the index because we want to not process the 'filename' on the next loop
			} else {
				return csp, fmt.Errorf("--service-manifest is missing a manifest filename argument")
			}
		case "--no-service-manifest": // Specify that no service manifests exists and so don't process any service creations, just push the app, if no-push was not specified.
			if processedFlags["--service-manifest"] {
				return csp, fmt.Errorf("--no-service-manifest cannot be used in conjunction with --service-manifest")
			}

			csp.DoNotCreateServices = true

			processedFlags[args[argIdx]] = true
		case "--no-push": // Specify that the services should be created but the app should not be pushed.
			processedFlags[args[argIdx]] = true
			csp.DoNotPush = true
			continue

		default: // The parameter was not recognized, so we'll just push it into the remaining args list
			csp.OtherCFArgs = append(csp.OtherCFArgs, args[argIdx])
		}
	}
	return csp, nil
}
