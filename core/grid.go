package core

import (
	"fmt"
	"math"
	"strings"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/node"
)

type area struct {
	currentFrame int
	name         string
	frame        []*gridRect
}

func (area area) String() string {
	var ret []string
	for _, f := range area.frame {
		ret = append(ret, f.String())
	}
	return fmt.Sprintf("%s: %s", area.name, strings.Join(ret, "|"))
}

// CurrentRow returns the current row of the current active frame.
func (area *area) CurrentRow() coord {
	return area.frame[area.currentFrame].row
}

// CurrentCol returns the current column of the current active frame.
func (area *area) CurrentCol() coord {
	return area.frame[area.currentFrame].col
}

// SetCurrentRow sets the current row in the active frame in the area.
func (area *area) SetCurrentRow(row coord) {
	area.frame[area.currentFrame].row = row
}

// SetCurrentCol sets the current column in the active frame in the area.
func (area *area) SetCurrentCol(col coord) {
	area.frame[area.currentFrame].col = col
}

type allocationMatrix map[gridCoord]int

type gridRect struct {
	row    coord
	col    coord
	width  coord
	height coord
}

func (gr *gridRect) String() string {
	return fmt.Sprintf("%d/%d wd: %d ht: %d", gr.col, gr.row, gr.width, gr.height)
}

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

type gridCoord int64
type coord int32

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
	marginLeft      bag.ScaledPoint
	marginRight     bag.ScaledPoint
	marginTop       bag.ScaledPoint
	marginBottom    bag.ScaledPoint
	allocatedBlocks allocationMatrix
	areas           map[string]*area
	inGroup         bool
}

func newGrid(xd *xtsDocument) *grid {
	g := &grid{
		gridWidth:  xd.defaultGridWidth,
		gridHeight: xd.defaultGridHeight,
		gridGapX:   xd.defaultGridGapX,
		gridGapY:   xd.defaultGridGapY,
		areas:      make(map[string]*area),
		inGroup:    true,
	}

	return g
}

// Connect the grid to a page and initialize the allocation matrix.
func (g *grid) setPage(p *page) {
	g.page = p
	g.allocatedBlocks = make(allocationMatrix)
	g.areas[pageAreaName] = &area{
		name:  pageAreaName,
		frame: []*gridRect{{1, 1, coord(g.nx), coord(g.ny)}},
	}
}

func (g *grid) String() string {
	return fmt.Sprintf("grid %dx%d", g.nx, g.ny)
}

// posX returns the horizontal offset relative to the left page border. Column 1
// returns the margin left.
func (g *grid) posX(column coord, area *area) bag.ScaledPoint {
	if column == 0 {
		return bag.ScaledPoint(0)
	}
	posx := g.marginLeft + bag.ScaledPoint(column-1)*g.gridWidth
	if column > 1 {
		posx += bag.ScaledPoint(column-2) * g.gridGapX
	}
	return posx
}

// posY returns the vertical offset relative to the top page border. Row 1
// returns the top margin.
func (g *grid) posY(row coord, area *area) bag.ScaledPoint {
	if row == 0 {
		return bag.ScaledPoint(0)
	}
	posy := g.marginTop + bag.ScaledPoint(row-1)*g.gridHeight
	if row > 1 {
		posy += bag.ScaledPoint(row-2) * g.gridGapY
	}
	return posy
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

func (g *grid) allocate(x, y coord, area *area, wd, ht bag.ScaledPoint) {
	var warningTopRaised, warningLeftRaised, warningRightRaised, warningBottomRaised bool
	var offsetX coord
	var offsetY coord

	offsetX = area.CurrentCol()
	offsetY = area.CurrentRow()

	for col := coord(1); col <= g.widthToColumns(wd); col++ {
		for row := coord(1); row <= g.heightToRows(ht); row++ {
			if posX, posY := col+x+offsetX-2, row+y+offsetY-2; posX >= 1 && posY >= 1 && posX <= coord(g.nx) && posY <= coord(g.ny) {
				g.allocatedBlocks.allocate(posX, posY)
			} else {
				if posX < 1 && !warningLeftRaised && !g.inGroup {
					bag.Logger.Warn("object protrudes into the left margin")
					warningLeftRaised = true
				}
				if posY < 1 && !warningTopRaised && !g.inGroup {
					bag.Logger.Warn("object protrudes into the top margin")
					warningTopRaised = true
				}
				if posX > coord(g.nx) && !warningRightRaised && !g.inGroup {
					bag.Logger.Warn("object protrudes into the right margin")
					warningRightRaised = true
				}
				if posY > coord(g.ny) && !warningBottomRaised && !g.inGroup {
					bag.Logger.Warn("object protrudes into the bottom margin")
					warningBottomRaised = true
				}
			}
		}
	}
	col := x + g.widthToColumns(wd)
	if col > coord(g.nx) {
		area.SetCurrentCol(1)
		area.SetCurrentRow(y + g.heightToRows(ht))
	}
}
func (g *grid) findFreeSpaceForObject(vl *node.VList, area *area) (gridCoord, error) {
	if area.CurrentRow() == 0 {
		area.SetCurrentRow(1)
	}
	if area.CurrentCol() == 0 {
		area.SetCurrentCol(1)
	}
	// this doesn't make sense, check
	rowOffset := area.CurrentRow()

	row := area.CurrentRow() + rowOffset - 1
	wdCols := g.widthToColumns(vl.Width)

	// if g.currentCol >= coord(g.nx) {
	// 	g.nextRow()
	// }
	col := g.fitsInRow(row, wdCols, area)
	if col > 0 {
		col = col - area.CurrentCol() + 1
	}
	xy := newGridCoord(col, row-rowOffset+1)

	return xy, nil
}

func (g *grid) nextRow() {
	// g.currentCol = 1
	// g.currentRow++
}

func (g *grid) fitsInRow(y coord, wdCols coord, area *area) coord {
	col := area.CurrentCol()
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
