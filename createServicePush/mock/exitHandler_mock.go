package createService_mock

// MockExitHandler is the struct holding the exit hander
type MockExitHandler struct {
	Exit1WasCalled bool
	Exit0WasCalled bool
}

// NewMockExitHandler creates a NewMockExitHandler struct
func NewMockExitHandler() *MockExitHandler {
	return &MockExitHandler{}
}

// HandleError is the method to deal with errors on exit.
func (eh *MockExitHandler) HandleError() {
	eh.Exit1WasCalled = true
}

// HandleOK is the method to deal with exiting the plugin while it was ok to
func (eh *MockExitHandler) HandleOK() {
	eh.Exit0WasCalled = true
}
