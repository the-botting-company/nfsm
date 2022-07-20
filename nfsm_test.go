package nfsm_test

import (
	"context"
	"testing"

	"github.com/nigzht/nfsm"
)

func TestFactoryPurity(t *testing.T) {
	n := nfsm.NewNfsm(context.Background(), "test", nfsm.Handlers{
		"test": func(nfsm nfsm.Machine) (string, error) {
			return "done", nil
		},
		"done": func(nfsm nfsm.Machine) (string, error) {
			return "", nil
		},
	})

	nf := n.Factory(context.Background())

	n.Metadata().Set("test", "data")

	if v := nf.Metadata().Get("test"); v != nil {
		t.Errorf("Factory copy inherited data")
	}
}
