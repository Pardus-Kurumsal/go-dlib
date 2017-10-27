/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dbus

import "testing"

func TestSplitObjectPath(t *testing.T) {
	data := []struct {
		Path   ObjectPath
		Parent ObjectPath
		Base   string
	}{
		{"/com/deepin/Test", "/com/deepin", "Test"},
		{"/com/deepin/Test/0", "/com/deepin/Test", "0"},
		{"/com/deepin/Test/a", "/com/deepin/Test", "a"},
		{"/com/deepin/Test/abc", "/com/deepin/Test", "abc"},
		{"/", "/", ""},
		{"/abc", "/", "abc"},
		{"/com/deepin/Test/", "", ""},
		{"", "", ""},
	}

	for _, i := range data {
		p, b := splitObjectPath(i.Path)
		if p != i.Parent || b != i.Base {
			t.Errorf("splitObjectPath:(%q) get (%q,%q). It should be (%q, %q)", i.Path, p, b, i.Parent, i.Base)
		}
	}
}
