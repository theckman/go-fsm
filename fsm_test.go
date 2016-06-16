// Copyright 2016 Tim Heckman. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package fsm

import (
	"runtime"
	"testing"

	. "gopkg.in/check.v1"
)

type TestSuite struct {
	m *Machine
	e *Error
}

var _ = Suite(&TestSuite{})

func Test(t *testing.T) { TestingT(t) }

func (t *TestSuite) setUpMachine(c *C) {
	t.m = &Machine{}
	c.Assert(t.m.AddStateTransitionRules("start", "started", "never_exist"), IsNil)
	c.Assert(t.m.AddStateTransitionRules("started", "finishing", "aborted"), IsNil)
	c.Assert(t.m.AddStateTransitionRules("aborted"), IsNil)
	c.Assert(t.m.AddStateTransitionRules("finishing", "finished"), IsNil)
	c.Assert(t.m.AddStateTransitionRules("finished"), IsNil)
}

func (t *TestSuite) setUpError(c *C) {
	t.e = &Error{
		message: "testMessage",
		code:    ^ErrorCode(0),
	}
}

func (t *TestSuite) SetUpSuite(c *C) {
	t.setUpMachine(c)
	t.setUpError(c)
}

type testCallback struct {
	prevState, curState State
}

func (tc *testCallback) StateTransitionCallback(s State) error {
	tc.prevState = tc.curState
	tc.curState = s
	return nil
}

func (t *TestSuite) TestMachine_AddStateTransitionRules(c *C) {
	var ok bool
	var err error
	var trs TransitionRuleSet

	// reset the machine
	defer t.setUpMachine(c)

	//
	// Test adding a state that allow multiple transitions
	//
	err = t.m.AddStateTransitionRules("testing", "passing", "failing")
	c.Assert(err, IsNil)

	trs, err = t.m.StateTransitionRules("testing")
	c.Assert(err, IsNil)
	c.Assert(len(trs), Equals, 2)

	_, ok = trs["passing"]
	c.Check(ok, Equals, true)

	_, ok = trs["failing"]
	c.Check(ok, Equals, true)

	//
	// Test adding a state with no transitions works
	//
	err = t.m.AddStateTransitionRules("passing")
	c.Assert(err, IsNil)

	trs, err = t.m.StateTransitionRules("passing")
	c.Assert(err, IsNil)
	c.Check(len(trs), Equals, 0)
}

