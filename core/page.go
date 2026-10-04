package core

import (
	"fmt"
	"log/slog"

	"github.com/boxesandglue/boxesandglue/backend/bag"
	"github.com/boxesandglue/boxesandglue/backend/document"
	"github.com/boxesandglue/boxesandglue/backend/node"
	"github.com/boxesandglue/boxesandglue/frontend"
	"github.com/boxesandglue/htmlbag"
	"github.com/speedata/goxml"
)

const (
	pageAreaName    string = "__pagearea"
	defaultAreaName string = "__pagearea"
)

type pagetype struct {
	name         string
	test         string
	areas        map[string]area
	marginLeft   bag.ScaledPoint // the default left margin
	marginRight  bag.ScaledPoint // the default right margin
	marginTop    bag.ScaledPoint // the default top margin
	marginBottom bag.ScaledPoint // the default bottom margin
	// marginFromXML is true when the margin attribute of DefineMasterPage
	// is set. The attribute owns the page geometry (it defines the grid);
	// without it the margins come from the matching CSS @page rule.
	marginFromXML bool
	// cssPage is the @page rule matched by name (see resolvePageCSS), nil
	// when the stylesheets carry no @page rule at all.
	cssPage *htmlbag.Page
	// cssMarginWarned prevents repeating the margin mismatch warning on
	// every page.
	cssMarginWarned bool
	layoutElt       *goxml.Element
}

// cssDefaultPageMargin is the margin used when neither the DefineMasterPage
// attribute nor the @page rule sets one, matching the CSS default in
// htmlbag.
var cssDefaultPageMargin = bag.MustSP("1cm")

// resolvePageCSS couples the master page to the CSS @page rule of the same
// name (with the generic @page rule as base). It applies the margin rule:
// the margin attribute of DefineMasterPage is authoritative, the @page
// margins fill in only when the attribute is absent, and a conflict between
// the two is reported once.
//
// The lookup runs at page creation rather than at DefineMasterPage time,
// so a StyleSheet that follows the master page definition in the layout is
// still honoured.
func (xd *xtsDocument) resolvePageCSS(pt *pagetype) error {
	pg, ok := xd.layoutcss.NamedPage(pt.name)
	if !ok {
		pt.cssPage = nil
		if !pt.marginFromXML {
			pt.marginTop, pt.marginBottom, pt.marginLeft, pt.marginRight = cssDefaultPageMargin, cssDefaultPageMargin, cssDefaultPageMargin, cssDefaultPageMargin
		}
		return nil
	}
	pt.cssPage = pg
	type side struct {
		name string
		css  string
		xml  *bag.ScaledPoint
	}
	sides := []side{
		{"top", pg.MarginTop, &pt.marginTop},
		{"bottom", pg.MarginBottom, &pt.marginBottom},
		{"left", pg.MarginLeft, &pt.marginLeft},
		{"right", pg.MarginRight, &pt.marginRight},
	}
	for _, sd := range sides {
		if sd.css == "" {
			if !pt.marginFromXML {
				*sd.xml = cssDefaultPageMargin
			}
			continue
		}
		v, err := bag.SP(sd.css)
		if err != nil {
			return fmt.Errorf("master page %q: cannot parse @page margin-%s %q: %w", pt.name, sd.name, sd.css, err)
		}
		if !pt.marginFromXML {
			*sd.xml = v
			continue
		}
		if v != *sd.xml && !pt.cssMarginWarned {
			pt.cssMarginWarned = true
			slog.Warn("Margin of master page differs from the CSS @page rule, the margin attribute of DefineMasterPage wins", "masterpage", pt.name, "side", sd.name, "attribute", sd.xml.String()+"pt", "css", sd.css)
		}
	}
	return nil
}

func (xd *xtsDocument) newPagetype(name string, test string) (*pagetype, error) {
	slog.Info("Define new page type", "type", name)
	pt := &pagetype{
		name: name,
		test: test,
	}

	xd.masterpages = append(xd.masterpages, pt)
	return pt, nil
}

func (xd *xtsDocument) detectPagetype() (*pagetype, error) {
	var thispagetype *pagetype
	for i := len(xd.masterpages) - 1; i >= 0; i-- {
		thispagetype = xd.masterpages[i]

		seq, err := evaluateXPath(xd, xd.layoutNS, thispagetype.test)
		var eval, ok bool
		if err == nil && len(seq) != 1 {
			err = fmt.Errorf("the test gives %d values instead of one", len(seq))
		}
		if err == nil {
			if eval, ok = seq[0].(bool); !ok {
				err = fmt.Errorf("the test does not give true or false")
			}
		}
		if err != nil {
			// A test that cannot be evaluated does not match, so the page
			// gets another type. Returning the error left the page nil.
			slog.Error(fmt.Sprintf("page type %q, test %q: %s", thispagetype.name, thispagetype.test, err))
			continue
		}
		if eval {
			break
		}
	}
	slog.Debug("DetectPagetype: chose page type", "name", thispagetype.name)
	return thispagetype, nil

}

