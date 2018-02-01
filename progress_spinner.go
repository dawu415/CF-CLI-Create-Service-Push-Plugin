package main

import (
	"fmt"
	"io"
)

type ProgressSpinner struct {
	state                string
	loadingMessage       string
	updateLoadingMessage bool
	writer               io.Writer
}

func NewProgressSpinner(writer io.Writer) *ProgressSpinner {
	ipb := &ProgressSpinner{"|", "", true,writer}
	return ipb
}
func (ps *ProgressSpinner) Next(loadingMessage string) {
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
		ps.updateLoadingMessage = true
	} else {
		ps.updateLoadingMessage = false
	}
	ps.write()
}
func (ps *ProgressSpinner) write() {
	if ps.updateLoadingMessage {
		fmt.Fprint(ps.writer, ps.loadingMessage+"\n")
	}
	fmt.Fprint(ps.writer, ps.state+"\r")
}