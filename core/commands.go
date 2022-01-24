package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/node"
	"github.com/speedata/boxesandglue/csshtml"
	"github.com/speedata/boxesandglue/document"
	"github.com/speedata/boxesandglue/pdfbackend/pdf"
	"github.com/speedata/goxml"
	"github.com/speedata/goxpath/xpath"
)

type commandFunc func(*xtsDocument, *goxml.Element) (xpath.Sequence, error)

var (
	dataDispatcher = make(map[string]map[string]*goxml.Element)
	dispatchTable  map[string]commandFunc
)

func init() {
	dispatchTable = map[string]commandFunc{
		"A":                cmdA,
		"B":                cmdB,
		"Box":              cmdBox,
		"Color":            cmdColor,
		"DefineFontfamily": cmdDefineFontfamily,
		"DefineFontsize":   cmdDefineFontsize,
		"DefinePagetype":   cmdDefinePagetype,
		"Image":            cmdImage,
		"LoadFontfile":     cmdLoadFontfile,
		"Options":          cmdOptions,
		"Pageformat":       cmdPageformat,
		"Paragraph":        cmdParagraph,
		"PlaceObject":      cmdPlaceObject,
		"Record":           cmdRecord,
		"SetGrid":          cmdSetGrid,
		"Stylesheet":       cmdStylesheet,
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
				seq, err := f(xd, elt)
				if err != nil {
					return nil, err
				}
				retSequence = append(retSequence, seq...)
			} else {
				bag.Logger.Errorf("layout: element %q unknown", elt.Name)
				return nil, fmt.Errorf("layout: element %q unknown", elt.Name)
			}
		}
	}
	return retSequence, nil
}

func cmdA(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Href string
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	seq, err := dispatch(xd, layoutelt, xd.data)
	if err != nil {
		return nil, err
	}
	hl := document.Hyperlink{URI: attValues.Href}

	te := &document.TypesettingElement{
		Settings: document.TypesettingSettings{
			document.SettingHyperlink: hl,
		},
	}
	getTextvalues(te, seq, "cmdA")

	return xpath.Sequence{te}, err
}

func cmdB(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	seq, err := dispatch(xd, layoutelt, xd.data)

	te := &document.TypesettingElement{
		Settings: document.TypesettingSettings{
			document.SettingFontWeight: 700,
		},
	}
	getTextvalues(te, seq, "cmdBold")

	return xpath.Sequence{te}, err
}

func cmdBox(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	attValues := &struct {
		Class           string
		ID              string
		Backgroundcolor string          `sdxml:"default:black"`
		Width           bag.ScaledPoint `sdxml:"mustexist"`
		Height          bag.ScaledPoint `sdxml:"mustexist"`
	}{}

	if err := getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	htmlString := `<box `
	if class := attValues.Class; class != "" {
		htmlString += fmt.Sprintf("class=%q ", class)
	}
	if id := attValues.ID; id != "" {
		htmlString += fmt.Sprintf("id=%q", id)
	}
	htmlString += ">"

	a, err := xd.layoutcss.ParseHTMLFragment(htmlString, "")
	if err != nil {
		return nil, err
	}

	attrs, _ := csshtml.ResolveAttributes(a.Find("box").First().Nodes[0].Attr)
	var bgcolor *document.Color
	if bgc, ok := attrs["background-color"]; ok {
		bgcolor = xd.doc.GetColor(bgc)
	} else {
		bgcolor = xd.doc.GetColor(attValues.Backgroundcolor)
	}

	r := node.NewRule()
	if bgcolor.Space != document.ColorNone {
		r.Pre = fmt.Sprintf(" %s 0 0 %s %s re f ", bgcolor.PDFStringFG(), attValues.Width, -attValues.Height)
	}

	r.Width = attValues.Width
	r.Height = attValues.Height
	vl := node.Vpack(r)
	return xpath.Sequence{vl}, err
}

func cmdColor(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Name string
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	seq, err := dispatch(xd, layoutelt, xd.data)

	te := &document.TypesettingElement{
		Settings: document.TypesettingSettings{
			document.SettingColor: attValues.Name,
		},
	}
	getTextvalues(te, seq, "cmdColor")

	return xpath.Sequence{te}, err
}

