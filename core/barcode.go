package core

import (
	"image/color"

	"github.com/speedata/barcode"
	"github.com/speedata/barcode/code128"
	"github.com/speedata/barcode/ean"
	"github.com/speedata/barcode/qr"
	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/node"
	"github.com/speedata/boxesandglue/frontend"
	"github.com/speedata/boxesandglue/frontend/pdfdraw"
)

const (
	barcodeEAN13 int = iota
	barcodeCode128
	barcodeQR
)

func createBarcode(typ int, value string, width bag.ScaledPoint, height bag.ScaledPoint, xd *xtsDocument, ff *frontend.FontFamily, fontsize bag.ScaledPoint, useText bool) (node.Node, error) {
	var n node.Node
	switch typ {
	case barcodeQR:
		useText = false
	}
	var vl *node.VList
	var err error
	if height == 0 {
		height = width
	}
	if useText {
		txt := frontend.NewText()
		txt.Settings[frontend.SettingSize] = fontsize
		txt.Settings[frontend.SettingFontFamily] = ff
		txt.Settings[frontend.SettingHAlign] = frontend.HAlignCenter
		txt.Items = append(txt.Items, value)
		vl, _, err = xd.document.FormatParagraph(txt, width)
		if err != nil {
			return nil, err
		}
		height -= vl.Height
	}
	switch typ {
	case barcodeEAN13:
		bc, err := ean.Encode(value)
		if err != nil {
			return nil, err
		}
		if n, err = barcodeCreate(bc, width, height, xd); err != nil {
			return nil, err
		}
	case barcodeCode128:
		bc, err := code128.Encode(value)
		if err != nil {
			return nil, err
		}
		if n, err = barcodeCreate(bc, width, height, xd); err != nil {
			return nil, err
		}
	case barcodeQR:
		bc, err := qr.Encode(value, qr.M, qr.Auto)
		if err != nil {
			return nil, err
		}
		return barcodeCreate3d(bc, width, xd)
	}
	if useText {
		n = node.InsertAfter(n, n, vl)
	}
	vl = node.Vpack(n)
	vl.Attributes = node.H{
		"origin": "barcode",
	}

	return vl, nil
}

func barcodeCreate(bc barcode.Barcode, width bag.ScaledPoint, height bag.ScaledPoint, xd *xtsDocument) (node.Node, error) {
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
			d.Rect(curX, 0, wdBar, -height)
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
	rule.Height = height
	rule.Width = width
	rule.Hide = true
	rule.Post = pdfdraw.New().Restore().String()
	return rule, nil
}

func barcodeCreate3d(bc barcode.Barcode, width bag.ScaledPoint, xd *xtsDocument) (node.Node, error) {
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
