# nfsm

## Description

nfsm is a type of non-deterministic state machine. Next states are determined by handlers with no pre-determined final state. The API docs can be found here: https://pkg.go.dev/github.com/nigzht/nfsm

## Installation

```
go get github.com/the-botting-company/nfsm@latest
```

## Example

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/the-botting-company/nfsm"
)

func main() {
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
		fmt.Errorf("%v", err)
	}
}

```