type page struct {
	pagenumber    int
	bagPage       *document.Page
	xd            *xtsDocument
	pagetype      *pagetype
	pageWidth     bag.ScaledPoint // total width of the (PDF) page
	pageHeight    bag.ScaledPoint // total height of the (PDF) page
	pagegrid      *grid
	markerid      int
	atPageShipout func()
	// underlay counts the objects the page's AtPageCreation drew: what
	// layer="behind" goes over, and everything placed since goes over it.
	underlay int
}

func clearPage(xd *xtsDocument) {
	if xd.currentPage == nil {
		return
	}
	cp := xd.currentPage
	if cp.atPageShipout != nil {
		cp.atPageShipout()
	}
	// CSS page margin boxes (@top-left, @bottom-center, ...) from the @page
	// rule coupled to this master page. They live in the page margin, so
	// they never touch the grid, and they render at shipout like
	// AtPageShipout, so counter(page) is final.
	if pt := cp.pagetype; pt.cssPage != nil {
		pd := htmlbag.PageDimensions{
			Width:        cp.pageWidth,
			Height:       cp.pageHeight,
			MarginTop:    pt.marginTop,
			MarginBottom: pt.marginBottom,
			MarginLeft:   pt.marginLeft,
			MarginRight:  pt.marginRight,
		}
		if err := xd.cssbuilder.OutputMarginBoxes(pd, pt.cssPage); err != nil {
			slog.Error("Cannot output CSS page margin boxes", "masterpage", pt.name, "page", cp.pagenumber, "error", err)
		}
	}
	xd.currentPage.bagPage.Shipout()
	xd.currentPage = nil
}

// newPage returns the new page object, the AtPageCreation function and an error.
func newPage(xd *xtsDocument) (*page, func(), error) {
	slog.Debug("New page")
	xd.currentPagenumber++
	g := newGrid(xd)
	pt, err := xd.detectPagetype()
	if err != nil {
		return nil, nil, err
	}
	if err = xd.resolvePageCSS(pt); err != nil {
		return nil, nil, err
	}
	d := xd.document.Doc
	g.marginLeft = pt.marginLeft
	g.marginBottom = pt.marginBottom
	g.marginTop = pt.marginTop
	g.marginRight = pt.marginRight

	// Set nx,ny. Either to the default values or to the calculated values.
	gridAreaWidth := d.DefaultPageWidth - g.marginLeft - g.marginRight - bag.ScaledPoint(xd.defaultGridNx-1)*g.gridGapX
	if xd.defaultGridNx > 0 {
		g.gridWidth = gridAreaWidth / bag.ScaledPoint(xd.defaultGridNx)
		g.nx = xd.defaultGridNx
	} else {
		g.nx = int(gridAreaWidth+g.gridGapX) / int(g.gridWidth+g.gridGapX)
	}
	gridAreaHeight := d.DefaultPageHeight - g.marginTop - g.marginBottom - bag.ScaledPoint(xd.defaultGridNy-1)*g.gridGapY
	if xd.defaultGridNy > 0 {
		g.gridHeight = gridAreaHeight / bag.ScaledPoint(xd.defaultGridNy)
		g.ny = xd.defaultGridNy
	} else {
		g.ny = int(gridAreaHeight+g.gridGapY) / int(g.gridHeight+g.gridGapY)
	}

	pg := &page{
		xd:         xd,
		bagPage:    d.NewPage(),
		pagetype:   pt,
		pagegrid:   g,
		pageWidth:  d.DefaultPageWidth,
		pageHeight: d.DefaultPageHeight,
		pagenumber: xd.currentPagenumber,
	}
	g.setPage(pg)

	var atPageCreation func()
	xd.currentGrid = pg.pagegrid
	if pt.layoutElt != nil {
		for _, node := range pt.layoutElt.Children() {
			switch t := node.(type) {
			case *goxml.Element:
				switch t.Name {
				case "AtPageCreation":
					xd.checkAttributes(t, nil)
					slog.Debug(fmt.Sprintf("Call %s (line %d)", t.Name, t.Line))
					atPageCreation = func() { dispatch(xd, t) }
				case "AtPageShipout":
					xd.checkAttributes(t, nil)
					pg.atPageShipout = func() {
						slog.Debug(fmt.Sprintf("Call %s (line %d)", t.Name, t.Line))
						dispatch(xd, t)
					}
				case "PositioningArea":
					attValues := &struct {
						Name string `sdxml:"mustexist"`
					}{}
					if err = getXMLAttributes(xd, t, attValues); err != nil {
						return nil, nil, err
					}
					var rects []*gridRect
					if rects, err = parsePositioningFrames(xd, t); err != nil {
						return nil, nil, err
					}
					if len(rects) == 0 {
						slog.Info(fmt.Sprintf("PositioningArea %s has no frames on page %d, area not created", attValues.Name, xd.currentPagenumber))
						continue
					}
					xd.currentGrid.areas[attValues.Name] = &area{
						name:  attValues.Name,
						frame: rects,
					}
				}
			}
		}

	}
	// CHECK
	docPage := pg.bagPage
	docPage.Userdata = make(map[any]any)
	docPage.Userdata["xtspage"] = pg
	return pg, atPageCreation, nil
}

