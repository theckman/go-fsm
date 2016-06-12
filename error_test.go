// Copyright 2016 Tim Heckman. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package fsm

import (
	"fmt"

	. "gopkg.in/check.v1"
)

func (*TestSuite) TestErrorCode_String(c *C) {
	c.Check(ErrorUnknown.String(), Equals, "Unknown")
	c.Check(ErrorMachineNotInitialized.String(), Equals, "MachineNotInitialized")
	c.Check(ErrorTransitionNotPermitted.String(), Equals, "TransitionNotPermitted")
	c.Check(ErrorStateUndefined.String(), Equals, "StateUndefined")
	c.Check((ErrorUnknown + 100).String(), Equals, "Unknown")
}

func (t *TestSuite) TestError_Message(c *C) {
	var str string

	str = t.e.Message()
	c.Check(str, Equals, "testMessage")
}

func (t *TestSuite) TestError_Code(c *C) {
	var code ErrorCode

	code = t.e.Code()
	c.Check(code, Equals, ^ErrorCode(0))
}

func (t *TestSuite) TestError_Error(c *C) {
	var errStr string

	expected := fmt.Sprintf("Unknown (%d): testMessage", ^uint(0))

	errStr = t.e.Error()
	c.Check(errStr, Equals, expected)
}
