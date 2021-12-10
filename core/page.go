package core

import (
	"fmt"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/node"
)

type grid struct {
	page       *page
	gridWidth  bag.ScaledPoint
	gridHeight bag.ScaledPoint
	gridGapX   bag.ScaledPoint
	gridGapY   bag.ScaledPoint
}

func (g *grid) posX(column int) bag.ScaledPoint {
	return bag.ScaledPoint(column)*g.gridWidth + bag.ScaledPoint(column-1)*g.gridGapX + g.page.marginLeft
}

func (g *grid) posY(row int) bag.ScaledPoint {
	return bag.ScaledPoint(row)*g.gridHeight + bag.ScaledPoint(row-1)*g.gridGapY + g.page.marginTop
}

func (g *grid) width(columns int) bag.ScaledPoint {
	return bag.ScaledPoint(columns)*g.gridWidth + bag.ScaledPoint(columns-1)*g.gridGapX
}

type pagetype struct {
	name         string
	test         string
	marginLeft   bag.ScaledPoint
	marginRight  bag.ScaledPoint
	marginTop    bag.ScaledPoint
	marginBottom bag.ScaledPoint
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
	for i := len(xd.pagetypes) - 1; i >= 0; i-- {
		thispagetype := xd.pagetypes[i]
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
			bag.Logger.Debugf("Page type %q chosen", thispagetype.name)
			return thispagetype, nil
		}
	}
	return nil, nil
}

type page struct {
	pagenumber   int
	xd           *xtsDocument
	pagetype     *pagetype
	pageWidth    bag.ScaledPoint
	pageHeight   bag.ScaledPoint
	marginLeft   bag.ScaledPoint
	marginRight  bag.ScaledPoint
	marginTop    bag.ScaledPoint
	marginBottom bag.ScaledPoint
	pagegrid     *grid
}

func newPage(xd *xtsDocument) (*page, error) {
	g := &grid{
		gridWidth:  xd.defaultGridWidth,
		gridHeight: xd.defaultGridHeight,
		gridGapX:   xd.defaultGridGapX,
		gridGapY:   xd.defaultGridGapY,
	}
	pt, err := xd.detectPagetype("")
	if err != nil {
		return nil, err
	}
	if xd.defaultGridNx > 0 {
		gridAreaWidth := xd.defaultPageWidth - pt.marginLeft - pt.marginRight - bag.ScaledPoint(xd.defaultGridNx-1)*g.gridGapX
		g.gridWidth = gridAreaWidth / bag.ScaledPoint(xd.defaultGridNx)
	}
	if xd.defaultGridNy > 0 {
		gridAreaHeight := xd.defaultPageHeight - pt.marginTop - pt.marginBottom - bag.ScaledPoint(xd.defaultGridNy-1)*g.gridGapY
		g.gridHeight = gridAreaHeight / bag.ScaledPoint(xd.defaultGridNy)
	}

	pg := &page{
		xd:           xd,
		pagetype:     pt,
		pagegrid:     g,
		pageWidth:    xd.defaultPageWidth,
		pageHeight:   xd.defaultPageHeight,
		marginLeft:   pt.marginLeft,
		marginRight:  pt.marginRight,
		marginTop:    pt.marginTop,
		marginBottom: pt.marginBottom,
	}
	g.page = pg
	xd.currentGrid = pg.pagegrid
	docPage := xd.doc.NewPage()
	docPage.Userdata = make(map[interface{}]interface{})
	docPage.Userdata["xtspage"] = pg
	return pg, nil
}

func (p *page) outputAbsolute(x, y bag.ScaledPoint, vl *node.VList) {
	p.xd.doc.OutputAt(x, p.pageHeight-y-vl.Height, vl)
}

func (p *page) String() string {
	return fmt.Sprintf("XTS page %d wd/ht: %s/%s margins: %s %s %s %s", p.pagenumber, p.pageWidth, p.pageHeight, p.marginLeft, p.marginTop, p.marginRight, p.marginBottom)
}
