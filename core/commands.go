package core

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	pdf "github.com/speedata/baseline-pdf"
	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/color"
	"github.com/speedata/boxesandglue/backend/document"
	"github.com/speedata/boxesandglue/backend/node"
	"github.com/speedata/boxesandglue/csshtml"
	"github.com/speedata/boxesandglue/frontend"
	"github.com/speedata/boxesandglue/frontend/pdfdraw"
	"github.com/speedata/goxml"
	"github.com/speedata/goxpath"
	xpath "github.com/speedata/goxpath"
	"github.com/speedata/textlayout/harfbuzz"
)

type commandFunc func(*xtsDocument, *goxml.Element) (xpath.Sequence, error)

var (
	dataDispatcher = make(map[string]map[string]*goxml.Element)
	dispatchTable  map[string]commandFunc
	onlyUnitRE     = regexp.MustCompile(`^(sp|mm|cm|in|pt|px|pc|m)$`)
	unitRE         = regexp.MustCompile(`(.*?)(sp|mm|cm|in|pt|px|pc|m)`)
	astRE          = regexp.MustCompile(`(\d*)\*`)
)

func init() {
	dispatchTable = map[string]commandFunc{
		"A":                cmdA,
		"Action":           cmdAction,
		"Attribute":        cmdAttribute,
		"B":                cmdB,
		"Br":               cmdBr,
		"Barcode":          cmdBarcode,
		"Bookmark":         cmdBookmark,
		"Box":              cmdBox,
		"Circle":           cmdCircle,
		"ClearPage":        cmdClearpage,
		"Column":           cmdColumn,
		"Columns":          cmdColumns,
		"Contents":         cmdContents,
		"DefineColor":      cmdDefineColor,
		"DefineFontfamily": cmdDefineFontfamily,
		"DefineFontsize":   cmdDefineFontsize,
		"DefineMasterpage": cmdDefineMasterpage,
		"DefineTextformat": cmdDefineTextformat,
		"Element":          cmdElement,
		"ForAll":           cmdForall,
		"Function":         cmdFunction,
		"Group":            cmdGroup,
		"I":                cmdI,
		"Image":            cmdImage,
		"LoadDataset":      cmdLoadDataset,
		"LoadFontfile":     cmdLoadFontfile,
		"Loop":             cmdLoop,
		"Mark":             cmdMark,
		"Message":          cmdMessage,
		"NextFrame":        cmdNextFrame,
		"NextRow":          cmdNextRow,
		"Options":          cmdOptions,
		"Pageformat":       cmdPageformat,
		"Paragraph":        cmdParagraph,
		"Param":            ignoreFunction,
		"PDFOptions":       cmdPDFOptions,
		"PlaceObject":      cmdPlaceObject,
		"ProcessNode":      cmdProcessNode,
		"Record":           cmdRecord,
		"SaveDataset":      cmdSaveDataset,
		"SetGrid":          cmdSetGrid,
		"SetVariable":      cmdSetVariable,
		"Span":             cmdSpan,
		"Stylesheet":       cmdStylesheet,
		"Switch":           cmdSwitch,
		"Table":            cmdTable,
		"Textblock":        cmdTextblock,
		"Td":               cmdTd,
		"Trace":            cmdTrace,
		"Tr":               cmdTr,
		"U":                cmdU,
		"Value":            cmdValue,
	}
}

func ignoreFunction(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	return nil, nil
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
				return nil, newTypesettingErrorFromString(fmt.Sprintf("layout: element %q unknown", elt.Name))
			}
		}
	}
	return retSequence, nil
}

func cmdA(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Href string
		Link string
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	seq, err := dispatch(xd, layoutelt, xd.data)
	if err != nil {
		return nil, err
	}
	hl := document.Hyperlink{URI: attValues.Href, Local: attValues.Link}

	te := &frontend.Text{
		Settings: frontend.TypesettingSettings{
			frontend.SettingHyperlink: hl,
		},
	}
	err = getTextvalues(te, seq, "cmdA", layoutelt.Line)

	return xpath.Sequence{te}, err
}

func cmdAction(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	seq, err := dispatch(xd, layoutelt, xd.data)
	if err != nil {
		return nil, err
	}
	var ret xpath.Sequence
	for _, itm := range seq {
		if m, ok := itm.(marker); ok {
			var dest *node.StartStop
			if m.pdftarget {
				dest = getNameDest(m.name)
			} else {
				dest = node.NewStartStop()
			}

			dest.Attributes = node.H{
				"page": xd.currentPage,
			}
			dest.ShipoutCallback = func(n node.Node) string {
				startStop := n.(*node.StartStop)
				cp := startStop.Attributes["page"].(*page)
				m.pagenumber = cp.pagenumber
				m.id = <-cp.markerids
				xd.marker[m.name] = m
				return ""
			}

			ret = append(ret, dest)
		}
	}
	return ret, nil
}

func cmdAttribute(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	attValues := &struct {
		Select string `sdxml:"noescape"`
		Name   string `sdxml:"mustexist"`
	}{}
	if err := getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	var eval xpath.Sequence
	var err error
	eval, err = xd.data.Evaluate(attValues.Select)
	if err != nil {
		return nil, newTypesettingErrorFromStringf("Attribute (line %d): error parsing select XPath expression %s", layoutelt.Line, err)
	}

	attr := goxml.Attribute{
		Name:  attValues.Name,
		Value: eval.Stringvalue(),
	}
	return xpath.Sequence{attr}, nil
}

func cmdB(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	seq, err := dispatch(xd, layoutelt, xd.data)
	if err != nil {
		return nil, err
	}
	te := &frontend.Text{
		Settings: frontend.TypesettingSettings{
			frontend.SettingFontWeight: frontend.FontWeight700,
		},
	}
	err = getTextvalues(te, seq, "cmdBold", layoutelt.Line)

	return xpath.Sequence{te}, err
}

func cmdBr(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	te := frontend.NewText()
	te.Items = append(te.Items, "\n")
	return xpath.Sequence{te}, nil
}

func cmdBarcode(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Select     string `sdxml:"mustexist"`
		Type       string
		FontFamily *string
		FontSize   string
		Height     bag.ScaledPoint `sdxml:"mustexist"`
		Width      bag.ScaledPoint `sdxml:"mustexist"`
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	ff := xd.document.FindFontFamily("text")
	if af := attValues.FontFamily; af != nil {
		if fontfamily := xd.document.FindFontFamily(*af); fontfamily != nil {
			ff = fontfamily
		}
	}

	fontsize, _, err := xd.getFontSizeLeading(attValues.FontSize)
	if err != nil {
		return nil, err
	}

	var eval xpath.Sequence
	eval, err = xd.data.Evaluate(attValues.Select)
	if err != nil {
		return nil, newTypesettingErrorFromStringf("Barcode (line %d): error parsing select XPath expression %s", layoutelt.Line, err)
	}

	var bcType int
	switch attValues.Type {
	case "EAN13":
		bcType = barcodeEAN13
	case "Code128":
		bcType = barcodeCode128
	case "QRCode":
		bcType = barcodeQR
	default:
		return nil, fmt.Errorf("Unknown barcode type %q", attValues.Type)
	}
	var bc node.Node
	if bc, err = createBarcode(bcType, eval.Stringvalue(), attValues.Width, attValues.Height, xd, ff, fontsize); err != nil {
		return nil, err
	}

	return xpath.Sequence{bc}, nil
}

