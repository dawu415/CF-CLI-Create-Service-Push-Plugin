package cspArguments

import (
	"fmt"
	"os"
	"strings"
)

// Interface describes the interface to process input commandline arguments
type Interface interface {
	Process(args []string) (*CSPArguments, error)
	GetUsage() string
	GetArgumentsDescription() map[string]string
}

// CSPFlagProperty defines a flag in create-service-push, its properties and a handler to decode output to a CSPArgument struct
type CSPFlagProperty struct {
	description   string
	argumentCount int
	handler       func(int, []string, *CSPArguments, *error) // func(index, argument list, outputArgument, error)
	processed     bool
	shouldDefer   bool // Specifies if this flag processing should be defered to the end of our Process method
}

// CSPArguments holds the Processed input arguments
type CSPArguments struct {
	IsUninstallingPlugin     bool
	ServiceManifestFilename  string
	DoNotCreateServices      bool
	DoNotPush                bool
	PushAsSubProcess         bool
	StaticVariablesFilePaths []string
	StaticVariables          map[string]string
	OtherCFArgs              []string                    // Holds other commandline arguments that isn't used by CSP. This will be passed to cf push.
	cspFlags                 map[string]*CSPFlagProperty // Private variable
}

// NewCSPArguments returns an initialized CSPArguments struct
func NewCSPArguments() *CSPArguments {
	return &CSPArguments{
		ServiceManifestFilename:  "services-manifest.yml",
		DoNotCreateServices:      false,
		DoNotPush:                false,
		PushAsSubProcess:         false,
		StaticVariablesFilePaths: []string{},
		StaticVariables:          map[string]string{},
		OtherCFArgs:              []string{},

		cspFlags: map[string]*CSPFlagProperty{
			/////////////////////////////////////////////////
			"--var": &CSPFlagProperty{
				description:   "Takes one input being a variable key value pair for variable substitution, (e.g., name=app1); can specify multiple times",
				argumentCount: 1,
				handler: func(index int, args []string, csp *CSPArguments, err *error) {
					if (index + 1) < len(args) { // Ensure var has a key=value parameter
						if strings.Contains(args[index+1], "=") {
							tokenString := strings.SplitN(args[index+1], "=", 2)

							if tokenString[1] == " " || len(tokenString[1]) == 0 {
								// Better check to make sure there's no hanging argument next, e.g., hanging value: ["--var","key=", "value"]
								if (index + 2) < len(args) {
									if !strings.HasPrefix(args[index+2], "-") {
										// We've got a hanging input here
										*err = fmt.Errorf("%s seems to be a hanging input. Ensure there are no spaces between the equals sign", args[index+2])
										return
									}
								}
							}

							csp.StaticVariables[tokenString[0]] = tokenString[1]

							// We expect the next argument to be on of the form 'a=b', and arg array should be like ["--var", "key=value"]
							// If the push as sub-process was actually processed, we can these flags as well
							if csp.cspFlags["--push-as-subprocess"].processed {
								// Pass the --var to the cf push command as well.
								csp.OtherCFArgs = append(csp.OtherCFArgs, args[index])
								csp.OtherCFArgs = append(csp.OtherCFArgs, args[index+1])
							}

							csp.cspFlags["--var"].processed = true
						} else {
							*err = fmt.Errorf("%s does not seem to be of the form key=value. Ensure there are no spaces between the equals sign", args[index+1])
							return
						}
					} else {
						*err = fmt.Errorf("--var is missing a key=value pair argument")
						return
					}
					*err = nil
				},
				processed:   false,
				shouldDefer: true, // We need to defer because we want to ensure push-as-subprocess is processed first
			},
			/////////////////////////////////////////////////
			"--vars-file": &CSPFlagProperty{
				description:   "Takes one input being the path to a variables file; can specify multiple times",
				argumentCount: 1,
				handler: func(index int, args []string, csp *CSPArguments, err *error) {
					if (index + 1) < len(args) { // Ensure vars-file has a filename parameter
						if strings.HasPrefix(args[index+1], "-") {
							*err = fmt.Errorf(
								"--vars-file requires a filename argument. \"%s\" was found instead", args[index+1])
							return
						}

						// If the push as sub-process was actually processed, we can these flags as well
						if csp.cspFlags["--push-as-subprocess"].processed {
							// Pass the --vars-file to the cf push command as well.
							csp.OtherCFArgs = append(csp.OtherCFArgs, args[index])
							csp.OtherCFArgs = append(csp.OtherCFArgs, args[index+1])
						}

						csp.StaticVariablesFilePaths = append(csp.StaticVariablesFilePaths, args[index+1])
						csp.cspFlags["--vars-file"].processed = true

					} else {
						*err = fmt.Errorf("--vars-file is missing a variable yaml filename argument")
						return
					}
					*err = nil
				},
				processed:   false,
				shouldDefer: true, // We need to defer because we want to ensure push-as-subprocess is processed first
			},
			/////////////////////////////////////////////////
			"--use-env-vars-prefixed-with": &CSPFlagProperty{
				description:   "Use environment variables that have a given prefix as substitution variables, i.e. --use-env-vars-prefixed-with APP_ will get all environment variables prefixed with APP_",
				argumentCount: 1,
				handler: func(index int, args []string, csp *CSPArguments, err *error) {
					if (index + 1) < len(args) { // Ensure a prefix has been specified
						// Get all environments prefixed with an input specified by args[index+1]

						for _, env := range os.Environ() {
							if strings.HasPrefix(env, args[index+1]) {
								tokenString := strings.SplitN(env, "=", 2)
								csp.StaticVariables[tokenString[0]] = tokenString[1]
							}
						}

						csp.cspFlags["--use-env-vars-prefixed-with"].processed = true
					} else {
						*err = fmt.Errorf("--use-env-vars-prefixed-with is missing a prefix input")
						return
					}
					*err = nil
				},
				processed:   false,
				shouldDefer: false,
			},
			/////////////////////////////////////////////////
			"--no-push": &CSPFlagProperty{
				description:   "Create the services but do not push the application",
				argumentCount: 0,
				handler: func(index int, args []string, csp *CSPArguments, err *error) {
					if csp.cspFlags["--push-as-subprocess"].processed {
						*err = fmt.Errorf("--no-push cannot be used in conjunction with --push-as-subprocess")
						return
					}
					*err = nil
					csp.DoNotPush = true
					csp.cspFlags["--no-push"].processed = true
				},
				processed:   false,
				shouldDefer: false,
			},
			/////////////////////////////////////////////////
			"--push-as-subprocess": &CSPFlagProperty{
				description:   "Perform cf push as a sub-process",
				argumentCount: 0,
				handler: func(index int, args []string, csp *CSPArguments, err *error) {

					if csp.cspFlags["--no-push"].processed {
						*err = fmt.Errorf("--push-as-subprocess cannot be used in conjunction with --no-push")
						return
					}
					*err = nil
					csp.PushAsSubProcess = true
					csp.cspFlags["--push-as-subprocess"].processed = true
				},
				processed:   false,
				shouldDefer: false,
			},
			/////////////////////////////////////////////////
			"--no-service-manifest": &CSPFlagProperty{
				description:   "Specifies that there is no service creation manifest",
				argumentCount: 0,
				handler: func(index int, args []string, csp *CSPArguments, err *error) {
					if csp.cspFlags["--service-manifest"].processed {
						*err = fmt.Errorf("--no-service-manifest cannot be used in conjunction with --service-manifest")
						return
					}
					*err = nil
					csp.DoNotCreateServices = true
					csp.cspFlags["--no-service-manifest"].processed = true
				},
				processed:   false,
				shouldDefer: false,
			},
			/////////////////////////////////////////////////
			"--service-manifest": &CSPFlagProperty{
				description:   "Takes one input specifying the fullpath and filename of the services creation manifest. e.g., --service-manifest my-manifest.yml. Defaults to services-manifest.yml.",
				argumentCount: 1,
				handler: func(index int, args []string, csp *CSPArguments, err *error) {
					if (index + 1) < len(args) { // Ensure service-manifest has a filename parameter
						if csp.cspFlags["--no-service-manifest"].processed {
							*err = fmt.Errorf("--service-manifest cannot be used in conjunction with --no-service-manifest")
							return
						}

						if strings.HasPrefix(args[index+1], "-") {
							*err = fmt.Errorf(
								"--service-manifest requires a filename argument. \"%s\" was found instead",
								csp.ServiceManifestFilename)
							return
						}

						csp.ServiceManifestFilename = args[index+1]
						csp.cspFlags["--service-manifest"].processed = true
					} else {
						*err = fmt.Errorf("--service-manifest is missing a manifest filename argument")
						return
					}
					*err = nil
				},
				processed:   false,
				shouldDefer: false,
			},
		},
	}
}

