package core

import (
	"image/color"

	"github.com/speedata/barcode"
	"github.com/speedata/barcode/code128"
	"github.com/speedata/barcode/ean"
	"github.com/speedata/barcode/qr"
	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/node"
	"github.com/speedata/boxesandglue/frontend/pdfdraw"
)

const (
	barcodeEAN13 int = iota
	barcodeCode128
	barcodeQR
)

func createBarcode(typ int, value string, width bag.ScaledPoint, xd *xtsDocument) (node.Node, error) {
	switch typ {
	case barcodeEAN13:
		bc, err := ean.Encode(value)
		if err != nil {
			return nil, err
		}
		return barcodeCreate(bc, value, width, xd)
	case barcodeCode128:
		bc, err := code128.Encode(value)
		if err != nil {
			return nil, err
		}
		return barcodeCreate(bc, value, width, xd)
	case barcodeQR:
		bc, err := qr.Encode(value, qr.M, qr.Auto)
		if err != nil {
			return nil, err
		}
		return barcodeCreate3d(bc, value, width, xd)

	}

	return nil, nil
}

func barcodeCreate(bc barcode.Barcode, value string, width bag.ScaledPoint, xd *xtsDocument) (node.Node, error) {
	dx := bc.Bounds().Dx()
	wdBar := width / bag.ScaledPoint(dx)
	bgcolor := xd.document.GetColor("black")

	d := pdfdraw.New()
	d.Save()
	d.Color(*bgcolor)
	curX := bag.ScaledPoint(0)
	i := 0
	for {
		at := bc.At(i, 0)
		col, _, _, _ := at.RGBA()
		if col == 0 {
			d.Rect(curX, 0, wdBar, -60*bag.Factor)
			d.Fill()
		}
		curX += wdBar
		i++
		if i == dx {
			break
		}
	}

	rule := node.NewRule()
	rule.Pre = d.String()
	rule.Hide = true
	rule.Post = pdfdraw.New().Restore().String()
	vl := node.Vpack(rule)
	if vl.Attributes == nil {
		vl.Attributes = node.H{}
	}
	vl.Attributes = node.H{
		"origin": "barcode",
	}
	vl.Width = width
	vl.Height = 60 * bag.Factor

	return vl, nil
}

func barcodeCreate3d(bc barcode.Barcode, value string, width bag.ScaledPoint, xd *xtsDocument) (node.Node, error) {
	dx := bc.Bounds().Dx()
	dy := bc.Bounds().Dy()

	wdRect := width / bag.ScaledPoint(dx)
	bgcolor := xd.document.GetColor("black")

	curX := bag.ScaledPoint(0)
	curY := bag.ScaledPoint(-wdRect)

	d := pdfdraw.New()
	d.Save()
	d.Color(*bgcolor)
	delta := bag.Factor / 100

	for y := 0; y < dy; y++ {
		curX = 0
		for x := 0; x < dx; x++ {
			col := bc.At(x, y).(color.Gray16)
			if col.Y == 0 {
				d.Rect(curX-delta, curY-delta, wdRect+2*delta, wdRect+2*delta).Fill()
			}
			curX += wdRect
		}
		curY -= wdRect
	}
	rule := node.NewRule()
	rule.Pre = d.String()
	rule.Hide = true
	rule.Post = pdfdraw.New().Restore().String()
	vl := node.Vpack(rule)
	if vl.Attributes == nil {
		vl.Attributes = node.H{}
	}
	vl.Attributes = node.H{
		"origin": "barcode",
	}
	vl.Width = width
	vl.Height = width

	return vl, nil
}
