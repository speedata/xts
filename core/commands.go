package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/node"
	"github.com/speedata/boxesandglue/document"
	"github.com/speedata/boxesandglue/pdfbackend/pdf"
	"github.com/speedata/goxml"
	"github.com/speedata/goxpath/xpath"
)

type commandFunc func(*xtsDocument, *goxml.Element, *xpath.Parser) (xpath.Sequence, error)

var (
	dataDispatcher = make(map[string]map[string]*goxml.Element)
	dispatchTable  map[string]commandFunc
)

func init() {
	dispatchTable = map[string]commandFunc{
		"B":                cmdB,
		"Color":            cmdColor,
		"DefineFontfamily": cmdDefineFontfamily,
		"DefineFontsize":   cmdDefineFontsize,
		"Image":            cmdImage,
		"LoadFontfile":     cmdLoadFontfile,
		"Options":          cmdOptions,
		"Paragraph":        cmdParagraph,
		"PlaceObject":      cmdPlaceObject,
		"Record":           cmdRecord,
		"SetGrid":          cmdSetGrid,
		"Textblock":        cmdTextblock,
		"Trace":            cmdTrace,
		"Value":            cmdValue,
	}
}

func dispatch(xd *xtsDocument, layoutelement *goxml.Element, data *xpath.Parser) (xpath.Sequence, error) {
	var retSequence xpath.Sequence
	for _, cld := range layoutelement.Children() {
		if elt, ok := cld.(*goxml.Element); ok {
			if f, ok := dispatchTable[elt.Name]; ok {
				bag.Logger.Debugf("Call %s (line %d)", elt.Name, elt.Line)
				seq, err := f(xd, elt, data)
				if err != nil {
					return nil, err
				}
				retSequence = append(retSequence, seq...)
			} else {
				bag.Logger.DPanicf("dispatch: element %q unknown", elt.Name)
			}
		}
	}
	return retSequence, nil
}

func cmdB(xd *xtsDocument, layoutelt *goxml.Element, dataelt *xpath.Parser) (xpath.Sequence, error) {
	seq, err := dispatch(xd, layoutelt, dataelt)

	te := &document.TypesettingElement{
		Settings: document.TypesettingSettings{
			document.SettingFontWeight: 700,
		},
	}
	getTextvalues(te, seq, "cmdBold")

	return xpath.Sequence{te}, err
}

func cmdColor(xd *xtsDocument, layoutelt *goxml.Element, dataelt *xpath.Parser) (xpath.Sequence, error) {
	var err error
	var colorname string
	if colorname, err = xd.getAttributeString("name", layoutelt, true, true, ""); err != nil {
		return nil, err
	}

	seq, err := dispatch(xd, layoutelt, dataelt)

	te := &document.TypesettingElement{
		Settings: document.TypesettingSettings{
			document.SettingColor: colorname,
		},
	}
	getTextvalues(te, seq, "cmdColor")

	return xpath.Sequence{te}, err
}

func cmdDefineFontfamily(xd *xtsDocument, layoutelt *goxml.Element, dataelt *xpath.Parser) (xpath.Sequence, error) {
	var err error
	var name, fontface string
	if name, err = xd.getAttributeString("name", layoutelt, true, true, ""); err != nil {
		return nil, err
	}

	ff := xd.doc.NewFontFamily(name)

	for _, cld := range layoutelt.Children() {
		if c, ok := cld.(*goxml.Element); ok {
			if fontface, err = xd.getAttributeString("fontface", c, true, true, ""); err != nil {
				return nil, err
			}

			switch c.Name {
			case "Regular":
				ff.AddMember(xd.fontsources[fontface], document.FontWeight400, document.FontStyleNormal)
			case "Italic":
				ff.AddMember(xd.fontsources[fontface], document.FontWeight400, document.FontStyleItalic)
			case "Bold":
				ff.AddMember(xd.fontsources[fontface], document.FontWeight700, document.FontStyleNormal)
			case "BoldItalic":
				ff.AddMember(xd.fontsources[fontface], document.FontWeight700, document.FontStyleItalic)
			}
		}
	}
	return nil, nil
}

