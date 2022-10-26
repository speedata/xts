package core

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/color"
	"github.com/speedata/boxesandglue/backend/document"
	"github.com/speedata/boxesandglue/backend/node"
	"github.com/speedata/boxesandglue/csshtml"
	"github.com/speedata/boxesandglue/frontend"
	"github.com/speedata/boxesandglue/frontend/pdfdraw"
	"github.com/speedata/boxesandglue/pdfbackend/pdf"
	"github.com/speedata/goxml"
	xpath "github.com/speedata/goxpath"
	"github.com/speedata/textlayout/harfbuzz"
)

type commandFunc func(*xtsDocument, *goxml.Element) (xpath.Sequence, error)

var (
	dataDispatcher = make(map[string]map[string]*goxml.Element)
	dispatchTable  map[string]commandFunc
	unitRE         = regexp.MustCompile(`(.*?)(sp|mm|cm|in|pt|px|pc|m)`)
	astRE          = regexp.MustCompile(`(\d*)\*`)
)

func init() {
	dispatchTable = map[string]commandFunc{
		"A":                cmdA,
		"Attribute":        cmdAttribute,
		"B":                cmdB,
		"Bookmark":         cmdBookmark,
		"Box":              cmdBox,
		"Circle":           cmdCircle,
		"ClearPage":        cmdClearpage,
		"Column":           cmdColumn,
		"Columns":          cmdColumns,
		"Contents":         cmdContents,
		"Copy-of":          cmdCopyof,
		"DefineColor":      cmdDefineColor,
		"DefineFontfamily": cmdDefineFontfamily,
		"DefineFontsize":   cmdDefineFontsize,
		"DefineMasterpage": cmdDefineMasterpage,
		"DefineTextformat": cmdDefineTextformat,
		"Element":          cmdElement,
		"ForAll":           cmdForall,
		"Group":            cmdGroup,
		"I":                cmdI,
		"Image":            cmdImage,
		"LoadDataset":      cmdLoadDataset,
		"LoadFontfile":     cmdLoadFontfile,
		"Loop":             cmdLoop,
		"Message":          cmdMessage,
		"NextFrame":        cmdNextFrame,
		"Options":          cmdOptions,
		"Pageformat":       cmdPageformat,
		"Paragraph":        cmdParagraph,
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
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	seq, err := dispatch(xd, layoutelt, xd.data)
	if err != nil {
		return nil, err
	}
	hl := document.Hyperlink{URI: attValues.Href}

	te := &frontend.Paragraph{
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
	if err := getXMLAttributes(xd, layoutelt, attValues); err != nil {
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

	te := &frontend.Paragraph{
		Settings: frontend.TypesettingSettings{
			frontend.SettingFontWeight: frontend.FontWeight700,
		},
	}
	getTextvalues(te, seq, "cmdBold", layoutelt.Line)

	return xpath.Sequence{te}, err
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
	dest := getNumDest()
	dest.Callback = func(n node.Node) string {
		startStop := n.(*node.StartStop)
		num := startStop.Value.(int)
		destObj := xd.document.Doc.GetDest(num)
		destObj.X = 0
		title := eval.Stringvalue()
		outline := &pdf.Outline{
			Title: title,
			Dest:  destObj,
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
		Backgroundcolor string          `sdxml:"default:-"`
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

	var bgcolor *color.Color
	if bgc, ok := attrs["background-color"]; ok {
		bgcolor = xd.document.GetColor(bgc)
	} else {
		bgcolor = xd.document.GetColor(attValues.Backgroundcolor)
	}

	r := node.NewRule()
	r.Hide = true
	r.Pre, r.Post = xd.document.HTMLBorder(attValues.Width, attValues.Height, 0, attrs)
	if bgcolor.Space != color.ColorNone {
		r.Pre += pdfdraw.New().Color(*bgcolor).Rect(0, 0, attValues.Width, -attValues.Height).Fill().String()
	} else {

	}

	vl := node.Vpack(r)
	vl.Width = attValues.Width
	vl.Height = attValues.Height

	return xpath.Sequence{vl}, err
}

func cmdCircle(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	attValues := &struct {
		Class           string
		ID              string
		Style           string
		Backgroundcolor *string
		RadiusX         *bag.ScaledPoint `sdxml:"mustexist"`
		RadiusY         *bag.ScaledPoint
	}{}
	if err := getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	if attValues.RadiusY == nil {
		attValues.RadiusY = attValues.RadiusX
	}

	attrs, err := xd.applyLayoutStylesheet(attValues.Class, attValues.ID, attValues.Style, "circle")
	var bgcolor *color.Color
	if attValues.Backgroundcolor == nil {
		if bgc, ok := attrs["background-color"]; ok {
			bgcolor = xd.document.GetColor(bgc)
		}
	} else {
		bgcolor = xd.document.GetColor(*attValues.Backgroundcolor)
	}

	r := node.NewRule()
	if bgcolor == nil {
		bgcolor = xd.document.GetColor("black")
	}

	if bgcolor.Space != color.ColorNone {
		str := pdfdraw.New().Save().Color(*bgcolor).Circle(0, 0, *attValues.RadiusX, *attValues.RadiusY).Fill().String()
		r.Pre = str
		r.Hide = true
		r.Post = pdfdraw.New().Restore().String()
	}

	r.Width = *attValues.RadiusX * 2
	r.Height = *attValues.RadiusY * 2
	vl := node.Vpack(r)
	if vl.Attributes == nil {
		vl.Attributes = node.H{}
	}
	vl.Attributes["origin"] = "circle"

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

func cmdCopyof(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Select string `sdxml:"mustexist"`
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
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

	te := &frontend.Paragraph{
		Settings: frontend.TypesettingSettings{
			frontend.SettingStyle: frontend.FontStyleItalic,
		},
	}
	getTextvalues(te, seq, "cmdBold", layoutelt.Line)

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
	oldContext := xd.data.Ctx.SetContext(xpath.Sequence{})
	xd.data.Ctx.Store = map[any]any{
		"xd": xd,
	}
	if dd, ok := dataDispatcher["root"]; ok {
		if rec, ok := dd[""]; ok {
			_, err = dispatch(xd, rec, xd.data)
		}
	}
	xd.data.Ctx.SetContext(xpath.Sequence{oldContext})
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

func cmdMessage(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
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
	bag.Logger.Infof("Message (line %d): %s", layoutelt.Line, eval.Stringvalue())
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
		ShowHyperlinks *bool
		Title          *string
		Author         *string
		Creator        *string
		Subject        *string
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
		Textformat string
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	seq, err := dispatch(xd, layoutelt, xd.data)
	if err != nil {
		return nil, err
	}

	te := &frontend.Paragraph{
		Settings: make(frontend.TypesettingSettings),
	}
	if attValues.Color != "" {
		te.Settings[frontend.SettingColor] = attValues.Color
	}
	if attValues.Features != "" {
		te.Settings[frontend.SettingOpenTypeFeature] = attValues.Features
	}
	if name := attValues.Textformat; name != "" {
		if tf, ok := xd.textformats[name]; ok {
			if tf.halignment != frontend.HAlignDefault {
				te.Settings[frontend.SettingHAlign] = tf.halignment
			}
		}
	}
	getTextvalues(te, seq, "cmdParagraph", layoutelt.Line)
	return xpath.Sequence{te}, nil
}

func cmdPlaceObject(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	xd.setupPage()
	var err error
	attValues := &struct {
		Allocate  bool `sdxml:"default:yes"`
		Column    string
		Row       string
		Area      string
		Groupname string
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
			row = area.CurrentRow()
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

	xd.store["maxwidth"] = xd.currentGrid.nx - columnInt + 1

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
		if vl.Attributes != nil {
			origin = vl.Attributes["origin"].(string)
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

	if pos == positioningUnknown {
		pos = positioningGrid
		xy, err := xd.currentGrid.findFreeSpaceForObject(vl, area)
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
		xd.OutputAt(vl, col, row, attValues.Allocate, area, origin)
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
			bag.Logger.Errorf("SaveDataset (line %d): error parsing select XPath expression %s", layoutelt.Line, err)
			return nil, err
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

	te := &frontend.Paragraph{
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

	getTextvalues(te, seq, "cmdSpan", layoutelt.Line)

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

func cmdTable(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	var err error
	attValues := &struct {
		Width      bag.ScaledPoint
		Stretch    string `sdxml:"default:no"`
		FontFamily string
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	if attValues.Width == 0 {
		attValues.Width = xd.currentGrid.width(coord(xd.store["maxwidth"].(int)))
	}
	tbl := frontend.Table{}
	tbl.MaxWidth = attValues.Width

	ff := xd.document.FindFontFamily("text")
	if af := attValues.FontFamily; af != "" {
		if fontfamily := xd.document.FindFontFamily(af); fontfamily != nil {
			ff = fontfamily
		}
	}
	tbl.FontFamily = ff

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
	var err error
	attValues := &struct {
		Fontsize   string
		Width      bag.ScaledPoint
		FontFamily string
		Parsep     bag.ScaledPoint
	}{}
	if err = getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}

	leading := xd.fontsizes["text"][1]
	fontsize := xd.fontsizes["text"][0]
	attrFontsize := attValues.Fontsize

	ff := xd.document.FindFontFamily("text")
	if af := attValues.FontFamily; af != "" {
		if fontfamily := xd.document.FindFontFamily(af); fontfamily != nil {
			ff = fontfamily
		}
	}

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
	textblock := node.NewVList()

	for i, itm := range seq {
		te := &frontend.Paragraph{
			Settings: frontend.TypesettingSettings{
				frontend.SettingFontFamily: ff,
				frontend.SettingSize:       fontsize,
			},
		}

		switch t := itm.(type) {
		case *frontend.Paragraph:
			if align := t.Settings[frontend.SettingHAlign]; align != 0 {
				te.Settings[frontend.SettingHAlign] = align
			}
			te.Items = append(te.Items, t)
		case node.Node:
			te.Items = append(te.Items, t)
		default:
			bag.Logger.DPanicf("cmdTextblock: unknown type %T", t)
		}
		if attValues.Width == 0 {
			attValues.Width = xd.currentGrid.width(coord(xd.store["maxwidth"].(int)))
		}

		vlist, _, err := xd.document.FormatParagraph(te, attValues.Width,
			frontend.Leading(leading),
		)
		if err != nil {
			return nil, err
		}
		textblock.List = node.InsertAfter(textblock.List, node.Tail(textblock.List), vlist)
		if i < len(seq) && attValues.Parsep != 0 {
			g := node.NewGlue()
			g.Width = attValues.Parsep
			textblock.List = node.InsertAfter(textblock.List, vlist, g)
		}
	}

	if textblock.List == nil {
		return nil, nil
	}
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
		c := itm.(*frontend.Paragraph)
		tc.Contents = append(tc.Contents, c)
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
