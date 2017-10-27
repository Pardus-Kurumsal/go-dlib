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
	"image/draw"
	"sort"
)

type compositeInfoSorter []CompositeInfo

func (a compositeInfoSorter) Len() int           { return len(a) }
func (a compositeInfoSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a compositeInfoSorter) Less(i, j int) bool { return a[i].Z < a[j].Z }

type CompositeInfo struct {
	File    string
	X, Y, Z int
}

// CompositeImageSet composite a set of image files.
func CompositeImageSet(srcfile string, compinfos []CompositeInfo, dstfile string, f Format) (err error) {
	srcimg, err := LoadImage(srcfile)
	if err != nil {
		return
	}
	dstimg := convertToRGBA(srcimg)

	sort.Sort(compositeInfoSorter(compinfos))
	for _, compinfo := range compinfos {
		var compimg image.Image
		compimg, err = LoadImage(compinfo.File)
		if err != nil {
			return
		}
		Composite(dstimg, compimg, compinfo.X, compinfo.Y)
	}
	err = SaveImage(dstfile, dstimg, f)
	dstimg.Pix = nil
	return
}

// CompositeImage composite two image files.
func CompositeImage(srcfile, compfile, dstfile string, x, y int, f Format) (err error) {
	srcimg, err := LoadImage(srcfile)
	if err != nil {
		return
	}
	compimg, err := LoadImage(compfile)
	if err != nil {
		return
	}
	dstimg := convertToRGBA(srcimg)
	Composite(dstimg, compimg, x, y)
	err = SaveImage(dstfile, dstimg, f)
	dstimg.Pix = nil
	return
}

// CompositeImageUri composite two images which format in data uri.
func CompositeImageUri(srcDatauri, compDataUri string, x, y int, f Format) (dstDataUri string, err error) {
	srcimg, err := LoadImageFromDataUri(srcDatauri)
	if err != nil {
		return
	}
	compimg, err := LoadImageFromDataUri(compDataUri)
	if err != nil {
		return
	}
	dstimg := convertToRGBA(srcimg)
	Composite(dstimg, compimg, x, y)
	dstDataUri, err = ConvertImageObjectToDataUri(dstimg, f)
	dstimg.Pix = nil
	return
}

// CompositeImage composite two image objects.
func Composite(dstimg draw.Image, compimg image.Image, x, y int) {
	w, h := GetSize(compimg)
	r := image.Rect(x, y, x+w, y+h)
	draw.Draw(dstimg, r, compimg, image.Point{0, 0}, draw.Over)
	return
}
