/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package passwd

import (
	"testing"

	C "launchpad.net/gocheck"
)

type testWrapper struct{}

func init() {
	C.Suite(&testWrapper{})
}

func Test(t *testing.T) {
	C.TestingT(t)
}

func (*testWrapper) TestGetPasswdByName(c *C.C) {
	name := "root"
	passwd, err := GetPasswdByName(name)
	c.Check(passwd.Home, C.Equals, "/root")
	c.Check(passwd.Uid, C.Equals, uint32(0))
	c.Check(err, C.IsNil)

	name = "root2"
	passwd, err = GetPasswdByName(name)
	c.Check(passwd, C.IsNil)
	c.Check(err, C.DeepEquals, &UserNotFoundError{Name: name})
}

func (*testWrapper) TestGetPasswdByUid(c *C.C) {
	uid := uint32(0)
	passwd, err := GetPasswdByUid(uid)
	c.Check(err, C.IsNil)
	c.Check(passwd.Name, C.Equals, "root")
}

func (*testWrapper) TestGetPasswdEntry(c *C.C) {
	passwds := GetPasswdEntry()
	c.Check(len(passwds), C.Not(C.Equals), 0)
	c.Check(passwds[0].Name, C.Equals, "root")
}
