package core

import (
	"math"

	pdf "github.com/speedata/baseline-pdf"
	"github.com/speedata/boxesandglue/backend/bag"
	bagimage "github.com/speedata/boxesandglue/backend/image"
	"github.com/speedata/boxesandglue/backend/node"
)

func createImageHlist(xd *xtsDocument, width, height, minwidth, maxwidth, minheight, maxheight *bag.ScaledPoint, stretch bool, imgfile *pdf.Imagefile, pagenumer int) *node.HList {
	if pagenumer == 0 {
		pagenumer = 1
	}
	ii := xd.document.Doc.CreateImage(imgfile, pagenumer)

	wd, ht := calculateImageSize(ii, width, height, minwidth, maxwidth, minheight, maxheight, stretch)
	imgNode := node.NewImage()
	imgNode.Img = ii
	imgNode.Width = wd
	imgNode.Height = ht
	hlist := node.Hpack(imgNode)
	return hlist
}

var posInf = math.Inf(1)

// calculateImageSize calculates the width and the height from the given
// parameters. Unset parameters should be nil or 0 for min.. and MaxSP for
// max... See https://www.w3.org/TR/CSS2/visudet.html#min-max-widths for rules.
func calculateImageSize(img *bagimage.Image, requestedWidth, requestedHeight, minWidth, maxWidth, minHeight, maxHeight *bag.ScaledPoint, stretch bool) (bag.ScaledPoint, bag.ScaledPoint) {
	// Constraint Violation                                                           Resolved Width                      Resolved Height
	// ===================================================================================================================================================
	//  1 none                                                                        w                                   h
	//  2 w > max-width                                                               max-width                           max(max-width * h/w, min-height)
	//  3 w < min-width                                                               min-width                           min(min-width * h/w, max-height)
	//  4 h > max-height                                                              max(max-height * w/h, min-width)    max-height
	//  5 h < min-height                                                              min(min-height * w/h, max-width)    min-height
	//  6 (w > max-width) and (h > max-height), where (max-width/w ≤ max-height/h)    max-width                           max(min-height, max-width * h/w)
	//  7 (w > max-width) and (h > max-height), where (max-width/w > max-height/h)    max(min-width, max-height * w/h)    max-height
	//  8 (w < min-width) and (h < min-height), where (min-width/w ≤ min-height/h)    min(max-width, min-height * w/h)    min-height
	//  9 (w < min-width) and (h < min-height), where (min-width/w > min-height/h)    min-width                           min(max-height, min-width * h/w)
	// 10 (w < min-width) and (h > max-height)                                        min-width                           max-height
	// 11 (w > max-width) and (h < min-height)                                        max-width                           min-height

	var imgwd, imght, wd, ht, minwd, maxwd, minht, maxht float64

	imgwd = img.Width.ToPT()
	imght = img.Height.ToPT()

	if requestedWidth == nil {
		wd = img.Width.ToPT()
	} else {
		wd = requestedWidth.ToPT()
	}

	if requestedHeight == nil {
		ht = img.Height.ToPT()
	} else {
		ht = requestedHeight.ToPT()
	}

	if minWidth == nil {
		minwd = 0
	} else {
		minwd = minWidth.ToPT()
	}

	if maxWidth == nil || *maxWidth == bag.MaxSP {
		maxwd = posInf
	} else {
		maxwd = maxWidth.ToPT()
	}

	if minHeight == nil {
		minht = 0
	} else {
		minht = minHeight.ToPT()
	}

	if maxHeight == nil || *maxHeight == bag.MaxSP {
		maxht = posInf
	} else {
		maxht = maxHeight.ToPT()
	}
	// if stretch and max{height,width} then the image should grow as needed
	if stretch && maxht < posInf && maxwd < posInf {
		stretchamount := math.Min(maxwd/imgwd, maxht/imght)
		if stretchamount > 1 {
			return bag.ScaledPointFromFloat(imgwd * stretchamount), bag.ScaledPointFromFloat(imght * stretchamount)
		}
	}

	// If one of height or width is given, the other one should
	// be adjusted to keep the aspect ratio
	if ht == imght {
		if wd != imgwd {
			ht = ht * wd / imgwd
		}
	} else if wd == imgwd {
		if ht != imght {
			wd = wd * ht / imght
		}
	}
	if wd < minwd && ht > maxht {
		// fmt.Println("10")
		wd = minwd
		ht = maxht
	} else if wd > maxwd && ht < minht {
		// fmt.Println("11")
		wd = maxwd
		ht = minht
	} else if wd > maxwd && ht > maxht && maxwd/wd <= maxht/ht {
		// fmt.Println("6")
		ht = math.Max(minht, maxwd*ht/wd)
		wd = maxwd
	} else if wd > maxwd && ht > maxht && maxwd/wd > maxht/ht {
		// fmt.Println("7")
		wd = math.Max(minwd, maxht*wd/ht)
		ht = maxht
	} else if wd < minwd && ht < minht && minwd/wd <= minht/ht {
		// fmt.Println("8")
		wd = math.Min(maxwd, minht*wd/ht)
		ht = minht
	} else if wd < minwd && ht < minht && minwd/wd > minht/ht {
		// fmt.Println("9")
		wd = minwd
		ht = math.Min(maxht, minwd*ht/wd)
	} else if wd > maxwd {
		// fmt.Println("2")
		ht = math.Max(maxwd*ht/wd, minht)
		wd = maxwd
	} else if wd < minwd {
		// fmt.Println("3")
		ht = math.Min(minwd*ht/wd, maxht)
		wd = minwd
	} else if ht > maxht {
		// fmt.Println("4")
		wd = math.Max(maxht*wd/ht, minwd)
		ht = maxht
	} else if ht < minht {
		// fmt.Println("5")
		wd = math.Min(minht*wd/ht, maxwd)
	}

	return bag.ScaledPointFromFloat(wd), bag.ScaledPointFromFloat(ht)
}
