# nfsm

## Description

nfsm is a type of non-deterministic state machine. Next states are determined by handlers with no pre-determined final state.

## Installation

```
go get github.com/nigzht/nfsm@latest
```

## Example

```go
	nfsm.NewNfsm(context.Background(), "generate", nfsm.Handlers{
		"generate": func(nfsm nfsm.Machine) (string, error) {
			rand.Seed(time.Now().UnixNano())

			nfsm.Metadata().Set("random_number", rand.Intn(2-0)+0)
			return "determine", nil
		},
		"determine": func(nfsm nfsm.Machine) (string, error) {
            fmt.Printf("current state: %s \n previous state: %s", nfsm.Current(), nfsm.Previous())

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
	}).Execute()
```
