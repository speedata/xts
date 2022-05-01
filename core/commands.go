package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/document"
	"github.com/speedata/boxesandglue/backend/node"
	"github.com/speedata/boxesandglue/csshtml"
	"github.com/speedata/boxesandglue/frontend"
	"github.com/speedata/boxesandglue/pdfbackend/pdf"
	"github.com/speedata/goxml"
	xpath "github.com/speedata/goxpath"
)

type commandFunc func(*xtsDocument, *goxml.Element) (xpath.Sequence, error)

var (
	dataDispatcher = make(map[string]map[string]*goxml.Element)
	dispatchTable  map[string]commandFunc
)

func init() {
	dispatchTable = map[string]commandFunc{
		"A":                cmdA,
		"Attribute":        cmdAttribute,
		"B":                cmdB,
		"Box":              cmdBox,
		"ClearPage":        cmdClearpage,
		"Color":            cmdColor,
		"Contents":         cmdContents,
		"Copy-of":          cmdCopyof,
		"DefineFontfamily": cmdDefineFontfamily,
		"DefineFontsize":   cmdDefineFontsize,
		"DefinePagetype":   cmdDefinePagetype,
		"Element":          cmdElement,
		"ForAll":           cmdForall,
		"Group":            cmdGroup,
		"Image":            cmdImage,
		"LoadFontfile":     cmdLoadFontfile,
		"Message":          cmdMessage,
		"NextFrame":        cmdNextFrame,
		"Options":          cmdOptions,
		"Pageformat":       cmdPageformat,
		"Paragraph":        cmdParagraph,
		"PlaceObject":      cmdPlaceObject,
		"ProcessNode":      cmdProcessNode,
		"Record":           cmdRecord,
		"SetGrid":          cmdSetGrid,
		"SetVariable":      cmdSetVariable,
		"Stylesheet":       cmdStylesheet,
		"Switch":           cmdSwitch,
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

	te := &frontend.TypesettingElement{
		Settings: frontend.TypesettingSettings{
			frontend.SettingHyperlink: hl,
		},
	}
	getTextvalues(te, seq, "cmdA", layoutelt.Line)

	return xpath.Sequence{te}, err
}

func cmdAttribute(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	attValues := &struct {
		Select string `sdxml:"noescape"`
		Name   string `sdxml:"mustexist"`
	}{}
	if err := getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	var eval xpath.Sequence
	var err error
	eval, err = xd.data.Evaluate(attValues.Select)
	if err != nil {
		bag.Logger.Errorf("Attribute (line %d): error parsing select XPath expression %s", layoutelt.Line, err)
		return nil, err
	}

	attr := goxml.Attribute{
		Name:  attValues.Name,
		Value: eval.Stringvalue(),
	}
	return xpath.Sequence{attr}, nil
}

func cmdB(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	seq, err := dispatch(xd, layoutelt, xd.data)

	te := &frontend.TypesettingElement{
		Settings: frontend.TypesettingSettings{
			frontend.SettingFontWeight: 700,
		},
	}
	getTextvalues(te, seq, "cmdBold", layoutelt.Line)

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
	var bgcolor *frontend.Color
	if bgc, ok := attrs["background-color"]; ok {
		bgcolor = xd.document.GetColor(bgc)
	} else {
		bgcolor = xd.document.GetColor(attValues.Backgroundcolor)
	}

	r := node.NewRule()
	if bgcolor.Space != frontend.ColorNone {
		r.Pre = fmt.Sprintf(" %s 0 0 %s %s re f ", bgcolor.PDFStringFG(), attValues.Width, -attValues.Height)
	}

	r.Width = attValues.Width
	r.Height = attValues.Height
	vl := node.Vpack(r)
	return xpath.Sequence{vl}, err
}

func cmdClearpage(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	clearPage(xd)
	return nil, nil
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

	te := &frontend.TypesettingElement{
		Settings: frontend.TypesettingSettings{
			frontend.SettingColor: attValues.Name,
		},
	}
	getTextvalues(te, seq, "cmdColor", layoutelt.Line)

	return xpath.Sequence{te}, err
}

func cmdContents(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	_, err := dispatch(xd, layoutelt, xd.data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func cmdCopyof(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Select string `sdxml:"mustexist"`
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	var eval xpath.Sequence
	eval, err = xd.data.Evaluate(attValues.Select)

	return eval, nil
}

func cmdDefineFontfamily(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Name string
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	ff := xd.document.NewFontFamily(attValues.Name)
	var fontface string
	for _, cld := range layoutelt.Children() {
		if c, ok := cld.(*goxml.Element); ok {
			if fontface, err = xd.getAttributeString("fontface", c, true, true, ""); err != nil {
				return nil, err
			}

			switch c.Name {
			case "Regular":
				ff.AddMember(xd.fontsources[fontface], frontend.FontWeight400, frontend.FontStyleNormal)
			case "Italic":
				ff.AddMember(xd.fontsources[fontface], frontend.FontWeight400, frontend.FontStyleItalic)
			case "Bold":
				ff.AddMember(xd.fontsources[fontface], frontend.FontWeight700, frontend.FontStyleNormal)
			case "BoldItalic":
				ff.AddMember(xd.fontsources[fontface], frontend.FontWeight700, frontend.FontStyleItalic)
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
		Margin string `sdxml:"mustexist"`
		Name   string `sdxml:"mustexist"`
		Test   string `sdxml:"mustexist"`
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	pt, err := xd.newPagetype(attValues.Name, attValues.Test)
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

	pt.layoutElt = layoutelt
	return xpath.Sequence{}, nil
}

func cmdElement(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Name string `sdxml:"mustexist"`
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	elt := goxml.Element{}
	elt.Name = attValues.Name

	seq, err := dispatch(xd, layoutelt, xd.data)
	if err != nil {
		return nil, err
	}
	for _, itm := range seq {
		switch t := itm.(type) {
		case goxml.XMLNode:
			elt.Append(goxml.XMLNode(t))
		default:
			bag.Logger.DPanicf("Element (line %d): don't know how to append %T", layoutelt.Line, t)
		}

	}
	return xpath.Sequence{elt}, nil
}

func cmdForall(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Select string `sdxml:"noescape"`
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	var eval xpath.Sequence
	eval, err = xd.data.Evaluate(attValues.Select)
	if err != nil {
		bag.Logger.Errorf("ForAll (line %d): error parsing select XPath expression %s", layoutelt.Line, err)
		return nil, err
	}

	oldContext := xd.data.Ctx.SetContext(xpath.Sequence{})
	for i, itm := range eval {
		xd.data.Ctx.SetContext(xpath.Sequence{itm})
		xd.data.Ctx.Pos = i + 1
		eval, err = dispatch(xd, layoutelt, xd.data)
		if err != nil {
			return nil, err
		}
	}
	xd.data.Ctx.SetContext(xpath.Sequence{oldContext})

	return nil, nil
}

func cmdGroup(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Name string `sdxml:"mustexist"`
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	saveGrid := xd.currentGrid
	xd.currentGroup = xd.newGroup(attValues.Name)
	xd.currentGrid = xd.currentGroup.grid
	_, err = dispatch(xd, layoutelt, xd.data)
	if err != nil {
		return nil, err
	}
	xd.currentGroup = nil
	xd.currentGrid = saveGrid
	return nil, nil
}

func cmdImage(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Href      string `sdxml:"mustexist"`
		Height    *bag.ScaledPoint
		Width     *bag.ScaledPoint
		MinHeight *bag.ScaledPoint
		MinWidth  *bag.ScaledPoint
		MaxHeight *bag.ScaledPoint
		MaxWidth  *bag.ScaledPoint
		Stretch   bool
		Page      int
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	filename, err := xd.cfg.FindFile(attValues.Href)
	if err != nil {
		return nil, err
	}
	var imgObj *pdf.Imagefile
	imgObj, err = xd.document.Doc.LoadImageFile(filename)
	if err != nil {
		return nil, err
	}

	hl := createImageHlist(
		xd,
		attValues.Width,
		attValues.Height,
		attValues.MinWidth,
		attValues.MaxWidth,
		attValues.MinHeight,
		attValues.MaxHeight,
		attValues.Stretch,
		imgObj,
		attValues.Page)
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
	fs := frontend.FontSource{
		Name:   attValues.Name,
		Source: fn,
	}
	// Not necessary when default fonts are initialized
	if xd.fontsources == nil {
		xd.fontsources = make(map[string]*frontend.FontSource)
	}
	xd.fontsources[attValues.Name] = &fs
	return nil, nil
}

func cmdProcessNode(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Select string `sdxml:"mustexist"`
		Mode   string
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	var eval xpath.Sequence
	eval, err = evaluateXPath(xd, layoutelt, attValues.Select)
	if err != nil {
		bag.Logger.Errorf("ProcessNode (line %d): error parsing select XPath expression %s", layoutelt.Line, err)
		return nil, err
	}

	oldContext := xd.data.Ctx.SetContext(xpath.Sequence{})

	for i, itm := range eval {
		xd.data.Ctx.Pos = i + 1
		if elt, ok := itm.(*goxml.Element); ok {
			if dd, ok := dataDispatcher[elt.Name]; ok {
				if rec, ok := dd[attValues.Mode]; ok {
					_, err = dispatch(xd, rec, xd.data)
				}
			}
		}
	}

	xd.data.Ctx.SetContext(xpath.Sequence{oldContext})
	return nil, nil
}

func cmdRecord(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Match string `sdxml:"mustexist"`
		Mode  string
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	dp := dataDispatcher[attValues.Match]
	if dp == nil {
		dataDispatcher[attValues.Match] = make(map[string]*goxml.Element)
	}
	dataDispatcher[attValues.Match][attValues.Mode] = layoutelt
	return nil, nil
}

func cmdMessage(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Select *string `sdxml:"noescape"`
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	var eval xpath.Sequence
	if attValues.Select != nil {
		eval, err = evaluateXPath(xd, layoutelt, *attValues.Select)
		if err != nil {
			bag.Logger.Errorf("Message (line %d): error parsing select XPath expression %s", layoutelt.Line, err)
			return nil, err
		}
	} else {
		eval, err = dispatch(xd, layoutelt, xd.data)
		if err != nil {
			return nil, err
		}
	}
	bag.Logger.Infof("Message (line %d): %s", layoutelt.Line, eval)
	return nil, nil
}

func cmdNextFrame(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Area string `sdxml:"mustexist"`
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	xd.setupPage()
	if area, ok := xd.currentPage.pagegrid.areas[attValues.Area]; ok {
		area.currentFrame++
		if area.currentFrame == len(area.frame) {
			area.currentFrame = 0
			clearPage(xd)
		}
		area.currentRow = 0
	}
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
		xd.document.Doc.DefaultLanguage = l
	}

	if attValues.Bleed != nil {
		xd.document.Doc.Bleed = *attValues.Bleed
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

	xd.document.Doc.DefaultPageWidth = attValues.Width
	xd.document.Doc.DefaultPageHeight = attValues.Height
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

	te := &frontend.TypesettingElement{
		Settings: make(frontend.TypesettingSettings),
	}
	if attValues.Color != "" {
		te.Settings[frontend.SettingColor] = attValues.Color
	}
	getTextvalues(te, seq, "cmdParagraph", layoutelt.Line)
	return xpath.Sequence{te}, nil
}

func cmdPlaceObject(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	xd.setupPage()
	var err error
	attValues := &struct {
		Column    string
		Row       string
		Area      string
		Groupname string
	}{}
	if err = getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	if attValues.Area == "" {
		attValues.Area = defaultAreaName
	}
	pos := positioningUnknown

	var columnInt, rowInt int
	var col, row coord
	var columnLength, rowLength bag.ScaledPoint

	var rowSet, colSet bool
	if columnInt, err = strconv.Atoi(attValues.Column); err == nil {
		colSet = true
		pos = positioningGrid
		col = coord(columnInt)
	}

	if rowInt, err = strconv.Atoi(attValues.Row); err == nil {
		rowSet = true
		pos = positioningGrid
		row = coord(rowInt)
	}

	if pos == positioningGrid && colSet != rowSet {
		bag.Logger.Errorf("line %d: both column and row must be integers with grid positioning", layoutelt.Line)
		if !colSet {
			col, columnInt = 1, 1
		}
		if !rowSet {
			row, rowInt = 1, 1
		}
	} else if pos == positioningUnknown {
		if columnLength, err = bag.Sp(attValues.Column); err == nil {
			pos = positioningAbsolute
		}
		if rowLength, err = bag.Sp(attValues.Row); err == nil {
			pos = positioningAbsolute
		}
	}
	var seq xpath.Sequence
	if attValues.Groupname != "" {
		seq = xpath.Sequence{xd.groups[attValues.Groupname].contents}
	} else {
		seq, err = dispatch(xd, layoutelt, xd.data)
		if err != nil {
			return nil, err
		}

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

	if pos == positioningUnknown {
		pos = positioningGrid
		xy, err := xd.currentGrid.findFreeSpaceForObject(vl, attValues.Area)
		if err != nil {
			return nil, err
		}
		bag.Logger.Debugf("looking for free space for %s", origin)
		col, row = xy.XY()
		columnInt = int(col)
		rowInt = int(row)
	}

	switch pos {
	case positioningAbsolute:
		xd.currentPage.outputAbsolute(columnLength, rowLength, vl)
	case positioningGrid:
		xd.OutputAt(vl, col, row, true, attValues.Area, origin)
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

func cmdSetVariable(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	attValues := &struct {
		Select   *string `sdxml:"noescape"`
		Variable string  `sdxml:"mustexist"`
		Trace    bool
	}{}
	if err := getXMLAtttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	var eval xpath.Sequence
	var err error
	if attValues.Select != nil {
		eval, err = xd.data.Evaluate(*attValues.Select)
		if err != nil {
			bag.Logger.Errorf("SetVariable (line %d): error parsing select XPath expression %s", layoutelt.Line, err)
			return nil, err
		}
		xd.data.SetVariable(attValues.Variable, eval)
	} else {
		eval, err = dispatch(xd, layoutelt, xd.data)
		if err != nil {
			return nil, err
		}
	}
	xd.data.SetVariable(attValues.Variable, eval)
	if attValues.Trace {
		bag.Logger.Infof("SetVariable (line %d): %s to %s", layoutelt.Line, attValues.Variable, eval)
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

func cmdSwitch(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error

	for _, cld := range layoutelt.Children() {
		if c, ok := cld.(*goxml.Element); ok {
			if c.Name == "Case" {
				attrs := c.Attributes()
				for _, attr := range attrs {
					if attr.Name == "test" {
						var eval xpath.Sequence
						eval, err = xd.data.Evaluate(attr.Value)
						if err != nil {
							return nil, err
						}
						var ok bool
						if ok, err = xpath.BooleanValue(eval); err != nil {
							return nil, err
						}
						if ok {
							return dispatch(xd, c, xd.data)
						}

					}
				}
			} else if c.Name == "Otherwise" {
				return dispatch(xd, c, xd.data)
			}
		}
	}
	return nil, nil
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

	te := &frontend.TypesettingElement{
		Settings: frontend.TypesettingSettings{
			frontend.SettingFontFamily: xd.document.GetFontFamily(0),
			frontend.SettingSize:       fontsize,
		},
	}

	for _, itm := range seq {
		switch t := itm.(type) {
		case *frontend.TypesettingElement:
			te.Items = append(te.Items, t)
		default:
			bag.Logger.DPanicf("cmdTextblock: unknown type %T", t)
		}
	}
	hlist, tail, err := xd.document.Mknodes(te)
	if err != nil {
		return nil, err
	}
	frontend.Hyphenate(hlist, xd.defaultLanguage)
	node.AppendLineEndAfter(tail)

	ls := node.NewLinebreakSettings()

	if attValues.Width == 0 {
		if xd.currentGrid.currentCol > coord(xd.currentGrid.nx) {
			xd.currentGrid.nextRow()
		}
		attValues.Width = xd.currentGrid.width(coord(xd.currentGrid.nx) - xd.currentGrid.currentCol + 1)
	}

	ls.HSize = attValues.Width
	ls.LineHeight = leading
	vlist, _ := node.Linebreak(hlist, ls)
	if vlist == nil {
		return nil, nil
	}
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
		eval, err := evaluateXPath(xd, layoutelt, *attValues.Select)
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