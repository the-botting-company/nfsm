package nfsm

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var (
	ErrStateMachineRunning = errors.New("state machine is running")
)

// MachineCtx provides an interface for which a state can access information on its machine.
type Machine interface {
	// Metadata provides a way to pass data across states.
	Metadata() *Metadata
	// Previous returns the state machines previous state.
	Previous() string
	// Current returns the state machines current state.
	Current() string
}

// Nfsm represents a type of non-deterministic state machine.
type Nfsm struct {
	flow *Flow

	metadata *Metadata

	previous string
	current  string

	running int32

	statesMu sync.Mutex
}

// NewNfsm creates a new instance of Nfsm.
func NewNfsm(flow *Flow) *Nfsm {
	return &Nfsm{
		flow:     flow,
		metadata: NewMetadata(),
	}
}

// Execute starts the state machine. It must not be running.
func (n *Nfsm) Execute(ctx context.Context) error {
	if !atomic.CompareAndSwapInt32(&n.running, 0, 1) {
		return ErrStateMachineRunning
	}

	defer atomic.SwapInt32(&n.running, 0)

	if n.flow.handlers[n.flow.initial] == nil {
		return fmt.Errorf("state %s does not exist", n.flow.initial)
	}

	next, err := n.callHandler(ctx, n.flow.initial)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if next == "" {
				return nil
			}

			if n.flow.handlers[next] == nil {
				return fmt.Errorf("state %s does not exist", next)
			}

			next, err = n.callHandler(ctx, next)
			if err != nil {
				return err
			}
		}
	}
}

func (n *Nfsm) callHandler(ctx context.Context, state string) (string, error) {
	n.statesMu.Lock()
	defer n.statesMu.Unlock()

	n.setCurrent(state)

	s, err := n.flow.handlers[state](ctx, n)
	if err != nil {
		return "", err
	}

	n.setPrevious(state)

	return s, nil
}

func (n *Nfsm) Metadata() *Metadata {
	return n.metadata
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
