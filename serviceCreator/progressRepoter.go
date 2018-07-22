package serviceCreator

import (
	"fmt"
)

// ProgressReporter describes the state and current message displayed
type ProgressReporter struct {
	state          string
	loadingMessage string
}

// NewProgressReporter initializes and creates a New Progress Reporter
func NewProgressReporter() *ProgressReporter {
	ipb := &ProgressReporter{"|", ""}
	return ipb
}

// Step updates the Progress Reporter State and Display message on screen
func (ps *ProgressReporter) Step(loadingMessage string) {
	var nextState string
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
	if loadingMessage != ps.loadingMessage {
		ps.loadingMessage = loadingMessage
		fmt.Printf(ps.loadingMessage + "\n")
	}
	fmt.Printf(ps.state + "\r")
}