func cmdDefineFontsize(xd *xtsDocument, layoutelt *goxml.Element, dataelt *xpath.Parser) (xpath.Sequence, error) {
	var name string
	var fontsize, leading bag.ScaledPoint
	var err error
	if name, err = xd.getAttributeString("name", layoutelt, true, true, ""); err != nil {
		return nil, err
	}
	if fontsize, err = xd.getAttributeSize("fontsize", layoutelt, true, true, ""); err != nil {
		return nil, err
	}
	if leading, err = xd.getAttributeSize("leading", layoutelt, true, true, ""); err != nil {
		return nil, err
	}

	if xd.fontsizes == nil {
		xd.fontsizes = make(map[string][2]bag.ScaledPoint)
	}
	xd.fontsizes[name] = [2]bag.ScaledPoint{fontsize, leading}
	return nil, nil
}

func cmdImage(xd *xtsDocument, layoutelt *goxml.Element, dataelt *xpath.Parser) (xpath.Sequence, error) {
	var filename string
	var err error
	if filename, err = xd.getAttributeString("file", layoutelt, true, true, ""); err != nil {
		return nil, err
	}

	filename, err = xd.cfg.FindFile(filename)
	if err != nil {
		return nil, err
	}
	var imgObj *pdf.Imagefile
	if imgfile, ok := loadedImages[filename]; !ok {

		imgObj, err = xd.doc.LoadImageFile(filename)
		if err != nil {
			return nil, err
		}
	} else {
		imgObj = imgfile
	}
	hl := createImageHlist(xd, imgObj)
	if hl.Attibutes == nil {
		hl.Attibutes = node.H{}
	}
	hl.Attibutes["origin"] = "image"

	return xpath.Sequence{hl}, nil

}

func cmdLoadFontfile(xd *xtsDocument, layoutelt *goxml.Element, dataelt *xpath.Parser) (xpath.Sequence, error) {
	var filename, name string
	var err error
	if filename, err = xd.getAttributeString("filename", layoutelt, true, true, ""); err != nil {
		return nil, err
	}
	if name, err = xd.getAttributeString("name", layoutelt, true, true, ""); err != nil {
		return nil, err
	}
	fn, err := xd.cfg.FindFile(filename)
	if err != nil {
		return nil, err
	}
	fs := document.FontSource{
		Name:   name,
		Source: fn,
	}
	// Not necessary when default fonts are initialized
	if xd.fontsources == nil {
		xd.fontsources = make(map[string]*document.FontSource)
	}
	xd.fontsources[name] = &fs
	return nil, nil
}

func cmdRecord(xd *xtsDocument, layoutelt *goxml.Element, dataelt *xpath.Parser) (xpath.Sequence, error) {
	var elt, mode string
	var err error
	if elt, err = xd.getAttributeString("element", layoutelt, true, false, ""); err != nil {
		return nil, err
	}
	if mode, err = xd.getAttributeString("mode", layoutelt, false, false, ""); err != nil {
		return nil, err
	}
	dp := dataDispatcher[elt]
	if dp == nil {
		dataDispatcher[elt] = make(map[string]*goxml.Element)
	}
	dataDispatcher[elt][mode] = layoutelt
	return nil, nil
}

func cmdOptions(xd *xtsDocument, layoutelt *goxml.Element, dataelt *xpath.Parser) (xpath.Sequence, error) {
	mainlanguageString, err := xd.getAttributeString("mainlanguage", layoutelt, false, true, "")
	if err != nil {
		return nil, err
	}
	l, err := xd.getLanguage(mainlanguageString)
	if err != nil {
		return nil, err
	}
	bag.Logger.Infof("Setting default language to %q", l.Name)
	xd.doc.DefaultLanguage = l

	return xpath.Sequence{}, nil
}