// GetUsage returns the usage instruction text to display when help is called.
func (csp *CSPArguments) GetUsage() string {
	return `
    cf create-service-push [APP_NAME] 
                           [ --service-manifest SERVICE_MANIFEST_FULL_PATH | --no-service-manifest ]
                           [ --no-push | --push-as-subprocess ]
                           [ --var KEY=VALUE ] [ --vars-file VARS_FILE_FULL_PATH ]
                           [ --use-env-vars-prefixed-with PREFIX ]
                           [CF_PUSH_ARGUMENTS]
    NOTES:
    a) APP_NAME is optional but should always be at the first position. cf push will validate this.
       If specified. It is simply passed to cf push in that order. If APP_NAME is not specified,
       cf push will expect to find it in an application manifest. 

    b) --push-as-subprocess overrides the cf cli plugin call to cf push and calls the cf cli installed on the machine instead.
       This flag will search for the cf cli in the current working directory first. If it is not present, it will search the 
       paths in the PATH environment variable. This was introduced as a workaround where the cf cli plugin architecture did not
       incorporate new features, such as --var and --vars-file. See https://github.com/cloudfoundry/cli/issues/1399#issuecomment-409061226 .

    c) By default --var and --vars-file will not be passed to cf push due to non-support in the cf plugin architecture. However,
       support for variable substitution is built into the plugin.  To pass --var and --vars-file parameters to cf push, use 
       it with the --push-as-subprocess flag. Please ensure that the cf cli installed on the machine is at least release 6.37.0.
       `
}

