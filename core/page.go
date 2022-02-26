package core

import (
	"fmt"
	"math"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/document"
	"github.com/speedata/boxesandglue/backend/node"
)

type allocationMatrix map[gridCoord]int

func (am allocationMatrix) allocate(x, y coord) {
	xy := newGridCoord(x, y)
	am[xy]++
}

func (am allocationMatrix) allocValue(x, y coord) int {
	return am[newGridCoord(x, y)]
}

type gridCoord uint64
type coord uint32

func (c coord) String() string {
	return fmt.Sprintf("%d", c)
}

// newGridCoord creates a grid coordinate from the x and y values where (x,y) is
// the upper left. To be used in the allocation matrix.
func newGridCoord(x, y coord) gridCoord {
	return gridCoord(x)<<32 + gridCoord(y)
}

// XY returns the x and the y coordinate.
func (gc gridCoord) XY() (coord, coord) {
	return coord(gc >> 32), coord(gc & 0xffffffff)
}

func (gc gridCoord) GoString() string {
	x, y := gc.XY()
	return fmt.Sprintf("(%d,%d)", x, y)
}

func (gc gridCoord) String() string {
	x, y := gc.XY()
	return fmt.Sprintf("(%d,%d)", x, y)
}

type grid struct {
	page            *page
	gridWidth       bag.ScaledPoint
	gridHeight      bag.ScaledPoint
	gridGapX        bag.ScaledPoint
	gridGapY        bag.ScaledPoint
	nx              int
	ny              int
	allocatedBlocks allocationMatrix
	currentCol      coord
	currentRow      coord
}

func newGrid(xd *xtsDocument) *grid {
	g := &grid{
		gridWidth:  xd.defaultGridWidth,
		gridHeight: xd.defaultGridHeight,
		gridGapX:   xd.defaultGridGapX,
		gridGapY:   xd.defaultGridGapY,
		currentRow: 1,
		currentCol: 1,
	}

	return g
}

// Connect the grid to a page and initialize the allocation matrix.
func (g *grid) setPage(p *page) {
	g.page = p
	pageAreaX := p.pageWidth - p.marginLeft - p.marginRight
	pageAreaY := p.pageHeight - p.marginTop - p.marginBottom
	// there might be more cells to the right and to the bottom, but those are
	// omitted, because they are only partly visible.
	if g.nx == 0 {
		g.nx = int(pageAreaX / (g.gridWidth + g.gridGapX))
	}
	if g.ny == 0 {
		g.ny = int(pageAreaY / (g.gridHeight + g.gridGapY))
	}
	g.allocatedBlocks = make(allocationMatrix)
}

func (g *grid) String() string {
	return fmt.Sprintf("grid %dx%d", g.nx, g.ny)
}

// posX returns the horizontal offset relative to the left page border. Column 1 returns the margin left.
func (g *grid) posX(column coord) bag.ScaledPoint {
	if column == 0 {
		return bag.ScaledPoint(0)
	}
	return bag.ScaledPoint(column-1)*g.gridWidth + bag.ScaledPoint(column-1)*g.gridGapX + g.page.marginLeft
}

// posY returns the vertical offset relative to the top page border. Row 1 returns the top margin.
func (g *grid) posY(row coord) bag.ScaledPoint {
	if row == 0 {
		return bag.ScaledPoint(0)
	}
	return bag.ScaledPoint(row-1)*g.gridHeight + bag.ScaledPoint(row-1)*g.gridGapY + g.page.marginTop
}

// height returns the height of the number of columns.
func (g *grid) height(columns coord) bag.ScaledPoint {
	return bag.ScaledPoint(columns)*g.gridHeight + bag.ScaledPoint(columns-1)*g.gridGapX
}

// width returns the width of the number of columns.
func (g *grid) width(columns coord) bag.ScaledPoint {
	return bag.ScaledPoint(columns)*g.gridWidth + bag.ScaledPoint(columns-1)*g.gridGapX
}

func (g *grid) widthToColumns(width bag.ScaledPoint) coord {
	r := float64(width) / float64(g.gridWidth+g.gridGapX)
	return coord(math.Ceil(r - 0.005))
}

func (g *grid) heightToRows(height bag.ScaledPoint) coord {
	r := float64(height) / float64(g.gridHeight+g.gridGapY)
	return coord(math.Ceil(r - 0.005))
}

