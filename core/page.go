package core

import (
	"fmt"
	"math"
	"strings"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/document"
	"github.com/speedata/boxesandglue/backend/node"
	"github.com/speedata/goxml"
)

const (
	pageAreaName    string = "__defaultarea"
	defaultAreaName string = "__defaultarea"
)

type allocationMatrix map[gridCoord]int

func (am allocationMatrix) allocate(x, y coord) {
	xy := newGridCoord(x, y)
	am[xy]++
}

func (am allocationMatrix) allocValue(x, y coord) int {
	return am[newGridCoord(x, y)]
}

func (am allocationMatrix) String() string {
	var maxX, maxY coord
	for k := range am {
		x, y := k.XY()
		if x > maxX {
			maxX = x
		}
		if y > maxY {
			maxY = y
		}
	}

	var ret strings.Builder
	for y := coord(1); y <= maxY; y++ {
		for x := coord(1); x <= maxX; x++ {
			if am.allocValue(x, y) > 0 {
				ret.WriteRune('▊')
			} else {
				ret.WriteRune('.')
			}
		}
		ret.WriteString("\n")
	}
	return ret.String()
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
	gridWidth       bag.ScaledPoint // width of the grid cells
	gridHeight      bag.ScaledPoint // height of the grid cells
	gridGapX        bag.ScaledPoint // horizontal space between two grid cells
	gridGapY        bag.ScaledPoint // vertical space between two grid cells
	nx              int             // nx * grid width + ( nx - 1 ) * grid gap + margin = paper width
	ny              int             // ny * grid height + ( ny - 1 ) * grid gap + margin = paper height
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
	g.allocatedBlocks = make(allocationMatrix)
}

func (g *grid) String() string {
	return fmt.Sprintf("grid %dx%d", g.nx, g.ny)
}

// posX returns the horizontal offset relative to the left page border. Column 1
// returns the margin left.
func (g *grid) posX(column coord, areaname string) bag.ScaledPoint {
	if column == 0 {
		return bag.ScaledPoint(0)
	}
	var offsetX coord
	if area, ok := g.page.areas[areaname]; ok {
		offsetX = area.frame[area.currentFrame].col
	}
	posX := bag.ScaledPoint(column + offsetX - 2)
	return posX*g.gridWidth + posX*g.gridGapX + g.page.pagetype.marginLeft
}

// posY returns the vertical offset relative to the top page border. Row 1
// returns the top margin.
func (g *grid) posY(row coord, areaname string) bag.ScaledPoint {
	if row == 0 {
		return bag.ScaledPoint(0)
	}
	var offsetY coord
	if area, ok := g.page.areas[areaname]; ok {
		offsetY = area.frame[area.currentFrame].row
	}
	posY := bag.ScaledPoint(row + offsetY - 2)
	return posY*g.gridHeight + posY*g.gridGapY + g.page.pagetype.marginTop
}

