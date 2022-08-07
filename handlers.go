package nfsm

// Handler performs some action during it's state defined in Handlers then returns the next state or error. An empty string or error ends the execution of the machine.
type Handler func(nfsm Machine) (string, error)

// Handlers maps a state to its Handler.
type Handlers map[string]Handler

// Flow represents the flow of the state machine.
type Flow struct {
	// initial is the initial state of the machine.
	initial string
	handlers Handlers
}

// NewFlow creates a new instance of Flow.
func NewFlow(initial string, handlers Handlers) Flow {
	return Flow{
		initial,
		handlers,
	}
}

// Initial returns the initial property of Flow.
func (f Flow) Initial() string {
	return f.initial
}

// Handlers returns the handlers property on Flow.
func (f Flow) Handlers() Handlers {
	return f.handlers
}