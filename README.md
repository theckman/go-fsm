# go-fsm
[![TravisCI Build Status](https://img.shields.io/travis/theckman/go-fsm/master.svg?style=flat)](https://travis-ci.org/theckman/go-fsm)
[![GoDoc](https://img.shields.io/badge/go--fsm-GoDoc-blue.svg?style=flat)](https://godoc.org/github.com/theckman/go-fsm)

`fsm` is a finite state machine written in Go. This package aims to be very simple, while making
it easy to maintain predictable state.

## License
This code is released for use under the [MIT License](https://tldrlegal.com/license/mit-license),
of which the full contents can be found within the [LICENSE](https://github.com/theckman/go-fsm/blob/master/LICENSE)
file.

This license is extremely permissve, and should allow the full use of this code in almost all situations.
When in doubt, consult a lawyer.

## Installation
```
go get github.com/theckman/go-fsm
```

## Usage
The purpose of the machine is to be a field within your struct. You then interact with the machine
to enforce the state of struct. For the full API documentation, check out the
[GoDoc](https://godoc.org/github.com/theckman/go-fsm) page. Here's an example
of a machine with some states you can transition between:

```Go
import "github.com/theckman/go-fsm"

type T struct {
	M *fsm.Machine
}

func main() {
	t := &T{M: &fsm.Machine{}}

	// add initial rule
	err := t.M.AddStateTransitionRules("started", "finished", "aborted", "exited")

	if err != nil {
		// handle
	}

	// add rest of rules
	t.M.AddStateTransitionRules("aborted", "started")
	t.M.AddStateTransitionRules("finished", "started")
	t.M.AddStateTransitionRules("exited") // final state

	// set initial state
	err = t.M.StateTransition("aborted") // nil

	// get the current state
	state := t.M.CurrentState() // "aborted"

	// try to transition to an non-whitelisted state
	err = t.M.StateTransition("finished") // ErrTransitionNotPermitted

	// try to transition to a permitted state
	err = t.M.StateTransition("started") // nil
}
```