// height returns the height of the number of columns.
func (g *grid) height(columns coord) bag.ScaledPoint {
	return bag.ScaledPoint(columns)*g.gridHeight + bag.ScaledPoint(columns-1)*g.gridGapY
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

func (g *grid) allocate(x, y coord, areaname string, wd, ht bag.ScaledPoint) {
	var warningTopRaised, warningLeftRaised, warningRightRaised, warningBottomRaised bool
	var offsetX coord
	var offsetY coord
	if area, ok := g.page.areas[areaname]; ok {
		offsetX = area.frame[area.currentFrame].col
	}
	if area, ok := g.page.areas[areaname]; ok {
		offsetY = area.frame[area.currentFrame].row
	}

	for col := coord(1); col <= g.widthToColumns(wd); col++ {
		for row := coord(1); row <= g.heightToRows(ht); row++ {
			if posX, posY := col+x+offsetX-2, row+y+offsetY-2; posX >= 1 && posY >= 1 && posX <= coord(g.nx) && posY <= coord(g.ny) {
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

func (g *grid) nextRow() {
	g.currentCol = 1
	g.currentRow++
}

func (g *grid) fitsInRow(y coord, wdCols coord, area *area) coord {
	col := area.currentCol + area.frame[area.currentFrame].col
	row := y
	for {
		if g.allocatedBlocks.allocValue(col, row) > 0 && int(col) <= g.nx {
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
	areas        map[string]area
	marginLeft   bag.ScaledPoint
	marginRight  bag.ScaledPoint
	marginTop    bag.ScaledPoint
	marginBottom bag.ScaledPoint
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
	areas      map[string]*area
	pagegrid   *grid
}

type gridRect struct {
	row    coord
	col    coord
	width  coord
	height coord
}

func (gr *gridRect) String() string {
	return fmt.Sprintf("%d/%d wd: %d ht: %d", gr.col, gr.row, gr.width, gr.height)
}

type area struct {
	currentRow   coord
	currentCol   coord
	currentFrame int
	name         string
	frame        []*gridRect
}

func (a area) String() string {
	var ret []string
	for _, f := range a.frame {
		ret = append(ret, f.String())
	}
	return fmt.Sprintf("%s: %s", a.name, strings.Join(ret, "|"))
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

	// Set nx,ny. Either to the default values or to the calculated values.
	gridAreaWidth := d.DefaultPageWidth - pt.marginLeft - pt.marginRight - bag.ScaledPoint(xd.defaultGridNx-1)*g.gridGapX
	if xd.defaultGridNx > 0 {
		g.gridWidth = gridAreaWidth / bag.ScaledPoint(xd.defaultGridNx)
		g.nx = xd.defaultGridNx
	} else {
		g.nx = int(gridAreaWidth+g.gridGapX) / int(g.gridWidth+g.gridGapX)
	}
	gridAreaHeight := d.DefaultPageHeight - pt.marginTop - pt.marginBottom - bag.ScaledPoint(xd.defaultGridNy-1)*g.gridGapY
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
		areas:      make(map[string]*area),
	}
	g.setPage(pg)

	pg.areas[pageAreaName] = &area{
		name:  pageAreaName,
		frame: []*gridRect{{1, 1, coord(g.nx), coord(g.ny)}},
	}

	var f func()
	xd.currentGrid = pg.pagegrid
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
				pg.areas[attValues.Name] = &area{
					name:  attValues.Name,
					frame: rects,
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

func (p *page) findFreeSpaceForObject(vl *node.VList, areaname string) (gridCoord, error) {
	var area *area
	var ok bool
	if area, ok = p.areas[areaname]; !ok {
		return 0, fmt.Errorf("area %q not found", areaname)
	}
	if area.currentRow == 0 {
		area.currentRow = 1
	}
	rowOffset := area.frame[area.currentFrame].row

	row := area.currentRow + rowOffset - 1
	currentGrid := p.pagegrid
	// htRows := currentGrid.heightToRows(vl.Height + vl.Depth)
	wdCols := currentGrid.widthToColumns(vl.Width)
	if p.pagegrid.currentCol >= coord(p.pagegrid.nx) {
		p.pagegrid.nextRow()
	}
	col := currentGrid.fitsInRow(row, wdCols, area)
	if col > 0 {
		col = col - area.frame[area.currentFrame].col + 1
	}
	xy := newGridCoord(col, row-rowOffset+1)

	return xy, nil
}

func (p *page) String() string {
	return fmt.Sprintf("XTS page %d wd/ht: %s/%s margins: %s %s %s %s", p.pagenumber, p.pageWidth, p.pageHeight, p.pagetype.marginLeft, p.pagetype.marginTop, p.pagetype.marginRight, p.pagetype.marginBottom)
}

func (xd *xtsDocument) OutputAt(vl *node.VList, col coord, row coord, allocate bool, areaname string, what string) error {
	areatext := ""
	if areaname != pageAreaName {
		areatext = fmt.Sprintf("%s [%d]: ", areaname, 1)
	}
	bag.Logger.Infof("PlaceObject: output %s at (%s%d,%d)", what, areatext, col, row)
	columnLength := xd.currentGrid.posX(col, areaname)
	rowLength := xd.currentGrid.posY(row, areaname)
	xd.currentPage.outputAbsolute(columnLength, rowLength, vl)

	xd.currentGrid.allocate(col, row, areaname, vl.Width, vl.Height)
	return nil
}
