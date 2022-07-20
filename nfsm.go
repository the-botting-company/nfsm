package nfsm

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
)

var (
	ErrStateMachineAlreadyRunning = errors.New("state machine is already running")
)

// MachineCtx provides an interface for which a state can access information on its machine.
type Machine interface {
	// Context returns the state machines context.
	Context() context.Context
	// Previous returns the state machines previous state.
	Previous() string
	// Current returns the state machines current state.
	Current() string
}

// Handler performs some action during it's state defined in Handlers then returns the next state or error. Error ends the execution of the machine.
type Handler func(nfsm Machine) (string, error)

// Handlers maps a state to its Handler.
type Handlers map[string]Handler

// Nfsm represents a type of non-deterministic state machine.
type Nfsm struct {
	initial string
	final   string

	handlers Handlers

	previous string
	current  string

	running int32

	ctx    context.Context
	cancel context.CancelFunc
}

// NewNfsm creates a new instance of Nfsm.
func NewNfsm(ctx context.Context, initial, final string, handlers Handlers) *Nfsm {
	ctx, c := context.WithCancel(ctx)
	return &Nfsm{
		initial:  initial,
		final:    final,
		handlers: handlers,
		ctx:      ctx,
		cancel:   c,
	}
}

// Execute starts the state machine. It must not be running.
func (n *Nfsm) Execute() error {
	if !atomic.CompareAndSwapInt32(&n.running, 0, 1) {
		return ErrStateMachineAlreadyRunning
	}

	defer atomic.SwapInt32(&n.running, 0)

	if n.handlers[n.initial] == nil {
		return fmt.Errorf("state %s does not exist", n.initial)
	}

	next, err := n.callHandler(n.initial)
	if err != nil {
		return err
	}

	for {
		select {
		case <-n.Context().Done():
			return nil
		default:
			if next == "" {
				return nil
			}

			if n.handlers[next] == nil {
				return fmt.Errorf("state %s does not exist", next)
			}

			next, err = n.callHandler(next)
			if err != nil {
				return err
			}
		}
	}
}

func (n *Nfsm) callHandler(state string) (string, error) {
	n.setCurrent(state)

	s, err := n.handlers[state](n)
	if err != nil {
		return "", err
	}

	n.setPrevious(state)

	return s, nil
}

// Final returns the name of the state machines final state.
func (n *Nfsm) Final() string {
	return n.final
}

// Previous returns the state machines previous state.
func (n *Nfsm) Previous() string {
	return n.previous
}

func (n *Nfsm) setPrevious(p string) {
	n.previous = p
}

// Current returns the state machines current state.
func (n *Nfsm) Current() string {
	return n.current
}

func (n *Nfsm) setCurrent(c string) {
	n.current = c
}

// Running returns whether the machine is currently running or not.
func (n *Nfsm) Running() bool {
	return atomic.LoadInt32(&n.running) == 1
}

// Context will return the state machines context.
func (n *Nfsm) Context() context.Context {
	return n.ctx
}

// Cancel will call the state machines CancelFunc.
func (n *Nfsm) Cancel() {
	n.cancel()
}

// Factory will create a new copy of Nfsm.
func (n *Nfsm) Factory(ctx context.Context) *Nfsm {
	return NewNfsm(ctx, n.initial, n.final, n.handlers)
}