func cmdDefineFontfamily(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Name string
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	ff := xd.doc.NewFontFamily(attValues.Name)
	var fontface string
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

func cmdDefineFontsize(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	attValues := &struct {
		Name     string
		Fontsize bag.ScaledPoint
		Leading  bag.ScaledPoint
	}{}

	if err := getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	if xd.fontsizes == nil {
		xd.fontsizes = make(map[string][2]bag.ScaledPoint)
	}
	xd.fontsizes[attValues.Name] = [2]bag.ScaledPoint{attValues.Fontsize, attValues.Leading}
	return nil, nil
}

func cmdDefinePagetype(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Margin string
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	pt, err := xd.newPagetype("mypage", "true()")
	if err != nil {
		return nil, err
	}
	fv, err := getFourValuesSP(attValues.Margin)
	if err != nil {
		return nil, err
	}
	pt.marginBottom = fv["bottom"]
	pt.marginLeft = fv["left"]
	pt.marginRight = fv["right"]
	pt.marginTop = fv["top"]
	return xpath.Sequence{}, nil
}

func cmdImage(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Href string
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	filename, err := xd.cfg.FindFile(attValues.Href)
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

func cmdLoadFontfile(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Filename string `sdxml:"mustexist"`
		Name     string `sdxml:"mustexist"`
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	fn, err := xd.cfg.FindFile(attValues.Filename)
	if err != nil {
		return nil, err
	}
	fs := document.FontSource{
		Name:   attValues.Name,
		Source: fn,
	}
	// Not necessary when default fonts are initialized
	if xd.fontsources == nil {
		xd.fontsources = make(map[string]*document.FontSource)
	}
	xd.fontsources[attValues.Name] = &fs
	return nil, nil
}

func cmdRecord(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Element string `sdxml:"mustexist"`
		Mode    string
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	dp := dataDispatcher[attValues.Element]
	if dp == nil {
		dataDispatcher[attValues.Element] = make(map[string]*goxml.Element)
	}
	dataDispatcher[attValues.Element][attValues.Mode] = layoutelt
	return nil, nil
}

func cmdOptions(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Mainlanguage *string
		Bleed        *bag.ScaledPoint
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	if attValues.Mainlanguage != nil {
		l, err := xd.getLanguage(*attValues.Mainlanguage)
		if err != nil {
			return nil, err
		}
		bag.Logger.Infof("Setting default language to %q", l.Name)
		xd.doc.DefaultLanguage = l
	}

	if attValues.Bleed != nil {
		xd.doc.Bleed = *attValues.Bleed
	}

	return xpath.Sequence{}, nil
}

func cmdPageformat(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Width  bag.ScaledPoint `sdxml:"mustexist"`
		Height bag.ScaledPoint `sdxml:"mustexist"`
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	xd.doc.DefaultPageWidth = attValues.Width
	xd.doc.DefaultPageHeight = attValues.Height
	return xpath.Sequence{}, nil
}

func cmdParagraph(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Color string
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	seq, err := dispatch(xd, layoutelt, xd.data)
	if err != nil {
		return nil, err
	}

	te := &document.TypesettingElement{
		Settings: make(document.TypesettingSettings),
	}
	if attValues.Color != "" {
		te.Settings[document.SettingColor] = attValues.Color
	}
	getTextvalues(te, seq, "cmdParagraph")
	return xpath.Sequence{te}, nil
}

func cmdPlaceObject(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	xd.setupPage()
	var err error
	attValues := &struct {
		Column string
		Row    string
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	pos := positioningUnknown

	var columnInt, rowInt int
	var col, row coord
	var columnLength, rowLength bag.ScaledPoint
	if columnInt, err = strconv.Atoi(attValues.Column); err == nil {
		pos = positioningGrid
		col = coord(columnInt)
	}

	if rowInt, err = strconv.Atoi(attValues.Row); err == nil {
		if pos != positioningGrid {
			return nil, fmt.Errorf("both column and row must be integers with grid positioning")
		}
		row = coord(rowInt)
	}

	if pos == positioningUnknown {
		if columnLength, err = bag.Sp(attValues.Column); err == nil {
			pos = positioningAbsolute
		}
		if rowLength, err = bag.Sp(attValues.Row); err == nil {
			pos = positioningAbsolute
		}
	}

	seq, err := dispatch(xd, layoutelt, xd.data)
	if err != nil {
		return nil, err
	}
	if len(seq) == 0 {
		bag.Logger.Warnf("line %d: no objects in PlaceObject", layoutelt.Line)
		return nil, nil
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
		vl = node.Vpack(t)
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
		bag.Logger.Infof("PlaceObject: output %s at (%d,%d)", origin, rowInt, columnInt)
		columnLength = xd.currentGrid.posX(col)
		rowLength = xd.currentGrid.posY(row)
		xd.currentPage.outputAbsolute(columnLength, rowLength, vl)
		xd.currentGrid.allocate(col, row, vl.Width, vl.Height)
	}

	return seq, nil
}

func cmdSetGrid(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	attValues := &struct {
		Nx     int
		Ny     int
		Dx     bag.ScaledPoint
		Dy     bag.ScaledPoint
		Width  bag.ScaledPoint
		Height bag.ScaledPoint
	}{}
	if err := getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	if height := attValues.Height; height > 0 {
		xd.defaultGridHeight = height
	}
	if width := attValues.Width; width > 0 {
		xd.defaultGridWidth = width
	}
	if dx := attValues.Dx; dx > 0 {
		xd.defaultGridGapX = dx
	}
	if dy := attValues.Dy; dy > 0 {
		xd.defaultGridGapY = dy
	}
	if nx := attValues.Nx; nx > 0 {
		xd.defaultGridNx = nx
	}
	if ny := attValues.Ny; ny > 0 {
		xd.defaultGridNy = ny
	}
	return nil, nil
}

func cmdStylesheet(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Scope string
		Href  string
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	var toks csshtml.Tokenstream
	if attrHref := attValues.Href; attrHref == "" {
		toks, err = xd.layoutcss.ParseCSSString(layoutelt.Stringvalue())
	} else {
		toks, err = xd.layoutcss.ParseCSSFile(attrHref)

	}
	if err != nil {
		bag.Logger.Error(err)
		return nil, nil
	}
	parsedStyles := csshtml.ConsumeBlock(toks, false)
	switch attValues.Scope {
	case "layout":
		xd.layoutcss.Stylesheet = append(xd.layoutcss.Stylesheet, parsedStyles)
	case "data":
		bag.Logger.Errorf("not implemented yet: scope=%q in Stylesheet (line %d)", attValues.Scope, layoutelt.Line)
	default:
		bag.Logger.Errorf("unknown scope: %q in Stylesheet (line %d)", attValues.Scope, layoutelt.Line)
	}

	return xpath.Sequence{nil}, nil
}

func cmdTextblock(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Fontsize string
		Width    bag.ScaledPoint
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	leading := 12 * bag.Factor
	fontsize := 10 * bag.Factor
	attrFontsize := attValues.Fontsize

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

	seq, err := dispatch(xd, layoutelt, xd.data)
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
	ls.HSize = attValues.Width
	ls.LineHeight = leading
	vlist, _ := node.Linebreak(hlist, ls)

	if vlist.Attibutes == nil {
		vlist.Attibutes = node.H{}
	}
	vlist.Attibutes["origin"] = "textblock"

	return xpath.Sequence{vlist}, nil
}

func cmdTrace(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Grid           *bool
		Hyphenation    *bool
		Gridallocation *bool
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	if attValues.Grid != nil {
		if *attValues.Grid {
			xd.SetVTrace(VTraceGrid)
		} else {
			xd.ClearVTrace(VTraceGrid)
		}
	}

	if attValues.Hyphenation != nil {
		if *attValues.Hyphenation {
			xd.SetVTrace(VTraceHyphenation)
		} else {
			xd.ClearVTrace(VTraceHyphenation)
		}
	}

	if attValues.Gridallocation != nil {
		if *attValues.Gridallocation {
			xd.SetVTrace(VTraceAllocation)
		} else {
			xd.ClearVTrace(VTraceAllocation)
		}
	}

	return nil, nil
}

func cmdValue(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Select *string
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	if attValues.Select != nil {
		eval, _ := xd.data.Evaluate(*attValues.Select)
		// if err != nil {
		// 	return nil, err
		// }
		fmt.Println(eval)
		return eval, nil
	}
	seq := xpath.Sequence{}
	for _, cld := range layoutelt.Children() {
		seq = append(seq, cld)
	}

	return seq, nil
}