func cmdBookmark(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Select string `sdxml:"mustexist"`
		Level  int
		Open   bool
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	var eval xpath.Sequence
	eval, err = xd.data.Evaluate(attValues.Select)
	if err != nil {
		return nil, newTypesettingErrorFromStringf("Bookmark (line %d): error parsing select XPath expression %s", layoutelt.Line, err)
	}

	dest := getNumDest()

	// this callback turns the dest object into an outline object by adding the
	// dest to the Outlines slice of the PDFWriter.
	dest.ShipoutCallback = func(n node.Node) string {
		startStop := n.(*node.StartStop)
		num := startStop.Value.(int)
		destObj := xd.document.Doc.GetNumDest(num)
		destObj.X = 0
		title := eval.Stringvalue()
		outline := &pdf.Outline{
			Title: title,
			Dest:  fmt.Sprintf("[ %s /XYZ %f %f 0]", destObj.PageObjectnumber.Ref(), destObj.X, destObj.Y),
			Open:  attValues.Open,
		}
		curOutlines := &xd.document.Doc.PDFWriter.Outlines
		for i := 1; i < attValues.Level; i++ {
			if len(*curOutlines) == 0 {
				bag.Logger.Errorf("level %d bookmark does not exist for new bookmark (title %s)", i, title)
			} else {
				curOutlines = &(*curOutlines)[len(*curOutlines)-1].Children
			}
		}
		*curOutlines = append(*curOutlines, outline)
		return ""
	}
	return xpath.Sequence{dest}, nil
}

func cmdBox(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	attValues := &struct {
		Class           string
		ID              string
		Style           string
		Backgroundcolor *string
		Width           bag.ScaledPoint `sdxml:"mustexist"`
		Height          bag.ScaledPoint `sdxml:"mustexist"`
	}{}

	if err := getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	attrs, err := xd.applyLayoutStylesheet(attValues.Class, attValues.ID, attValues.Style, "box")
	if err != nil {
		return nil, err
	}

	hv := xd.document.CSSPropertiesToValues(attrs)
	if attValues.Backgroundcolor != nil {
		hv.BackgroundColor = xd.document.GetColor(*attValues.Backgroundcolor)
	} else if _, ok := attrs["background-color"]; ok {
		// already parsed
	}

	vl := node.NewVList()
	vl.Width = attValues.Width - hv.BorderLeftWidth - hv.BorderRightWidth
	vl.Height = attValues.Height - hv.BorderTopWidth - hv.BorderBottomWidth
	vl = xd.document.HTMLBorder(vl, hv)
	return xpath.Sequence{vl}, err
}

func cmdCircle(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	attValues := &struct {
		Class           string
		ID              string
		Style           string
		Backgroundcolor *string
		Borderwidth     *bag.ScaledPoint
		Bordercolor     *string
		RadiusX         *bag.ScaledPoint `sdxml:"mustexist"`
		RadiusY         *bag.ScaledPoint
		OriginX         int `sdxml:"default:50"`
		OriginY         int `sdxml:"default:50"`
	}{}
	if err := getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	if attValues.RadiusY == nil {
		attValues.RadiusY = attValues.RadiusX
	}

	rx, ry := *attValues.RadiusX, *attValues.RadiusY
	attrs, err := xd.applyLayoutStylesheet(attValues.Class, attValues.ID, attValues.Style, "circle")

	var bgcolor, bordercolor *color.Color
	var borderwidth bag.ScaledPoint
	if attValues.Backgroundcolor == nil {
		if bgc, ok := attrs["background-color"]; ok {
			bgcolor = xd.document.GetColor(bgc)
		}
	} else {
		bgcolor = xd.document.GetColor(*attValues.Backgroundcolor)
	}

	if attValues.Bordercolor == nil {
		if bgc, ok := attrs["border-color"]; ok {
			bordercolor = xd.document.GetColor(bgc)
		}
	} else {
		bordercolor = xd.document.GetColor(*attValues.Bordercolor)
	}

	if attValues.Borderwidth == nil {
		if bgc, ok := attrs["border-width"]; ok {
			borderwidth = bag.MustSp(bgc)
		}
	} else {
		borderwidth = *attValues.Borderwidth
	}

	circ := pdfdraw.New().Save().Circle(rx, ry*-1, rx, ry)
	if bgcolor != nil && bgcolor.Space != color.ColorNone && bordercolor != nil && bordercolor.Space != color.ColorNone {
		circ.ColorNonstroking(*bgcolor)
		circ.ColorStroking(*bordercolor).LineWidth(borderwidth)
		circ.StrokeFill()
	} else if bordercolor != nil && bordercolor.Space != color.ColorNone {
		circ.ColorStroking(*bordercolor).LineWidth(borderwidth)
		circ.Stroke()
	} else if bgcolor != nil && bgcolor.Space != color.ColorNone {
		circ.ColorNonstroking(*bgcolor)
		circ.Fill()
	}
	r := node.NewRule()
	r.Pre = circ.String()
	r.Hide = true
	r.Post = pdfdraw.New().Restore().String()

	r.Width = rx * 2
	r.Height = ry * 2
	vl := node.Vpack(r)
	if vl.Attributes == nil {
		vl.Attributes = node.H{}
	}
	vl.Attributes["origin"] = "circle"
	vl.Attributes["shiftX"] = -1 * (r.Width * bag.ScaledPoint(attValues.OriginX) / 100)
	vl.Attributes["shiftY"] = -1 * (r.Height * bag.ScaledPoint(attValues.OriginY) / 100)

	return xpath.Sequence{vl}, err
}

func cmdColumn(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Width string
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	g := node.NewGlue()
	split := strings.Split(attValues.Width, "plus")
	var unitString string
	var stretchString string
	if len(split) == 1 {
		if unitRE.MatchString(split[0]) {
			unitString = split[0]
		} else if astRE.MatchString(split[0]) {
			stretchString = split[0]
		}
	} else {
		if unitRE.MatchString(split[0]) {
			unitString = split[0]
		}
		if astRE.MatchString(split[1]) {
			stretchString = split[1]
		}
	}

	if unitString != "" {
		g.Width = bag.MustSp(unitString)
	}
	if astRE.MatchString(stretchString) {
		astMatch := astRE.FindAllStringSubmatch(stretchString, -1)
		if c := astMatch[0][1]; c != "" {
			stretch, err := strconv.Atoi(c)
			if err != nil {
				return nil, err
			}
			g.Stretch = bag.ScaledPoint(stretch) * bag.Factor
		} else {
			g.Stretch = bag.Factor
		}
		g.StretchOrder = 1
	}
	cs := frontend.ColSpec{
		ColumnWidth: g,
	}
	return xpath.Sequence{cs}, nil
}