func (t *TestSuite) TestMachine_StateTransition(c *C) {
	var ok bool
	var err error
	var state State
	var fsmErr *Error

	// reset the machine
	defer t.setUpMachine(c)

	//
	// Test that when setting initial state, the requested state
	// must have already been registered
	//
	err = t.m.StateTransition("does_not_exist_shouldn't_ever_exist")
	c.Assert(err, NotNil)

	fsmErr, ok = err.(*Error)
	c.Assert(ok, Equals, true)
	c.Check(fsmErr.Code(), Equals, ErrorStateUndefined)

	//
	// Test that setting the inital state to a valid state succeeds
	//
	err = t.m.StateTransition("start")
	c.Assert(err, IsNil)

	state = t.m.CurrentState()
	c.Check(state, Equals, State("start"))

	//
	// Test that trying to transition to a state, that we are NOT
	// permitted to transition to, fails to transition
	//
	err = t.m.StateTransition("finished")
	c.Assert(err, NotNil)

	fsmErr, ok = err.(*Error)
	c.Assert(ok, Equals, true)
	c.Check(fsmErr.Code(), Equals, ErrorTransitionNotPermitted)

	//
	// Test that trying to transition to a state, that is permitted yet
	// hasn't been registered, fails to transition
	//
	err = t.m.StateTransition("never_exist")
	c.Assert(err, NotNil)

	fsmErr, ok = err.(*Error)
	c.Assert(ok, Equals, true)
	c.Check(fsmErr.Code(), Equals, ErrorStateUndefined)

	//
	// Test that transitioning to a valid state works
	//
	err = t.m.StateTransition("started")
	c.Assert(err, IsNil)

	state = t.m.CurrentState()
	c.Check(state, Equals, State("started"))

	//
	// Test that an uninitialized machine errors
	//
	machine := &Machine{}

	err = machine.StateTransition("starting")
	c.Assert(err, NotNil)

	fsmErr, ok = err.(*Error)
	c.Assert(ok, Equals, true)
	c.Check(fsmErr.Code(), Equals, ErrorMachineNotInitialized)

	//
	// Test that the callback is called synchronously
	//
	tc := &testCallback{}

	err = t.m.SetStateTransitionCallback(tc, true)
	c.Assert(err, IsNil)

	err = t.m.StateTransition("finishing")
	c.Assert(err, IsNil)

	state = t.m.CurrentState()
	c.Check(state, Equals, State("finishing"))

	c.Check(tc.curState, Equals, State("finishing"))
	c.Check(tc.prevState, Equals, State(""))

	//
	// Test that the callback is called asynchronously
	//
	err = t.m.SetStateTransitionCallback(tc, false)
	c.Assert(err, IsNil)

	err = t.m.StateTransition("finished")
	c.Assert(err, IsNil)

	// ensure callbacks haven't execute, yet (might be racy)
	c.Check(tc.curState, Equals, State("finishing"))
	c.Check(tc.prevState, Equals, State(""))

	state = t.m.CurrentState()
	c.Check(state, Equals, State("finished"))

	// schedule any goroutines to run
	runtime.Gosched()

	c.Check(tc.curState, Equals, State("finished"))
	c.Check(tc.prevState, Equals, State("finishing"))
}

func (t *TestSuite) TestMachine_CurrentState(c *C) {
	var state State

	// reset the machine
	defer t.setUpMachine(c)

	state = t.m.CurrentState()
	c.Check(state, Equals, State(""))

	err := t.m.StateTransition("start")
	c.Assert(err, IsNil)

	state = t.m.CurrentState()
	c.Check(state, Equals, State("start"))
}

func (t *TestSuite) TestMachine_StateTransitionRules(c *C) {
	var ok bool
	var err error
	var trs TransitionRuleSet
	var fsmErr *Error

	//
	// Test that retrieving a state with multiple valid destinations
	// returns all states
	//
	trs, err = t.m.StateTransitionRules("started")
	c.Assert(err, IsNil)
	c.Assert(len(trs), Equals, 2)

	_, ok = trs["finishing"]
	c.Check(ok, Equals, true)

	_, ok = trs["aborted"]
	c.Check(ok, Equals, true)

	//
	// Test that retrieving a state with no valid destinations returns
	// and empty TransitionRuleSet
	//
	trs, err = t.m.StateTransitionRules("aborted")
	c.Assert(err, IsNil)
	c.Check(len(trs), Equals, 0)

	//
	// Test that retreiving an unregistered state returns an error
	//
	trs, err = t.m.StateTransitionRules("never_exist")
	c.Assert(err, NotNil)
	c.Check(len(trs), Equals, 0)

	fsmErr, ok = err.(*Error)
	c.Assert(ok, Equals, true)
	c.Check(fsmErr.Code(), Equals, ErrorStateUndefined)

	//
	// Test that an uninitialized machine errors
	//
	machine := &Machine{}

	trs, err = machine.StateTransitionRules("")
	c.Assert(err, NotNil)
	c.Check(len(trs), Equals, 0)

	fsmErr, ok = err.(*Error)
	c.Assert(ok, Equals, true)
	c.Check(fsmErr.Code(), Equals, ErrorMachineNotInitialized)
}

func (t *TestSuite) TestMission_SetStateTransitionCallback(c *C) {
	var err error

	// reset the machine
	defer t.setUpMachine(c)

	tc := &testCallback{}

	err = t.m.SetStateTransitionCallback(tc, true)
	c.Assert(err, IsNil)
	c.Check(t.m.callback, DeepEquals, tc)
}
