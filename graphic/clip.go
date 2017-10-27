/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package graphic

import (
	"fmt"
	"image"
	"image/draw"
)

// ClipImage clip image file that format is recognized and save it with target format.
func ClipImage(srcfile, dstfile string, x, y, w, h int, f Format) (err error) {
	srcimg, err := LoadImage(srcfile)
	if err != nil {
		return
	}
	dstimg := Clip(srcimg, x, y, w, h)
	err = SaveImage(dstfile, dstimg, f)
	dstimg.Pix = nil
	return
}

// ClipImageCache clip any recognized format image and save to cache
// directory, if already exists, just return it.
func ClipImageCache(srcfile string, x, y, w, h int, f Format) (dstfile string, useCache bool, err error) {
	dstfile = generateCacheFilePath(fmt.Sprintf("ClipImageCache%s%d%d%d%d%s", srcfile, x, y, w, h, f))
	if isFileExists(dstfile) {
		useCache = true
		return
	}
	err = ClipImage(srcfile, dstfile, x, y, w, h, f)
	return
}

// ClipImage clip image object.
func Clip(srcimg image.Image, x, y, w, h int) (dstimg *image.RGBA) {
	dstimg = image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(dstimg, dstimg.Bounds(), srcimg, image.Point{x, y}, draw.Src)
	return
}