// parsePositioningFrames parses the children of a PositioningArea (or of a
// matching Switch branch within) into grid rectangles. Switch is evaluated at
// page creation time, so a Case test can use sd:current-page() to select
// different frames for example for even and odd pages.
func parsePositioningFrames(xd *xtsDocument, elt *goxml.Element) ([]*gridRect, error) {
	var rects []*gridRect
	for _, cld := range elt.Children() {
		c, ok := cld.(*goxml.Element)
		if !ok {
			continue
		}
		if c.Name == "Switch" {
			branch, err := selectSwitchBranch(xd, c)
			if err != nil {
				return nil, err
			}
			if branch == nil {
				continue
			}
			branchRects, err := parsePositioningFrames(xd, branch)
			if err != nil {
				return nil, err
			}
			rects = append(rects, branchRects...)
			continue
		}
		attValues := &struct {
			Width  int `sdxml:"mustexist"`
			Height int `sdxml:"mustexist"`
			Column int `sdxml:"mustexist"`
			Row    int `sdxml:"mustexist"`
		}{}
		if err := getXMLAttributes(xd, c, attValues); err != nil {
			return nil, err
		}
		rect := gridRect{
			row:        coord(attValues.Row),
			col:        coord(attValues.Column),
			width:      coord(attValues.Width),
			height:     coord(attValues.Height),
			currentCol: 1,
			currentRow: 1,
		}
		rects = append(rects, &rect)
	}
	return rects, nil
}

// nextMarkerID returns the page's next marker id. It is a counter rather than
// a goroutine feeding a channel, which never ended and kept the page, and
// with it the whole document, alive.
func (p *page) nextMarkerID() int {
	id := p.markerid
	p.markerid++
	return id
}

func (p *page) outputAbsolute(x, y bag.ScaledPoint, vl *node.VList) {
	p.bagPage.OutputAt(x, p.pageHeight-y, vl)
}

// outputBehind places vl under everything on the page but what its
// AtPageCreation drew, in the order such objects are placed.
func (p *page) outputBehind(x, y bag.ScaledPoint, vl *node.VList) {
	objs := p.bagPage.Objects
	at := min(p.underlay, len(objs))
	objs = append(objs, document.Object{})
	copy(objs[at+1:], objs[at:])
	objs[at] = document.Object{X: x, Y: p.pageHeight - y, Vlist: vl}
	p.bagPage.Objects = objs
	p.underlay = at + 1
}

func (p *page) String() string {
	g := p.pagegrid
	return fmt.Sprintf("XTS page %d wd/ht: %s/%s margins: %s %s %s %s", p.pagenumber, p.pageWidth, p.pageHeight, g.marginLeft, g.marginTop, g.marginRight, g.marginBottom)
}

func (xd *xtsDocument) OutputAt(vl *node.VList, col coord, row coord, allocate, behind bool, area *area, what string, halign frontend.HorizontalAlignment) error {
	var currentSlate *slate
	if currentSlate = xd.currentSlate; currentSlate != nil {
		if area.name != pageAreaName {
			slog.Error(fmt.Sprintf("Cannot use area (%s) within a slate (%s)", area.name, currentSlate.name))
		}
		g := xd.currentGrid
		shiftRight := bag.ScaledPoint(0)
		if halign == frontend.HAlignRight {
			f := area.frame[area.currentFrame]
			shiftRight = g.width(f.col+f.width-col) - vl.Width
		}
		// The slate's own coordinate space starts at its top left corner, so
		// the page margins built into posX/posY are removed again.
		x := g.posX(col, area) - g.marginLeft + shiftRight
		y := g.posY(row, area) - g.marginTop
		currentSlate.appendItem(slateItem{x: x, y: y, vl: vl, noRoom: !allocate, behind: behind})
	} else {
		slog.Info("PlaceObject", "obj", what, "col", col, "row", row, "area", area.name)

		shiftRight := bag.ScaledPoint(0)
		if halign == frontend.HAlignRight {
			f := area.frame[area.currentFrame]
			shiftRight = xd.currentGrid.width(f.col+f.width-col) - vl.Width
		}

		columnLength := xd.currentGrid.posX(col, area)
		rowLength := xd.currentGrid.posY(row, area)
		if behind {
			xd.currentPage.outputBehind(columnLength+shiftRight, rowLength, vl)
		} else {
			xd.currentPage.outputAbsolute(columnLength+shiftRight, rowLength, vl)
		}
	}
	if allocate {
		xd.currentGrid.allocate(col, row, area, vl.Width, vl.Height+vl.Depth)
	}
	return nil
}