func cmdParagraph(xd *xtsDocument, layoutelt *goxml.Element, dataelt *xpath.Parser) (xpath.Sequence, error) {
	colorString, err := xd.getAttributeString("color", layoutelt, false, true, "")
	if err != nil {
		return nil, err
	}

	seq, err := dispatch(xd, layoutelt, dataelt)
	if err != nil {
		return nil, err
	}

	te := &document.TypesettingElement{
		Settings: make(document.TypesettingSettings),
	}
	if colorString != "" {
		te.Settings[document.SettingColor] = colorString
	}
	getTextvalues(te, seq, "cmdParagraph")
	return xpath.Sequence{te}, nil
}

func cmdPlaceObject(xd *xtsDocument, layoutelt *goxml.Element, dataelt *xpath.Parser) (xpath.Sequence, error) {
	xd.setupPage()
	columnString, err := xd.getAttributeString("column", layoutelt, false, true, "")
	if err != nil {
		return nil, err
	}
	rowString, err := xd.getAttributeString("row", layoutelt, false, true, "")
	if err != nil {
		return nil, err
	}

	pos := positioningUnknown

	var columnGrid, rowGrid int
	var columnLength, rowLength bag.ScaledPoint
	if columnGrid, err = strconv.Atoi(columnString); err == nil {
		pos = positioningGrid
	}
	if rowGrid, err = strconv.Atoi(rowString); err == nil {
		if pos != positioningGrid {
			return nil, fmt.Errorf("both column and row must be integers with grid positioning")
		}
	}

	if pos == positioningUnknown {
		if columnLength, err = bag.Sp(columnString); err == nil {
			pos = positioningAbsolute
		}
		if rowLength, err = bag.Sp(rowString); err == nil {
			pos = positioningAbsolute
		}
	}

	seq, err := dispatch(xd, layoutelt, dataelt)
	if err != nil {
		return nil, err
	}
	var origin string
	var vl *node.VList
	switch t := seq[0].(type) {
	case *node.VList:
		vl = t
		if vl.Attibutes != nil {
			origin = vl.Attibutes["origin"].(string)
		}
	case *node.HList:
		vl = node.NewVList()
		vl.List = t
		if t.Attibutes != nil {
			origin = t.Attibutes["origin"].(string)
		}
	default:
		bag.Logger.DPanicf("PlaceObject: unknown node %v", t)
	}

	switch pos {
	case positioningAbsolute:
		xd.currentPage.outputAbsolute(columnLength, rowLength, vl)
	case positioningGrid:
		bag.Logger.Infof("PlaceObject: output %s at (%d,%d)", origin, rowGrid, columnGrid)
		columnLength = xd.currentGrid.posX(columnGrid)
		rowLength = xd.currentGrid.posY(rowGrid)
		xd.currentPage.outputAbsolute(columnLength, rowLength, vl)
	}

	return seq, nil
}

func cmdSetGrid(xd *xtsDocument, layoutelt *goxml.Element, dataelt *xpath.Parser) (xpath.Sequence, error) {
	var err error
	var height, width, dx, dy bag.ScaledPoint
	var nx, ny int
	if height, err = xd.getAttributeSize("height", layoutelt, false, true, ""); err != nil {
		return nil, err
	}
	if width, err = xd.getAttributeSize("width", layoutelt, false, true, ""); err != nil {
		return nil, err
	}
	if dx, err = xd.getAttributeSize("dx", layoutelt, false, true, ""); err != nil {
		return nil, err
	}
	if dy, err = xd.getAttributeSize("dy", layoutelt, false, true, ""); err != nil {
		return nil, err
	}
	if nx, err = xd.getAttributeInt("nx", layoutelt, false, true, ""); err != nil {
		return nil, err
	}
	if ny, err = xd.getAttributeInt("ny", layoutelt, false, true, ""); err != nil {
		return nil, err
	}

	if height > 0 {
		xd.defaultGridHeight = height
	}
	if width > 0 {
		xd.defaultGridWidth = width
	}
	if dx > 0 {
		xd.defaultGridGapX = dx
	}
	if dy > 0 {
		xd.defaultGridGapY = dy
	}
	if nx > 0 {
		xd.defaultGridNx = nx
	}
	if ny > 0 {
		xd.defaultGridNy = ny
	}
	return nil, nil
}

