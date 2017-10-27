/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package lunar

import (
	"testing"
)

func Test_CalcEarthLongitudeNutation(t *testing.T) {
	n := CalcEarthLongitudeNutation(1.2345)
	t.Log(n)
}

func Test_CalcEarthObliquityNutation(t *testing.T) {
	n := CalcEarthObliquityNutation(1.2345)
	t.Log(n)
}
