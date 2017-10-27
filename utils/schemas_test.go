/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package utils

import C "launchpad.net/gocheck"

type schemaTest struct {
	schema string
	exist  bool
}

func (*testWrapper) TestGSchemaIsExist(c *C.C) {
	schemas := []string{
		"com.deepin.wacom",
		"com.deepin.touchpad",
	}

	list1 := []schemaTest{
		schemaTest{"com.deepin.wacom", true},
		schemaTest{"com.deepin.touchpad", true},
	}

	list2 := []schemaTest{
		schemaTest{"org.123.123", false},
		schemaTest{"org/11/11", false},
		schemaTest{"sdsdsvfdsfs", false},
		schemaTest{"/dsfd/assasd", false},
		schemaTest{".sds.sadsd.", false},
		schemaTest{"-sds-sds-ss", false},
		schemaTest{"(jjjj)", false},
		schemaTest{"$fgg$", false},
	}

	for _, t := range list1 {
		c.Check(isSchemaInList(t.schema, schemas), C.Equals, t.exist)
	}

	for _, t := range list2 {
		c.Check(IsGSchemaExist(t.schema), C.Equals, t.exist)
		c.Check(isSchemaInList(t.schema, schemas), C.Equals, t.exist)
	}
}