// GetArgumentsDescription returns the usage instruction text to display when help is called.
func (csp *CSPArguments) GetArgumentsDescription() map[string]string {
	arguments := make(map[string]string)
	for flag, property := range csp.cspFlags {
		arguments[flag] = property.description
	}
	return arguments
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

	// If there were no other arguments, we can short circuit the Process method
	if len(args) == 0 {
		return csp, nil
	}

	// Iterate through the possible set of arguments to see if they're in the args list
	errArray := make([]error, len(args))

	func() { // Call in a func so that our deferred handlers get called and have an output we check for errors
		for idx := 0; idx < len(args); idx++ {
			arg := args[idx]
			if property, isCSPFlag := csp.cspFlags[arg]; isCSPFlag {
				if property.shouldDefer {
					// These properties should be handled at the end because
					// it is dependent on another property.
					defer property.handler(idx, args, csp, &errArray[idx])
				} else {
					property.handler(idx, args, csp, &errArray[idx])
					if errArray[idx] != nil {
						return
					}
				}
				idx += property.argumentCount
			} else {
				csp.OtherCFArgs = append(csp.OtherCFArgs, arg)
			}
		}
	}() // <-- Note that we're calling the function

	// Just return the first error encountered
	for _, err := range errArray {
		if err != nil {
			return csp, err
		}
	}

	return csp, nil
}
