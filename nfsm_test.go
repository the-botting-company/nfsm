package nfsm_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/the-botting-company/nfsm"
)

func TestExample(t *testing.T) {
	h := nfsm.Handlers{
		"generate": func(nfsm nfsm.Machine) (string, error) {
			rand.Seed(time.Now().UnixNano())

			nfsm.Metadata().Set("random_number", rand.Intn(2-0)+0)
			return "determine", nil
		},
		"determine": func(nfsm nfsm.Machine) (string, error) {
			fmt.Printf("current state: %s \nprevious state: %s \n", nfsm.Current(), nfsm.Previous())

			n := nfsm.Metadata().Get("random_number")

			n, ok := n.(int)
			if !ok {
				return "", errors.New("unexpected type")
			}

			fmt.Printf("random number: %d \n", n)

			if n == 0 {
				return "generate", nil
			}

			return "", nil
		},
	}

	if err := nfsm.NewNfsm(context.Background(), nfsm.NewFlow("generate", h)).Execute(); err != nil {
		t.Errorf("%v", err)
	}
}

func TestLinearFlow(t *testing.T) {

	h := nfsm.Handlers{
		"1": func(nfsm nfsm.Machine) (string, error) {
			nfsm.Metadata().Set("steps", []string{nfsm.Current()})
			return "2", nil
		},
		"2": func(nfsm nfsm.Machine) (string, error) {
			s := nfsm.Metadata().Get("steps")

			s = append(s.([]string), nfsm.Current())

			nfsm.Metadata().Set("steps", s)
			return "3", nil
		},
		"3": func(nfsm nfsm.Machine) (string, error) {
			s := nfsm.Metadata().Get("steps")

			s = append(s.([]string), nfsm.Current())

			nfsm.Metadata().Set("steps", s)
			return "4", nil
		},
		"4": func(nfsm nfsm.Machine) (string, error) {
			s := nfsm.Metadata().Get("steps")

			s = append(s.([]string), nfsm.Current())

			nfsm.Metadata().Set("steps", s)
			return "5", nil
		},
		"5": func(nfsm nfsm.Machine) (string, error) {
			s := nfsm.Metadata().Get("steps")

			s = append(s.([]string), nfsm.Current())

			if !reflect.DeepEqual(s, []string{"1", "2", "3", "4", "5"}) {
				return "", errors.New("steps are not equal")
			}

			return "", nil
		},
	}

	if err := nfsm.NewNfsm(context.Background(), nfsm.NewFlow("1", h)).Execute(); err != nil {
		t.Errorf("%v", err)
	}
}

func TestCancel(t *testing.T) {

	h := nfsm.Handlers{
		"1": func(nfsm nfsm.Machine) (string, error) {
			for {
				select {
				case <-nfsm.Context().Done():
					return "", nil
				case <-time.After(1 * time.Second):
					return "2", nil
				}
			}
		},
		"2": func(nfsm nfsm.Machine) (string, error) {
			return "", errors.New("reached undesired state")
		},
	}

	ctx, c := context.WithCancel(context.Background())

	n := nfsm.NewNfsm(ctx, nfsm.NewFlow("1", h))

	go func() {
		if err := n.Execute(); err != nil {
			t.Errorf("%v", err)
		}
	}()

	<-time.After(950 * time.Millisecond)

	c()
}
