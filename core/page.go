package core

import (
	"fmt"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/document"
	"github.com/speedata/boxesandglue/backend/node"
	"github.com/speedata/goxml"
)

const (
	pageAreaName    string = "__defaultarea"
	defaultAreaName string = "__defaultarea"
)

type pagetype struct {
	name         string
	test         string
	areas        map[string]area
	marginLeft   bag.ScaledPoint // the default left margin
	marginRight  bag.ScaledPoint // the default right margin
	marginTop    bag.ScaledPoint // the default top margin
	marginBottom bag.ScaledPoint // the default bottom margin
	layoutElt    *goxml.Element
}

func (xd *xtsDocument) newPagetype(name string, test string) (*pagetype, error) {
	bag.Logger.Infof("Define new page type %q", name)
	pt := &pagetype{
		name: name,
		test: test,
	}

	xd.pagetypes = append(xd.pagetypes, pt)
	return pt, nil
}

func (xd *xtsDocument) detectPagetype(name string) (*pagetype, error) {
	var thispagetype *pagetype
	for i := len(xd.pagetypes) - 1; i >= 0; i-- {
		thispagetype = xd.pagetypes[i]
		seq, err := xd.data.Evaluate(thispagetype.test)
		if err != nil {
			return nil, err
		}
		if len(seq) != 1 {
			return nil, fmt.Errorf("something is wrong with the page type")
		}
		var eval, ok bool
		if eval, ok = seq[0].(bool); !ok {
			return nil, fmt.Errorf("something is wrong with the page type: could not evaluate test to boolean value")
		}
		if eval {
			break
		}
	}
	bag.Logger.Debugf("DetectPagetype: chose page type %q", thispagetype.name)
	return thispagetype, nil

}

type page struct {
	pagenumber int
	bagPage    *document.Page
	xd         *xtsDocument
	pagetype   *pagetype
	pageWidth  bag.ScaledPoint // total width of the (PDF) page
	pageHeight bag.ScaledPoint // total height of the (PDF) page
	pagegrid   *grid
}

func clearPage(xd *xtsDocument) {
	xd.currentPage.bagPage.Shipout()
	xd.currentPage = nil
}

func newPage(xd *xtsDocument) (*page, func(), error) {
	bag.Logger.Debug("newPage")
	g := newGrid(xd)
	pt, err := xd.detectPagetype("")
	if err != nil {
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
	}
	g.setPage(pg)

	var f func()
	xd.currentGrid = pg.pagegrid
	if pt.layoutElt != nil {

		for _, node := range pt.layoutElt.Children() {
			switch t := node.(type) {
			case *goxml.Element:
				switch t.Name {
				case "AtPageCreation":
					bag.Logger.Debugf("Call %s (line %d)", t.Name, t.Line)
					f = func() { dispatch(xd, t, xd.data) }
				case "PositioningArea":
					attValues := &struct {
						Name string `sdxml:"mustexist"`
					}{}
					if err = getXMLAtttributes(xd, t, attValues); err != nil {
						return nil, nil, err
					}
					var rects []*gridRect
					for _, cld := range t.Children() {
						if c, ok := cld.(*goxml.Element); ok {
							attValues := &struct {
								Width  int `sdxml:"mustexist"`
								Height int `sdxml:"mustexist"`
								Column int `sdxml:"mustexist"`
								Row    int `sdxml:"mustexist"`
							}{}
							if err = getXMLAtttributes(xd, c, attValues); err != nil {
								return nil, nil, err
							}
							rect := gridRect{
								row:    coord(attValues.Row),
								col:    coord(attValues.Column),
								width:  coord(attValues.Width),
								height: coord(attValues.Height),
							}
							rects = append(rects, &rect)
						}
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
	docPage.Userdata = make(map[interface{}]interface{})
	docPage.Userdata["xtspage"] = pg
	return pg, f, nil
}

func (p *page) outputAbsolute(x, y bag.ScaledPoint, vl *node.VList) {
	p.bagPage.OutputAt(x, p.pageHeight-y, vl)
}

func (p *page) String() string {
	g := p.pagegrid
	return fmt.Sprintf("XTS page %d wd/ht: %s/%s margins: %s %s %s %s", p.pagenumber, p.pageWidth, p.pageHeight, g.marginLeft, g.marginTop, g.marginRight, g.marginBottom)
}

func (xd *xtsDocument) OutputAt(vl *node.VList, col coord, row coord, allocate bool, areaname string, what string) error {
	areatext := ""
	if areaname != pageAreaName {
		areatext = fmt.Sprintf("%s [%d]: ", areaname, 1)
	}
	var currentGroup *group
	if currentGroup = xd.currentGroup; currentGroup != nil {
		if areaname != pageAreaName {
			bag.Logger.Errorf("Cannot use area (%s) within a group (%s)", areaname, currentGroup.name)
		}
		if currentGroup.contents == nil {
			currentGroup.contents = vl
		}

	} else {
		bag.Logger.Infof("PlaceObject: output %s at (%s%d,%d)", what, areatext, col, row)
		columnLength := xd.currentGrid.posX(col, areaname)
		rowLength := xd.currentGrid.posY(row, areaname)

		xd.currentPage.outputAbsolute(columnLength, rowLength, vl)
	}

	xd.currentGrid.allocate(col, row, areaname, vl.Width, vl.Height)
	return nil
}
