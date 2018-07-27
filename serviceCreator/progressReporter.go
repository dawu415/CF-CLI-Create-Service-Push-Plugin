package serviceCreator

import (
	"fmt"
)

// LogFunc is logger function signature. We'll use this to set the function to output the progress output
type LogFunc func(format string, a ...interface{}) (n int, err error)

// ProgressReporter describes the state and current message displayed
type ProgressReporter struct {
	state          string
	loadingMessage string
	log            LogFunc
}

// NewProgressReporter initializes and creates a New Progress Reporter defaulting output to stdout
func NewProgressReporter() *ProgressReporter {
	return NewProgressReporterWithLoggerOut(fmt.Printf)
}

// NewProgressReporterWithLoggerOut initializes and creates a New Progress Reporter with a specific log output function
func NewProgressReporterWithLoggerOut(outputFunction LogFunc) *ProgressReporter {
	return &ProgressReporter{"|", "", outputFunction}
}

// Step updates the Progress Reporter State and Display message on screen
func (ps *ProgressReporter) Step(loadingMessage string) {
	var nextState string

	if loadingMessage != ps.loadingMessage {
		ps.loadingMessage = loadingMessage
		ps.log(ps.loadingMessage + "\n")
	}
	ps.log(ps.state + "\r")

	switch ps.state {
	case "|":
		nextState = "/"
		break
	case "/":
		nextState = "-"
		break
	case "-":
		nextState = "\\"
		break
	default:
		nextState = "|"
	}
	ps.state = nextState

}