func cmdColumns(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	seq, err := dispatch(xd, layoutelt, xd.data)
	if err != nil {
		return nil, err
	}
	return seq, nil
}

func cmdClearpage(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	clearPage(xd)
	return nil, nil
}

func cmdContents(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	_, err := dispatch(xd, layoutelt, xd.data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func cmdDefineFontfamily(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Name string
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
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

func cmdDefineColor(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Name                string
		Colorname           string
		Model               string
		Value               string
		R, G, B, C, M, Y, K string
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	col := color.Color{}
	switch attValues.Model {
	case "cmyk":
		col.Space = color.ColorCMYK
		var c, m, y, k int
		if c, err = strconv.Atoi(attValues.C); err != nil {
			return nil, fmt.Errorf("DefineColor: cannot parse value for c (line %d)", layoutelt.Line)
		}
		if m, err = strconv.Atoi(attValues.M); err != nil {
			return nil, fmt.Errorf("DefineColor: cannot parse value for m (line %d)", layoutelt.Line)
		}
		if y, err = strconv.Atoi(attValues.Y); err != nil {
			return nil, fmt.Errorf("DefineColor: cannot parse value for y (line %d)", layoutelt.Line)
		}
		if k, err = strconv.Atoi(attValues.K); err != nil {
			return nil, fmt.Errorf("DefineColor: cannot parse value for k (line %d)", layoutelt.Line)
		}
		col.C = float64(c) / 100.0
		col.M = float64(m) / 100.0
		col.Y = float64(y) / 100.0
		col.K = float64(k) / 100.0
	case "rgb":
		// 0-100
		col.Space = color.ColorRGB
		var r, g, b int
		if r, err = strconv.Atoi(attValues.R); err != nil {
			return nil, fmt.Errorf("DefineColor: cannot parse value for r (line %d)", layoutelt.Line)
		}
		if g, err = strconv.Atoi(attValues.G); err != nil {
			return nil, fmt.Errorf("DefineColor: cannot parse value for g (line %d)", layoutelt.Line)
		}
		if b, err = strconv.Atoi(attValues.B); err != nil {
			return nil, fmt.Errorf("DefineColor: cannot parse value for b (line %d)", layoutelt.Line)
		}
		col.R = float64(r) / 100.0
		col.G = float64(g) / 100.0
		col.B = float64(b) / 100.0
	case "RGB":
		// 0-255
		col.Space = color.ColorRGB
		var r, g, b int
		if r, err = strconv.Atoi(attValues.R); err != nil {
			return nil, fmt.Errorf("DefineColor: cannot parse value for r (line %d)", layoutelt.Line)
		}
		if g, err = strconv.Atoi(attValues.G); err != nil {
			return nil, fmt.Errorf("DefineColor: cannot parse value for g (line %d)", layoutelt.Line)
		}
		if b, err = strconv.Atoi(attValues.B); err != nil {
			return nil, fmt.Errorf("DefineColor: cannot parse value for b (line %d)", layoutelt.Line)
		}
		col.R = float64(r) / 255.0
		col.G = float64(g) / 255.0
		col.B = float64(b) / 255.0
	case "spotcolor":
		col.Space = color.ColorSpotcolor
	case "":
		// let's hope the user has provided a value field...
		if attValues.Value == "" {
			return nil, fmt.Errorf("DefineColor: empty model not recognized - you should provide a color model or a value attribute (line %d)", layoutelt.Line)
		}
		col = *xd.document.GetColor(attValues.Value)
	default:
		return nil, fmt.Errorf("DefineColor: model %q not recognized (line %d)", attValues.Model, layoutelt.Line)
	}
	xd.document.DefineColor(attValues.Name, &col)
	return nil, nil
}

func cmdDefineFontsize(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	attValues := &struct {
		Name     string
		Fontsize bag.ScaledPoint
		Leading  bag.ScaledPoint
	}{}

	if err := getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	if attValues.Fontsize == 0 || attValues.Leading == 0 {
		if attValues.Fontsize == 0 {
			bag.Logger.Warnf("DefineFontsize (line %d): fontsize is 0", layoutelt.Line)
		}
		if attValues.Leading == 0 {
			bag.Logger.Warnf("DefineFontsize (line %d): leading is 0", layoutelt.Line)
		}
	}
	if xd.fontsizes == nil {
		xd.fontsizes = make(map[string][2]bag.ScaledPoint)
	}
	xd.fontsizes[attValues.Name] = [2]bag.ScaledPoint{attValues.Fontsize, attValues.Leading}
	return nil, nil
}

func cmdDefineMasterpage(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Margin string `sdxml:"mustexist"`
		Name   string `sdxml:"mustexist"`
		Test   string `sdxml:"mustexist"`
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	xd.layoutNS = layoutelt.Namespaces
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

func cmdDefineTextformat(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Alignment string
		Name      string
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	tf := textformat{}
	switch attValues.Alignment {
	case "leftaligned":
		tf.halignment = frontend.HAlignLeft
	case "rightaligned":
		tf.halignment = frontend.HAlignRight
	case "centered":
		tf.halignment = frontend.HAlignCenter
	case "justified":
		tf.halignment = frontend.HAlignJustified
	}
	xd.defineTextformat(attValues.Name, tf)
	return nil, nil
}

func cmdElement(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Name string `sdxml:"mustexist"`
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
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
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	var eval xpath.Sequence
	eval, err = xd.data.Evaluate(attValues.Select)
	if err != nil {
		return nil, newTypesettingErrorFromStringf("ForAll (line %d): error parsing select XPath expression %s", layoutelt.Line, err)
	}
	var ret xpath.Sequence

	oldContext := xd.data.Ctx.SetContextSequence(xpath.Sequence{})
	for i, itm := range eval {
		xd.data.Ctx.SetContextSequence(xpath.Sequence{itm})
		xd.data.Ctx.Pos = i + 1
		neval, err := dispatch(xd, layoutelt, xd.data)
		if err != nil {
			return nil, err
		}
		for _, nitm := range neval {
			ret = append(ret, nitm)
		}
	}
	xd.data.Ctx.SetContextSequence(xpath.Sequence{oldContext})
	return ret, nil
}

func cmdFunction(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Name string `sdxml:"mustexist"`
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	params := []string{}
	for _, cld := range layoutelt.Children() {
		if cldElt, ok := cld.(*goxml.Element); ok && cldElt.Name == "Param" {
			paramValues := &struct {
				Name string `sdxml:"mustexist"`
			}{}
			if err = getXMLAttributes(xd, cldElt, paramValues); err != nil {
				return nil, err
			}
			params = append(params, paramValues.Name)
		}
	}

	prefixName := strings.Split(attValues.Name, ":")
	if len(prefixName) != 2 {
		return nil, newTypesettingErrorFromStringf("Function (line %d): function name needs a namespace prefix", layoutelt.Line)
	}
	var ns string
	var ok bool
	prefix := prefixName[0]
	name := prefixName[1]
	if ns, ok = layoutelt.Namespaces[prefix]; !ok {
		return nil, newTypesettingErrorFromStringf("Function: unknown name space prefix %s", prefix)
	}
	a := func(ctx *xpath.Context, args []xpath.Sequence) (xpath.Sequence, error) {
		sf := returnEvalBodyLater(layoutelt, xd, ctx)
		for i := 0; i < len(params); i++ {
			xd.data.SetVariable(params[i], args[i])
		}
		return sf()
	}
	minArg := len(params)
	maxArg := len(params)
	goxpath.RegisterFunction(&goxpath.Function{Name: name, Namespace: ns, F: a, MinArg: minArg, MaxArg: maxArg})
	return nil, nil
}

func cmdGroup(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Name string `sdxml:"mustexist"`
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
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

func cmdI(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	seq, err := dispatch(xd, layoutelt, xd.data)
	if err != nil {
		return nil, err
	}
	te := &frontend.Text{
		Settings: frontend.TypesettingSettings{
			frontend.SettingStyle: frontend.FontStyleItalic,
		},
	}
	err = getTextvalues(te, seq, "cmdBold", layoutelt.Line)

	return xpath.Sequence{te}, err
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
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	filename, err := xd.cfg.FindFile(attValues.Href)
	if err != nil {
		bag.Logger.Error(err)
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
	if hl.Attributes == nil {
		hl.Attributes = node.H{}
	}
	hl.Attributes["origin"] = "image"

	return xpath.Sequence{hl}, nil
}

func cmdLoadDataset(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Href *string
		Name *string
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	var filename string
	if attValues.Name != nil {
		filename = xd.jobname + "-" + *attValues.Name + ".xml"
	} else if attValues.Href != nil {
		filename = *attValues.Href
	}
	xmlPath, err := xd.cfg.FindFile(filename)
	if xmlPath == "" {
		bag.Logger.Infof("LoadDataset file %s does not exist", filename)
		return nil, nil
	}
	r, err := os.Open(xmlPath)
	if err != nil {
		return nil, err
	}
	bag.Logger.Infof("LoadDataset file %s loaded", filename)
	saveData := xd.data
	defer r.Close()
	xd.data, err = xpath.NewParser(r)
	if err != nil {
		return nil, err
	}
	oldContext := xd.data.Ctx.SetContextSequence(xpath.Sequence{})
	xd.data.Ctx.Store = map[any]any{
		"xd": xd,
	}
	if dd, ok := dataDispatcher["root"]; ok {
		if rec, ok := dd[""]; ok {
			_, err = dispatch(xd, rec, xd.data)
		}
	}
	xd.data.Ctx.SetContextSequence(xpath.Sequence{oldContext})
	xd.data = saveData

	return nil, nil
}

func cmdLoadFontfile(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Filename string `sdxml:"mustexist"`
		Name     string `sdxml:"mustexist"`
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	fn, err := xd.cfg.FindFile(attValues.Filename)
	if err != nil {
		bag.Logger.Errorf("error in line %d", layoutelt.Line)
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
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	var eval xpath.Sequence
	eval, err = evaluateXPath(xd, layoutelt.Namespaces, attValues.Select)
	if err != nil {
		bag.Logger.Errorf("ProcessNode (line %d): error parsing select XPath expression %s", layoutelt.Line, err)
		return nil, err
	}

	oldContext := xd.data.Ctx.SetContextSequence(xpath.Sequence{})

	for i, itm := range eval {
		xd.data.Ctx.Pos = i + 1
		if elt, ok := itm.(*goxml.Element); ok {
			xd.data.Ctx.SetContextSequence(xpath.Sequence{elt})
			if dd, ok := dataDispatcher[elt.Name]; ok {
				if rec, ok := dd[attValues.Mode]; ok {
					_, err = dispatch(xd, rec, xd.data)
				}
			}
		}
	}

	xd.data.Ctx.SetContextSequence(xpath.Sequence{oldContext})
	return nil, nil
}

func cmdRecord(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Element string `sdxml:"mustexist"`
		Mode    string
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	dp := dataDispatcher[attValues.Element]
	if dp == nil {
		dataDispatcher[attValues.Element] = make(map[string]*goxml.Element)
	}
	dataDispatcher[attValues.Element][attValues.Mode] = layoutelt
	return nil, nil
}

func cmdLoop(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Select   string `sdxml:"noescape,mustexist"`
		Variable string
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	var eval xpath.Sequence

	xd.layoutNS = layoutelt.Namespaces
	eval, err = evaluateXPath(xd, layoutelt.Namespaces, attValues.Select)
	if err != nil {
		return nil, err
	}
	f, err := strconv.ParseFloat(eval.Stringvalue(), 64)
	if err != nil {
		return nil, err
	}

	for i := 1; i < int(f)+1; i++ {
		xd.data.SetVariable(attValues.Variable, xpath.Sequence{i})
		eval, err = dispatch(xd, layoutelt, xd.data)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func cmdMark(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Select    string `sdxml:"noescape"`
		Append    bool
		PDFTarget bool
		ShiftUP   bag.ScaledPoint
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	eval, err := evaluateXPath(xd, layoutelt.Namespaces, attValues.Select)
	if err != nil {
		return nil, err
	}
	var seq xpath.Sequence
	for _, itm := range eval {
		seq = append(seq, marker{
			name:      itm.(string),
			append:    attValues.Append,
			pdftarget: attValues.PDFTarget,
			shiftup:   attValues.ShiftUP,
		})
	}

	return seq, nil
}

func cmdMessage(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Type   *string
		Select *string `sdxml:"noescape"`
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	var eval xpath.Sequence
	if attValues.Select != nil {
		eval, err = evaluateXPath(xd, layoutelt.Namespaces, *attValues.Select)
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
	f := bag.Logger.Infof
	if t := attValues.Type; t != nil {
		switch *t {
		case "debug":
			f = bag.Logger.Debugf
		case "warning":
			f = bag.Logger.Warnf
		case "error":
			f = bag.Logger.Errorf
		}
	}
	f("Message (line %d): %s", layoutelt.Line, eval.Stringvalue())
	return nil, nil
}

func cmdNextFrame(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Area string `sdxml:"mustexist"`
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	xd.setupPage()
	if area, ok := xd.currentPage.pagegrid.areas[attValues.Area]; ok {
		area.currentFrame++
		if area.currentFrame == len(area.frame) {
			area.currentFrame = 0
			clearPage(xd)
		}
	}
	return nil, nil
}

func cmdPDFOptions(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Author            *string
		Creator           *string
		DisplayMode       *string
		Duplex            *string
		PickTrayByPDFSize *bool
		PrintScaling      *string
		ShowHyperlinks    *bool
		Subject           *string
		Title             *string
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	if attValues.Author != nil {
		xd.document.Doc.Author = *attValues.Author
	}
	if attValues.Creator != nil {
		xd.document.Doc.Creator = *attValues.Creator
	}
	if dm := attValues.DisplayMode; dm != nil {
		switch *dm {
		case "attachments":
			xd.document.Doc.ViewerPreferences["PageMode"] = "/UseAttachments"
		case "bookmarks":
			xd.document.Doc.ViewerPreferences["PageMode"] = "/UseOutlines"
		case "fullscreen":
			xd.document.Doc.ViewerPreferences["PageMode"] = "/FullScreen"
		case "none":
			xd.document.Doc.ViewerPreferences["PageMode"] = "/UseNone"
		case "thumbnails":
			xd.document.Doc.ViewerPreferences["PageMode"] = "/UseThumbs"
		}
	}
	if dplx := attValues.Duplex; dplx != nil {
		switch *dplx {
		case "simplex":
			xd.document.Doc.ViewerPreferences["Duplex"] = "/Simplex"
		case "duplexflipshortedge":
			xd.document.Doc.ViewerPreferences["Duplex"] = "/DuplexFlipShortEdge"
		case "duplexfliplongedge":
			xd.document.Doc.ViewerPreferences["Duplex"] = "/DuplexFlipLongEdge"
		default:
			return nil, newTypesettingErrorFromStringf("Unknown PDFOptions setting %s", *dplx)
		}
	}
	if attValues.PickTrayByPDFSize != nil {
		xd.document.Doc.ViewerPreferences["PickTrayByPDFSize"] = fmt.Sprintf("%t", *attValues.PickTrayByPDFSize)
	}
	if ps := attValues.PrintScaling; ps != nil {
		switch *ps {
		case "appdefault":
			xd.document.Doc.ViewerPreferences["PrintScaling"] = "/AppDefault"
		case "none":
			xd.document.Doc.ViewerPreferences["PrintScaling"] = "/None"
		default:
			delete(xd.document.Doc.ViewerPreferences, "PrintScaling")
		}
	}
	if attValues.Title != nil {
		xd.document.Doc.Title = *attValues.Title
	}
	if attValues.Subject != nil {
		xd.document.Doc.Subject = *attValues.Subject
	}
	if attValues.ShowHyperlinks != nil {
		xd.document.Doc.ShowHyperlinks = *attValues.ShowHyperlinks
	}
	return nil, nil
}

func cmdNextRow(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	xd.setupPage()
	var err error
	attValues := &struct {
		Area string
		Row  *int
		Rows *int
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	if attValues.Area == "" {
		attValues.Area = defaultAreaName
	}
	var area *area
	var ok bool
	if area, ok = xd.currentGrid.areas[attValues.Area]; !ok {
		return nil, fmt.Errorf("area %s not found", attValues.Area)
	}
	if attValues.Row != nil {
		area.SetCurrentRow(coord(*attValues.Row))
	} else if r := attValues.Rows; r != nil {
		xd.currentGrid.nextRow(area)
		area.SetCurrentRow(area.CurrentRow() + coord(*r-1))
	} else {
		xd.currentGrid.nextRow(area)
	}
	area.SetCurrentCol(1)
	return nil, nil
}

func cmdOptions(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Mainlanguage *string
		Bleed        *bag.ScaledPoint
		Cutmarks     *bool
		Features     *string
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
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
	if attValues.Cutmarks != nil {
		xd.document.Doc.ShowCutmarks = *attValues.Cutmarks
	}
	if features := attValues.Features; features != nil {
		for _, str := range strings.Split(*features, ",") {
			f, err := harfbuzz.ParseFeature(str)
			if err != nil {
				bag.Logger.Errorf("cannot parse OpenType feature tag %q.", str)
			}
			xd.document.DefaultFeatures = append(xd.document.DefaultFeatures, f)
		}
	}

	return xpath.Sequence{}, nil
}

func cmdPageformat(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Width  bag.ScaledPoint `sdxml:"mustexist"`
		Height bag.ScaledPoint `sdxml:"mustexist"`
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	xd.document.Doc.DefaultPageWidth = attValues.Width
	xd.document.Doc.DefaultPageHeight = attValues.Height
	return xpath.Sequence{}, nil
}

func cmdParagraph(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Color      string
		Features   string
		Textformat *string
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	seq, err := dispatch(xd, layoutelt, xd.data)
	if err != nil {
		return nil, err
	}
	if seq == nil {
		seq = xpath.Sequence{}
	}
	te := &frontend.Text{
		Settings: make(frontend.TypesettingSettings),
	}
	if attValues.Color != "" {
		te.Settings[frontend.SettingColor] = attValues.Color
	}
	if attValues.Features != "" {
		te.Settings[frontend.SettingOpenTypeFeature] = attValues.Features
	}
	if name := attValues.Textformat; name != nil {
		if tf, ok := xd.textformats[*name]; ok {
			if tf.halignment != frontend.HAlignDefault {
				te.Settings[frontend.SettingHAlign] = tf.halignment
			}
		} else {
			bag.Logger.Warnf("text format %q not found", *name)
		}
	}
	err = getTextvalues(te, seq, "cmdParagraph", layoutelt.Line)
	return xpath.Sequence{te}, err
}

func cmdPlaceObject(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	xd.setupPage()
	var err error
	attValues := &struct {
		Allocate        bool `sdxml:"default:yes"`
		Area            string
		Background      bool
		BackgroundColor string
		Column          string
		Frame           bool
		Row             string
		Groupname       string
		HAlign          string
		HReference      string
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	if attValues.Area == "" {
		attValues.Area = defaultAreaName
	}
	var area *area
	var ok bool
	if area, ok = xd.currentGrid.areas[attValues.Area]; !ok {
		return nil, newTypesettingErrorFromStringf(fmt.Sprintf("area %s not found", attValues.Area))
	}

	pos := positioningUnknown

	var columnInt, rowInt int = 1, 1
	var col, row coord
	var columnLength, rowLength bag.ScaledPoint

	var rowSet, colSet bool
	if attValues.Column != "" {
		if colF, err := strconv.ParseFloat(attValues.Column, 64); err == nil {
			colSet = true
			pos = positioningGrid
			columnInt = int(colF)
			col = coord(columnInt)
		} else {
			bag.Logger.Debug(err)
		}
	}

	xd.store["maxwidth"] = int(area.frame[area.currentFrame].width-area.CurrentCol()) + 1

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
	var shiftX, shiftY bag.ScaledPoint
	var vl *node.VList
	switch t := seq[0].(type) {
	case *node.VList:
		vl = t
		if vl.Attributes != nil {
			origin = vl.Attributes["origin"].(string)
			if sx, ok := vl.Attributes["shiftX"]; ok {
				shiftX = sx.(bag.ScaledPoint)
			}
			if sy, ok := vl.Attributes["shiftY"]; ok {
				shiftY = sy.(bag.ScaledPoint)
			}
		}
	case []*node.VList:
		vl = t[0]
		if vl.Attributes != nil {
			origin = vl.Attributes["origin"].(string)
		}
	case *node.HList:
		vl = node.Vpack(t)
		if t.Attributes != nil {
			origin = t.Attributes["origin"].(string)
		}
	default:
		bag.Logger.DPanicf("PlaceObject: unknown node %T", t)
	}
	if attValues.Frame {
		r := node.NewRule()
		r.Hide = true
		r.Pre = pdfdraw.NewStandalone().Rect(0, 0, vl.Width, -vl.Height).Stroke().String()
		vl.List = node.InsertBefore(vl.List, vl.List, r)
	}
	if attValues.Background {
		col := xd.document.GetColor(attValues.BackgroundColor)
		r := node.NewRule()
		r.Hide = true
		r.Pre = pdfdraw.NewStandalone().Rect(0, 0, vl.Width, -vl.Height).ColorNonstroking(*col).Fill().String()
		vl.List = node.InsertBefore(vl.List, vl.List, r)
	}

	if rowInt, err = strconv.Atoi(attValues.Row); err == nil {
		rowSet = true
		pos = positioningGrid
		row = coord(rowInt)
	}
	if pos == positioningGrid && colSet != rowSet {
		if !colSet {
			col = area.CurrentCol()
			columnInt = int(col)
		}
		if !rowSet {
			wdCols := xd.currentGrid.widthToColumns(vl.Width)
			htCols := xd.currentGrid.heightToRows(vl.Height + vl.Depth)
			row = xd.currentGrid.findSuitableRow(wdCols, htCols, col, area)
			rowInt = int(row)
		}
	} else if pos == positioningUnknown {
		if columnLength, err = bag.Sp(attValues.Column); err == nil {
			pos = positioningAbsolute
		}
		if rowLength, err = bag.Sp(attValues.Row); err == nil {
			pos = positioningAbsolute
		}
	}

	if pos == positioningUnknown {
		pos = positioningGrid
		startCol := area.CurrentCol()
		if startCol+xd.currentGrid.widthToColumns(vl.Width) > coord(xd.currentGrid.nx) {
			startCol = 1
		}
		wdCols := xd.currentGrid.widthToColumns(vl.Width)
		htCols := xd.currentGrid.heightToRows(vl.Height + vl.Depth)
		row = xd.currentGrid.findSuitableRow(wdCols, htCols, startCol, area)
		bag.Logger.Debugf("looking for free space for %s", origin)
		col = startCol

	}
	halign := frontend.HAlignLeft
	if attValues.HAlign == "right" {
		halign = frontend.HAlignRight
	}

	switch pos {
	case positioningAbsolute:
		if attValues.HReference == "right" {
			columnLength -= vl.Width
		}
		columnLength += shiftX
		rowLength += shiftY
		xd.currentPage.outputAbsolute(columnLength, rowLength, vl)
	case positioningGrid:
		if attValues.HReference == "right" {
			wd := xd.currentGrid.widthToColumns(vl.Width)
			col = col - wd + 1
		}
		if shiftX != 0 {
			col = col + xd.currentGrid.widthToColumns(shiftX)
		}
		if shiftY != 0 {
			row = row + xd.currentGrid.heightToRows(shiftY)
		}
		xd.OutputAt(vl, col, row, attValues.Allocate, area, origin, halign)

		// if the current column is right of the area, go to the start of the
		// next row below the object.
		if col+xd.currentGrid.widthToColumns(vl.Width) > area.frame[area.currentFrame].width {
			area.SetCurrentRow(row + xd.currentGrid.heightToRows(vl.Height+vl.Depth))
			area.SetCurrentCol(1)
		}
	}
	return seq, nil
}

func cmdSaveDataset(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	attValues := &struct {
		Href        *string
		Name        *string
		Elementname string
		Select      *string
	}{}
	if err := getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	var eval xpath.Sequence
	var err error
	if attValues.Select != nil {
		eval, err = xd.data.Evaluate(*attValues.Select)
		if err != nil {
			return nil, newTypesettingErrorFromStringf("SaveDataset (line %d): error parsing select XPath expression %s", layoutelt.Line, err)
		}
	} else {
		eval, err = dispatch(xd, layoutelt, xd.data)
		if err != nil {
			return nil, err
		}
	}
	root := goxml.NewElement()
	root.Name = attValues.Elementname
	for _, itm := range eval {
		if elt, ok := itm.(goxml.Element); ok {
			root.Append(elt)
		}
	}

	var filename string
	if attValues.Name != nil {
		filename = xd.jobname + "-" + *attValues.Name + ".xml"
	} else if attValues.Href != nil {
		filename = *attValues.Href
	}

	if filename == "" {
		return nil, fmt.Errorf("SaveDataset (line %d) filename must be provided via name or href", layoutelt.Line)
	}

	w, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	defer w.Close()
	bag.Logger.Infof("Write XML file %s", filename)
	_, err = w.Write([]byte(root.ToXML()))

	return nil, err
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
	if err := getXMLAttributes(xd, layoutelt, attValues); err != nil {
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
	if err := getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	var eval xpath.Sequence
	var err error
	if attValues.Select != nil {
		eval, err = xd.data.Evaluate(*attValues.Select)
		if err != nil {
			return nil, newTypesettingErrorFromStringf("SetVariable (line %d): error parsing select XPath expression %s", layoutelt.Line, err)
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

func cmdSpan(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	attValues := &struct {
		Class string
		ID    string
		Style string
	}{}

	if err := getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	attrs, err := xd.applyLayoutStylesheet(attValues.Class, attValues.ID, attValues.Style, "span")
	if err != nil {
		return nil, err
	}

	seq, err := dispatch(xd, layoutelt, xd.data)
	if err != nil {
		return nil, err
	}

	te := &frontend.Text{
		Settings: frontend.TypesettingSettings{},
	}

	if val, ok := attrs["color"]; ok {
		te.Settings[frontend.SettingColor] = xd.document.GetColor(val)
	}
	if val, ok := attrs["font-weight"]; ok {
		te.Settings[frontend.SettingFontWeight] = frontend.ResolveFontWeight(val, frontend.FontWeight100)
	}
	if val, ok := attrs["font-style"]; ok {
		te.Settings[frontend.SettingStyle] = frontend.ResolveFontStyle(val)
	}

	err = getTextvalues(te, seq, "cmdSpan", layoutelt.Line)

	return xpath.Sequence{te}, err
}

func cmdStylesheet(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Scope string
		Href  string
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
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
							return nil, newTypesettingErrorFromStringf("Case (line %d): error parsing test XPath expression %s", layoutelt.Line, err)
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

func cmdTable(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	xd.setupPage()
	var err error
	attValues := &struct {
		Width      bag.ScaledPoint
		Stretch    bool
		FontFamily string
		FontSize   string
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	if attValues.Width == 0 {
		var mw coord
		if mwInt, ok := xd.store["maxwidth"].(int); ok {
			mw = coord(mwInt)
		} else {
			mw = coord(xd.currentGrid.nx)
		}
		attValues.Width = xd.currentGrid.width(mw)
	}

	fontsize, leading, err := xd.getFontSizeLeading(attValues.FontSize)
	if err != nil {
		return nil, err
	}

	tbl := frontend.Table{}
	tbl.MaxWidth = attValues.Width
	tbl.Stretch = attValues.Stretch

	ff := xd.document.FindFontFamily("text")
	if af := attValues.FontFamily; af != "" {
		if fontfamily := xd.document.FindFontFamily(af); fontfamily != nil {
			ff = fontfamily
		}
	}
	tbl.FontFamily = ff
	tbl.FontSize = fontsize
	tbl.Leading = leading

	seq, err := dispatch(xd, layoutelt, xd.data)
	if err != nil {
		return nil, err
	}

	for _, itm := range seq {
		switch t := itm.(type) {
		case frontend.TableRow:
			tbl.Rows = append(tbl.Rows, &t)
		case frontend.ColSpec:
			tbl.ColSpec = append(tbl.ColSpec, t)
		default:
			fmt.Println(t)
		}
	}
	vls, err := xd.document.BuildTable(&tbl)
	if err != nil {
		return nil, err
	}
	return xpath.Sequence{vls}, nil
}

func cmdTextblock(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	xd.setupPage()
	var err error
	attValues := &struct {
		Fontsize   string
		Width      bag.ScaledPoint
		FontFamily string
		Parsep     bag.ScaledPoint
		Textformat *string `sdxml:"default:text"`
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	ff := xd.document.FindFontFamily("text")
	if af := attValues.FontFamily; af != "" {
		if fontfamily := xd.document.FindFontFamily(af); fontfamily != nil {
			ff = fontfamily
		}
	}
	fontsize, leading, err := xd.getFontSizeLeading(attValues.Fontsize)
	if err != nil {
		return nil, err
	}
	seq, err := dispatch(xd, layoutelt, xd.data)
	if err != nil {
		return nil, err
	}
	var vlists node.Node
	for i, itm := range seq {
		te := &frontend.Text{
			Settings: frontend.TypesettingSettings{
				frontend.SettingFontFamily: ff,
				frontend.SettingSize:       fontsize,
			},
		}
		if name := attValues.Textformat; name != nil {
			if tf, ok := xd.textformats[*name]; ok {
				if tf.halignment != frontend.HAlignDefault {
					te.Settings[frontend.SettingHAlign] = tf.halignment
				}
			} else {
				bag.Logger.Warnf("text format %q not found", *name)
			}
		}

		switch t := itm.(type) {
		case *frontend.Text:
			if align, found := t.Settings[frontend.SettingHAlign]; found && align != frontend.HAlignDefault {
				te.Settings[frontend.SettingHAlign] = align
			}
			te.Items = append(te.Items, t)
		case node.Node:
			te.Items = append(te.Items, t)
		default:
			bag.Logger.DPanicf("cmdTextblock: unknown type %T", t)
		}
		// if no width is requested, we use the maximum available width
		if attValues.Width == 0 {
			var mw coord
			if mwInt, ok := xd.store["maxwidth"].(int); ok {
				mw = coord(mwInt)
			} else {
				mw = coord(xd.currentGrid.nx)
			}
			attValues.Width = xd.currentGrid.width(mw)
		}

		vlist, _, err := xd.document.FormatParagraph(te, attValues.Width, frontend.Leading(leading))
		if err != nil {
			return nil, err
		}
		vlists = node.InsertAfter(vlists, node.Tail(vlists), vlist)
		if i < len(seq) && attValues.Parsep != 0 {
			g := node.NewGlue()
			g.Width = attValues.Parsep
			vlists = node.InsertAfter(vlists, vlist, g)
		}
	}

	if vlists == nil {
		return nil, nil
	}
	textblock := node.Vpack(vlists)
	if textblock.Attributes == nil {
		textblock.Attributes = node.H{}
	}
	textblock.Attributes["origin"] = "textblock"
	return xpath.Sequence{textblock}, nil
}

func cmdTr(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Valign string
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	seq, err := dispatch(xd, layoutelt, xd.data)
	if err != nil {
		return nil, err
	}
	tr := frontend.TableRow{}
	for _, itm := range seq {
		c := itm.(frontend.TableCell)
		tr.Cells = append(tr.Cells, &c)
	}
	switch attValues.Valign {
	case "top":
		tr.VAlign = frontend.VAlignTop
	case "bottom":
		tr.VAlign = frontend.VAlignBottom
	case "middle":
		tr.VAlign = frontend.VAlignMiddle
	}

	return xpath.Sequence{tr}, nil
}

func cmdTd(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Align             string
		BorderBottom      *bag.ScaledPoint `sdxml:"attr:border-bottom"`
		BorderTop         *bag.ScaledPoint `sdxml:"attr:border-top"`
		BorderLeft        *bag.ScaledPoint `sdxml:"attr:border-left"`
		BorderRight       *bag.ScaledPoint `sdxml:"attr:border-right"`
		BorderBottomColor *string          `sdxml:"attr:border-bottom-color"`
		BorderTopColor    *string          `sdxml:"attr:border-top-color"`
		BorderLeftColor   *string          `sdxml:"attr:border-left-color"`
		BorderRightColor  *string          `sdxml:"attr:border-right-color"`
		Padding           *bag.ScaledPoint
		PaddingBottom     *bag.ScaledPoint `sdxml:"attr:padding-bottom"`
		PaddingTop        *bag.ScaledPoint `sdxml:"attr:padding-top"`
		PaddingLeft       *bag.ScaledPoint `sdxml:"attr:padding-left"`
		PaddingRight      *bag.ScaledPoint `sdxml:"attr:padding-right"`
		Colspan           int
		Rowspan           int
		Valign            string
		Class             string
		ID                string
		Style             string
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	seq, err := dispatch(xd, layoutelt, xd.data)
	if err != nil {
		return nil, err
	}

	tc := frontend.TableCell{}
	for _, itm := range seq {
		switch t := itm.(type) {
		case *frontend.Text:
			tc.Contents = append(tc.Contents, t)
		case *node.HList:
			tc.Contents = append(tc.Contents, t)
		default:
			bag.Logger.DPanicf("pack table cell, unknown type %T\n", t)
		}

	}

	attrs, err := xd.applyLayoutStylesheet(attValues.Class, attValues.ID, attValues.Style, "table", "tr", "td")
	if err != nil {
		return nil, err
	}

	if wd, ok := attrs["border-bottom-width"]; ok {
		tc.BorderBottomWidth, err = bag.Sp(wd)
		if err != nil {
			return nil, err
		}
	}
	if wd, ok := attrs["border-top-width"]; ok {
		tc.BorderTopWidth, err = bag.Sp(wd)
		if err != nil {
			return nil, err
		}
	}
	if wd, ok := attrs["border-left-width"]; ok {
		tc.BorderLeftWidth, err = bag.Sp(wd)
		if err != nil {
			return nil, err
		}
	}
	if wd, ok := attrs["border-right-width"]; ok {
		tc.BorderRightWidth, err = bag.Sp(wd)
		if err != nil {
			return nil, err
		}
	}
	if wd, ok := attrs["padding-bottom"]; ok {
		tc.PaddingBottom, err = bag.Sp(wd)
		if err != nil {
			return nil, err
		}
	}
	if wd, ok := attrs["padding-top"]; ok {
		tc.PaddingTop, err = bag.Sp(wd)
		if err != nil {
			return nil, err
		}
	}
	if wd, ok := attrs["padding-left"]; ok {
		tc.PaddingLeft, err = bag.Sp(wd)
		if err != nil {
			return nil, err
		}
	}
	if wd, ok := attrs["padding-right"]; ok {
		tc.PaddingRight, err = bag.Sp(wd)
		if err != nil {
			return nil, err
		}
	}
	if col, ok := attrs["border-bottom-color"]; ok {
		tc.BorderBottomColor = xd.document.GetColor(col)
		if err != nil {
			return nil, err
		}
	}
	if col, ok := attrs["border-top-color"]; ok {
		tc.BorderTopColor = xd.document.GetColor(col)
		if err != nil {
			return nil, err
		}
	}
	if col, ok := attrs["border-left-color"]; ok {
		tc.BorderLeftColor = xd.document.GetColor(col)
		if err != nil {
			return nil, err
		}
	}
	if col, ok := attrs["border-right-color"]; ok {
		tc.BorderRightColor = xd.document.GetColor(col)
		if err != nil {
			return nil, err
		}
	}
	if bb := attValues.BorderBottom; bb != nil {
		tc.BorderBottomWidth = *bb
	}
	if bb := attValues.BorderTop; bb != nil {
		tc.BorderTopWidth = *bb
	}
	if bb := attValues.BorderLeft; bb != nil {
		tc.BorderLeftWidth = *bb
	}
	if bb := attValues.BorderRight; bb != nil {
		tc.BorderRightWidth = *bb
	}
	if attValues.BorderBottomColor != nil {
		tc.BorderBottomColor = xd.document.GetColor(*attValues.BorderBottomColor)
	}
	if attValues.BorderTopColor != nil {
		tc.BorderTopColor = xd.document.GetColor(*attValues.BorderTopColor)
	}
	if attValues.BorderLeftColor != nil {
		tc.BorderLeftColor = xd.document.GetColor(*attValues.BorderLeftColor)
	}
	if attValues.BorderRightColor != nil {
		tc.BorderRightColor = xd.document.GetColor(*attValues.BorderRightColor)
	}
	if p := attValues.Padding; p != nil {
		tc.PaddingLeft = *p
		tc.PaddingRight = *p
		tc.PaddingTop = *p
		tc.PaddingBottom = *p
	}
	if p := attValues.PaddingRight; p != nil {
		tc.PaddingRight = *p
	}
	if p := attValues.PaddingLeft; p != nil {
		tc.PaddingLeft = *p
	}
	if p := attValues.PaddingTop; p != nil {
		tc.PaddingTop = *p
	}
	if p := attValues.PaddingBottom; p != nil {
		tc.PaddingBottom = *p
	}

	switch attValues.Valign {
	case "top":
		tc.VAlign = frontend.VAlignTop
	case "bottom":
		tc.VAlign = frontend.VAlignBottom
	case "middle":
		tc.VAlign = frontend.VAlignMiddle
	}

	switch attValues.Align {
	case "left":
		tc.HAlign = frontend.HAlignLeft
	case "right":
		tc.HAlign = frontend.HAlignRight
	case "center":
		tc.HAlign = frontend.HAlignCenter
	case "justified":
		tc.HAlign = frontend.HAlignJustified
	}

	if cs := attValues.Colspan; cs != 0 {
		tc.ExtraColspan = cs - 1
	}
	if rs := attValues.Rowspan; rs != 0 {
		tc.ExtraRowspan = rs - 1
	}
	return xpath.Sequence{tc}, nil
}

func cmdTrace(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Grid           *bool
		Hyphenation    *bool
		Gridallocation *bool
		Hyperlinks     *bool
		Dests          *bool
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	if attValues.Grid != nil {
		if *attValues.Grid {
			xd.SetVTrace(VTraceGrid)
		} else {
			xd.ClearVTrace(VTraceGrid)
		}
	}

	if attValues.Dests != nil {
		if *attValues.Dests {
			xd.document.Doc.SetVTrace(document.VTraceDest)
		} else {
			xd.document.Doc.ClearVTrace(document.VTraceDest)
		}
	}

	if attValues.Hyphenation != nil {
		if *attValues.Hyphenation {
			xd.SetVTrace(VTraceHyphenation)
		} else {
			xd.ClearVTrace(VTraceHyphenation)
		}
	}

	if attValues.Hyperlinks != nil {
		if *attValues.Hyperlinks {
			xd.document.Doc.SetVTrace(document.VTraceHyperlinks)
		} else {
			xd.document.Doc.ClearVTrace(document.VTraceHyperlinks)
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

func cmdU(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	seq, err := dispatch(xd, layoutelt, xd.data)
	if err != nil {
		return nil, err
	}
	te := &frontend.Text{
		Settings: frontend.TypesettingSettings{
			frontend.SettingTextDecorationLine: frontend.TextDecorationUnderline,
		},
	}
	err = getTextvalues(te, seq, "cmdUnderline", layoutelt.Line)

	return xpath.Sequence{te}, err
}

func cmdValue(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Select *string
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	if attValues.Select != nil {
		eval, err := evaluateXPath(xd, layoutelt.Namespaces, *attValues.Select)
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
