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
	"image"
)

// FlipImageHorizontal flip image in horizontal direction, and save as
// target format.
func FlipImageHorizontal(srcfile, dstfile string, f Format) (err error) {
	srcimg, err := LoadImage(srcfile)
	if err != nil {
		return err
	}
	dstimg := doFlipImageHorizontal(srcimg)
	err = SaveImage(dstfile, dstimg, f)
	dstimg.Pix = nil
	return
}

// FlipImageVertical  flip image  in vertical  direction, and  save as
// target format.
func FlipImageVertical(srcfile, dstfile string, f Format) (err error) {
	srcimg, err := LoadImage(srcfile)
	if err != nil {
		return err
	}
	dstimg := doFlipImageVertical(srcimg)
	err = SaveImage(dstfile, dstimg, f)
	dstimg.Pix = nil
	return
}

func doFlipImageHorizontal(srcimg image.Image) (dstimg *image.RGBA) {
	w, h := GetSize(srcimg)
	dstimg = image.NewRGBA(image.Rect(0, 0, w, h))

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			dstimg.Set(w-x-1, y, srcimg.At(x, y))
		}
	}

	return
}

func doFlipImageVertical(srcimg image.Image) (dstimg *image.RGBA) {
	w, h := GetSize(srcimg)
	dstimg = image.NewRGBA(image.Rect(0, 0, w, h))

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			dstimg.Set(x, h-y-1, srcimg.At(x, y))
		}
	}

	return
}
