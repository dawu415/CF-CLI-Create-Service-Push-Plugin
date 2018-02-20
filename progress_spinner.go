package main

import (
	"fmt"
	"io"
)

type ProgressSpinner struct {
	state                string
	loadingMessage       string
	writer               io.Writer
}

func NewProgressSpinner(writer io.Writer) *ProgressSpinner {
	ipb := &ProgressSpinner{"|", "", writer}
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
		fmt.Fprint(ps.writer, ps.loadingMessage+"\n")
	}
	fmt.Fprint(ps.writer, ps.state+"\r")
}