func (g *grid) allocate(x, y coord, wd, ht bag.ScaledPoint) {
	var warningTopRaised, warningLeftRaised, warningRightRaised, warningBottomRaised bool
	for col := coord(1); col <= g.widthToColumns(wd); col++ {
		for row := coord(1); row <= g.heightToRows(ht); row++ {
			if posX, posY := col+x-1, row+y-1; posX >= 1 && posY >= 1 && posX <= coord(g.nx) && posY <= coord(g.ny) {
				g.allocatedBlocks.allocate(posX, posY)
			} else {
				if posX < 1 && !warningLeftRaised {
					bag.Logger.Warn("object protrudes into the left margin")
					warningLeftRaised = true
				}
				if posY < 1 && !warningTopRaised {
					bag.Logger.Warn("object protrudes into the top margin")
					warningTopRaised = true
				}
				if posX > coord(g.nx) && !warningRightRaised {
					bag.Logger.Warn("object protrudes into the right margin")
					warningRightRaised = true
				}
				if posY > coord(g.ny) && !warningBottomRaised {
					bag.Logger.Warn("object protrudes into the bottom margin")
					warningBottomRaised = true
				}
			}
		}
	}
	g.currentCol = g.widthToColumns(wd) + x
	g.currentRow = y
}

func (g *grid) fitsInRow(y coord, wdCols coord) coord {
	col := g.currentCol
	row := y
	for {
		if g.allocatedBlocks.allocValue(col, row) > 0 || int(col) > g.nx {
			col++
		} else {
			break
		}
	}
	nowhere := newGridCoord(0, 0)
	for i := col; i < col+wdCols; i++ {
		if g.allocatedBlocks.allocValue(i, row) > 0 {
			return coord(nowhere)
		}
	}
	return col
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
			bag.Logger.Debugf("DetectPagetype: chose page type %q", thispagetype.name)
			return thispagetype, nil
		}
	}
	return nil, nil
}

type page struct {
	pagenumber   int
	bagPage      *document.Page
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
	bag.Logger.Debug("newPage")
	g := newGrid(xd)
	pt, err := xd.detectPagetype("")
	if err != nil {
		return nil, err
	}
	d := xd.frontend.Doc

	if xd.defaultGridNx > 0 {
		gridAreaWidth := d.DefaultPageWidth - pt.marginLeft - pt.marginRight - bag.ScaledPoint(xd.defaultGridNx-1)*g.gridGapX
		g.gridWidth = gridAreaWidth / bag.ScaledPoint(xd.defaultGridNx)
		g.nx = xd.defaultGridNx
	}
	if xd.defaultGridNy > 0 {
		gridAreaHeight := d.DefaultPageHeight - pt.marginTop - pt.marginBottom - bag.ScaledPoint(xd.defaultGridNy-1)*g.gridGapY
		g.gridHeight = gridAreaHeight / bag.ScaledPoint(xd.defaultGridNy)
		g.ny = xd.defaultGridNy
	}

	pg := &page{
		xd:           xd,
		bagPage:      d.NewPage(),
		pagetype:     pt,
		pagegrid:     g,
		pageWidth:    d.DefaultPageWidth,
		pageHeight:   d.DefaultPageHeight,
		marginLeft:   pt.marginLeft,
		marginRight:  pt.marginRight,
		marginTop:    pt.marginTop,
		marginBottom: pt.marginBottom,
	}
	g.setPage(pg)
	xd.currentGrid = pg.pagegrid
	// CHECK
	docPage := pg.bagPage
	docPage.Userdata = make(map[interface{}]interface{})
	docPage.Userdata["xtspage"] = pg
	return pg, nil
}

func (p *page) outputAbsolute(x, y bag.ScaledPoint, vl *node.VList) {
	p.bagPage.OutputAt(x, p.pageHeight-y, vl)
}

func (p *page) findFreeSpaceForObject(vl *node.VList) gridCoord {
	currentGrid := p.pagegrid
	// htRows := currentGrid.heightToRows(vl.Height + vl.Depth)
	wdCols := currentGrid.widthToColumns(vl.Width)
	col := currentGrid.fitsInRow(currentGrid.currentRow, wdCols)
	xy := newGridCoord(col, 1)
	return xy
}

func (p *page) String() string {
	return fmt.Sprintf("XTS page %d wd/ht: %s/%s margins: %s %s %s %s", p.pagenumber, p.pageWidth, p.pageHeight, p.marginLeft, p.marginTop, p.marginRight, p.marginBottom)
}
