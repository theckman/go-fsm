// Copyright 2016 Tim Heckman. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package fsm

import "fmt"

// ErrorCode is the type for package-specific error codes. This is used
// within the Error struct, which allows you to programatically determine
// the error cause.
type ErrorCode uint

func (e ErrorCode) String() string {
	switch e {
	case ErrorMachineNotInitialized:
		return "MachineNotInitialized"
	case ErrorTransitionNotPermitted:
		return "TransitionNotPermitted"
	case ErrorStateUndefined:
		return "StateUndefined"
	default:
		return "Unknown"
	}
}

const (
	// ErrorUnknown is the default value
	ErrorUnknown ErrorCode = iota

	// ErrorMachineNotInitialized is an error returned when actions are taken on
	// a machine before it has been initialized. A machine is initialized by
	// adding at least one state and setting it as the initial state.
	ErrorMachineNotInitialized

	// ErrorTransitionNotPermitted is the error returned when trying to
	// transition to an invalid state. In other words, the machine is not
	// permitted to transition from the current state to the one requested.
	ErrorTransitionNotPermitted

	// ErrorStateUndefined is the error returned when the requested state is
	// not defined within the machine.
	ErrorStateUndefined
)

// Error is the struct representing internal errors.
// This implements the error interface
type Error struct {
	message string
	code    ErrorCode
}

// newErrorStruct uses messge and code to create an *Error struct. The *Error
// struct implements the 'error' interface, so it should be able to be used
// wherever 'error' is expected.
func newErrorStruct(message string, code ErrorCode) *Error {
	return &Error{
		message: message,
		code:    code,
	}
}

// Message returns the error message.
func (e *Error) Message() string { return e.message }

// Code returns the error code.
func (e *Error) Code() ErrorCode { return e.code }

func (e *Error) Error() string {
	return fmt.Sprintf("%s (%d): %s", e.code, e.code, e.message)
}
