package cspArguments

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry/bosh-cli/director/template"
)

// Interface describes the interface to process input commandline arguments
type Interface interface {
	Process(args []string) (*CSPArguments, error)
}

// CSPArguments holds the Processed input arguments
type CSPArguments struct {
	IsUninstallingPlugin     bool
	ServiceManifestFilename  string
	DoNotCreateServices      bool
	DoNotPush                bool
	StaticVariablesFilePaths []string
	StaticVariables          []template.VarKV
	OtherCFArgs              []string // Holds other commandline arguments that isn't used by CSP. This will be passed to cf push.
}

// NewCSPArguments returns an initialized CSPArguments struct
func NewCSPArguments() *CSPArguments {
	return &CSPArguments{
		ServiceManifestFilename:  "services-manifest.yml",
		DoNotCreateServices:      false,
		DoNotPush:                false,
		StaticVariablesFilePaths: []string{},
		StaticVariables:          []template.VarKV{},
		OtherCFArgs:              []string{},
	}
}

// Process function is the entrypoint for processing Create-Service-Push Arguments
func (csp *CSPArguments) Process(args []string) (*CSPArguments, error) {

	if args[0] == "CLI-MESSAGE-UNINSTALL" {
		csp.IsUninstallingPlugin = true
		return csp, nil
	}

	if args[0] != "create-service-push" {
		return csp, fmt.Errorf("This plugin only works with create-service-push")
	}
	args = args[1:] // Remove the create-service-push tag

	// Process the input arguments specific to create-service-push
	processedFlags := map[string]bool{}

	for argIdx := 0; argIdx < len(args); argIdx++ {
		switch args[argIdx] {
		case "--vars-file": // A file containing static variables for service manifest interpolation
			if (argIdx + 1) < len(args) { // Ensure vars-file has a filename parameter
				if strings.HasPrefix(args[argIdx+1], "--") {
					return csp, fmt.Errorf(
						"--vars-file requires a filename argument. \"%s\" was found instead", args[argIdx+1])
				}
				csp.StaticVariablesFilePaths = append(csp.StaticVariablesFilePaths, args[argIdx+1])

				// Pass the --vars-file to the cf push command as well.
				csp.OtherCFArgs = append(csp.OtherCFArgs, args[argIdx])
				csp.OtherCFArgs = append(csp.OtherCFArgs, args[argIdx+1])
				argIdx++
			} else {
				return csp, fmt.Errorf("--vars-file is missing a variable yaml filename argument")
			}
		case "--var": // A key=value static variable for service manifest interpolation
			if (argIdx + 1) < len(args) { // Ensure var has a key=value parameter
				// We expect the next argument to be on of the form 'a=b'
				if strings.Contains(args[argIdx+1], "=") {
					tokenString := strings.SplitN(args[argIdx+1], "=", 2)

					if tokenString[1] == " " || len(tokenString[1]) == 0 {
						// Better check to make sure there's no hanging argument next
						if (argIdx + 2) < len(args) {
							if !strings.HasPrefix(args[argIdx+2], "-") {
								// We've got a hanging input here
								return csp, fmt.Errorf("%s seems to be a hanging input. Ensure there are no spaces between the equals sign", args[argIdx+2])
							}
						}
					}

					csp.StaticVariables = append(csp.StaticVariables,
						template.VarKV{
							Name:  tokenString[0],
							Value: tokenString[1],
						})

					// Pass the --var to the cf push command as well.
					csp.OtherCFArgs = append(csp.OtherCFArgs, args[argIdx])
					csp.OtherCFArgs = append(csp.OtherCFArgs, args[argIdx+1])
					argIdx++
				} else {
					return csp, fmt.Errorf("%s does not seem to be of the form key=value. Ensure there are no spaces between the equals sign", args[argIdx+1])
				}
			} else {
				return csp, fmt.Errorf("--var is missing a key=value pair argument")
			}
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
