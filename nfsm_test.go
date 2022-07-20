package nfsm_test

import (
	"context"
	"testing"

	"github.com/nigzht/nfsm"
)

func TestFactory(t *testing.T) {
	n := nfsm.NewNfsm(context.Background(), "test", nfsm.Handlers{
		"test": func(nfsm nfsm.Machine) (string, error) {
			return "done", nil
		},
		"done": func(nfsm nfsm.Machine) (string, error) {
			return "", nil
		},
	})

	nf := n.Factory(context.Background())

	nf.Context()
}
