/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package group

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

func (*testWrapper) TestGetGroupByName(c *C.C) {
	name := "root"
	group, err := GetGroupByName(name)
	c.Check(err, C.IsNil)
	c.Check(group.Gid, C.Equals, uint32(0))

	name = "root2"
	group, err = GetGroupByName(name)
	c.Check(group, C.IsNil)
	c.Check(err, C.DeepEquals, &GroupNotFoundError{Name: name})
}

func (*testWrapper) TestGetGroupByGid(c *C.C) {
	uid := uint32(0)
	group, err := GetGroupByGid(uid)
	c.Check(err, C.IsNil)
	c.Check(group.Name, C.Equals, "root")
}

func (*testWrapper) TestGetGroupEntry(c *C.C) {
	groups := GetGroupEntry()
	c.Check(len(groups), C.Not(C.Equals), 0)
	c.Check(groups[0].Name, C.Equals, "root")
}
