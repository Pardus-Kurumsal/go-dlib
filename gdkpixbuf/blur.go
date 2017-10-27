/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package gdkpixbuf

// #cgo pkg-config: gdk-pixbuf-2.0
// #cgo LDFLAGS: -lm
// #include <stdlib.h>
// #include "blur.h"
import "C"

import (
	"fmt"
	"pkg.deepin.io/lib/utils"
)

// BlurImage generate blur effect to an image file.
func BlurImage(srcFile, dstFile string, sigma, numSteps float64, f Format) (err error) {
	srcPixbuf, err := NewPixbufFromFile(srcFile)
	defer FreePixbuf(srcPixbuf)
	if err != nil {
		return
	}
	err = Blur(srcPixbuf, sigma, numSteps)
	if err != nil {
		return
	}
	err = Save(srcPixbuf, dstFile, f)
	return
}

// BlurImageCache generate and save the blurred image file to cache
// directory, if target file already exists, just return it.
func BlurImageCache(srcFile string, sigma, numSteps float64, f Format) (dstFile string, useCache bool, err error) {
	dstFile = generateCacheFilePath(fmt.Sprintf("BlurImageCache%s%f%f%s", srcFile, sigma, numSteps, f))
	if utils.IsFileExist(dstFile) {
		// return cache file
		useCache = true
		return
	}
	err = BlurImage(srcFile, dstFile, sigma, numSteps, f)
	return
}

// Blur generate blur effect to pixbuf object.
func Blur(pixbuf *C.GdkPixbuf, sigma, numSteps float64) (err error) {
	defaultError := fmt.Errorf("blur gdkpixbuf failed, pixbuf=%v, sigma=%v, numSteps=%v", pixbuf, sigma, numSteps)
	ret := C.blur(pixbuf, C.double(sigma), C.double(numSteps))
	if ret == 0 {
		err = defaultError
		return
	}
	return
}
