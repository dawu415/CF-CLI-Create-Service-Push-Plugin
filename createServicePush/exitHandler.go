package createServicePush

import (
	"os"
)

// ExitInterface is an interface to an Exit Handler
type ExitInterface interface {
	HandleError()
	HandleOK()
}

// ExitHandler is the struct holding the exit hander
type ExitHandler struct {
	Exit ExitInterface
}

// NewExitHandler creates the default ExitHandler
func NewExitHandler() *ExitHandler {
	return &ExitHandler{&ExitHandler{}}
}

// HandleError is the method to deal with errors on exit.
func (eh *ExitHandler) HandleError() {
	os.Exit(1)
}

// HandleOK is the method to deal with exiting the plugin while it was ok to
func (eh *ExitHandler) HandleOK() {
	os.Exit(0)
}
