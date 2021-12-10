package core

import (
	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/node"
	"github.com/speedata/boxesandglue/pdfbackend/pdf"
)

var (
	loadedImages = make(map[string]*pdf.Imagefile)
)

func createImageHlist(xd *xtsDocument, imgfile *pdf.Imagefile) *node.HList {
	ii := xd.doc.CreateImage(imgfile)
	imgNode := node.NewImage()
	imgNode.Img = ii
	imgNode.Width = bag.MustSp("4cm")
	imgNode.Height = bag.MustSp("3cm")
	hlist := node.Hpack(imgNode)
	return hlist
}