func cmdTextblock(xd *xtsDocument, layoutelt *goxml.Element, dataelt *xpath.Parser) (xpath.Sequence, error) {
	attrFontsize, err := xd.getAttributeString("fontsize", layoutelt, false, true, "")
	if err != nil {
		return nil, err
	}

	attrWidth, err := xd.getAttributeWidth("width", layoutelt, false, true, "")
	if err != nil {
		return nil, err
	}

	leading := 12 * bag.Factor
	fontsize := 10 * bag.Factor

	if sp := strings.Split(attrFontsize, "/"); len(sp) == 2 {
		if fontsize, err = bag.Sp(sp[0]); err != nil {
			return nil, err
		}
		if leading, err = bag.Sp(sp[1]); err != nil {
			return nil, err
		}
	} else if fs, ok := xd.fontsizes[attrFontsize]; ok {
		fontsize = fs[0]
		leading = fs[1]
	} else if attrFontsize == "" {
		// ok, ignore
		bag.Logger.Debug("use default font size text")
	} else {
		return nil, fmt.Errorf("unknown font size %s", attrFontsize)
	}

	seq, err := dispatch(xd, layoutelt, dataelt)
	if err != nil {
		return nil, err
	}

	te := &document.TypesettingElement{
		Settings: document.TypesettingSettings{
			document.SettingFontFamily: xd.doc.GetFontFamily(0),
			document.SettingSize:       fontsize,
		},
	}

	for _, itm := range seq {
		switch t := itm.(type) {
		case *document.TypesettingElement:
			te.Items = append(te.Items, t)
		default:
			bag.Logger.DPanicf("cmdTextblock: unknown type %T", t)

		}
	}
	hlist, tail, err := xd.doc.Mknodes(te)
	if err != nil {
		return nil, err
	}
	xd.doc.Hyphenate(hlist)
	node.AppendLineEndAfter(tail)

	ls := node.NewLinebreakSettings()
	ls.HSize = attrWidth
	ls.LineHeight = leading
	vlist, _ := node.Linebreak(hlist, ls)

	if vlist.Attibutes == nil {
		vlist.Attibutes = node.H{}
	}
	vlist.Attibutes["origin"] = "textblock"

	return xpath.Sequence{vlist}, nil
}

func cmdTrace(xd *xtsDocument, layoutelt *goxml.Element, dataelt *xpath.Parser) (xpath.Sequence, error) {
	if traceGrid, err := xd.getAttributeBool("grid", layoutelt, false, false, ""); err == nil && traceGrid {
		xd.SetVTrace(VTraceGrid)
	}
	if traceHyphenation, err := xd.getAttributeBool("hyphenation", layoutelt, false, false, ""); err == nil && traceHyphenation {
		xd.SetVTrace(VTraceHyphenation)
	}

	return nil, nil
}

func cmdValue(xd *xtsDocument, layoutelt *goxml.Element, dataelt *xpath.Parser) (xpath.Sequence, error) {
	var selection string
	var err error

	if selection, err = xd.getAttributeString("select", layoutelt, false, false, ""); err != nil {
		return nil, err
	}
	if selection != "" {
		eval, err := dataelt.Evaluate(selection)
		if err != nil {
			return nil, err
		}
		return eval, nil
	}
	seq := xpath.Sequence{}
	for _, cld := range layoutelt.Children() {
		seq = append(seq, cld)
	}

	return seq, nil
}
