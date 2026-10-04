package core

import (
	"testing"

	"github.com/boxesandglue/boxesandglue/backend/bag"
)

// threeFrames is an area of three frames side by side on a bare grid, with the
// document around it: enough for advanceToFit and nextRow without a page.
func threeFrames(widths ...coord) (*xtsDocument, *area) {
	a := &area{name: "cols"}
	col := coord(1)
	for _, wd := range widths {
		a.frame = append(a.frame, &gridRect{row: 1, col: col, width: wd, height: 6, currentRow: 1, currentCol: 1})
		col += wd + 1
	}
	g := &grid{allocatedBlocks: make(allocationMatrix), areas: map[string]*area{"cols": a}}
	xd := &xtsDocument{currentPage: &page{}, currentGrid: g}
	g.page = xd.currentPage
	xd.currentPage.xd = xd
	return xd, a
}

// fill allocates columns x0..x1 of rows y0..y1.
func fill(g *grid, x0, x1, y0, y1 coord) {
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			g.allocatedBlocks.allocate(x, y)
		}
	}
}

// TestAdvanceToFitPassesOverAFullFrame checks that an object moving on from a
// full frame skips a frame whose free rows are too few for it, instead of
// being put at that frame's row 1 on top of what is already there.
func TestAdvanceToFitPassesOverAFullFrame(t *testing.T) {
	xd, a := threeFrames(5, 5, 5)
	// Something placed earlier covers the top four rows of the middle frame.
	fill(xd.currentGrid, 7, 11, 1, 4)

	got, row := xd.advanceToFit(a, 5, 3, 1)
	if got != a {
		t.Fatalf("got another area %v, want the same one", got)
	}
	if got.currentFrame != 2 || row != 1 {
		t.Errorf("frame %d row %d, want frame 2 row 1: the middle frame has only one free row", got.currentFrame, row)
	}

	a.currentFrame = 0
	got, row = xd.advanceToFit(a, 5, 1, 1)
	if got.currentFrame != 1 || row != 5 {
		t.Errorf("frame %d row %d, want frame 1 row 5, the first free row below what is there", got.currentFrame, row)
	}
}

// TestAdvanceToFitTakesAFreeFrameForAnOversizedObject: an object taller than
// any frame fits none of them, and must still land in the next free frame
// rather than skip them all and leave the page empty.
func TestAdvanceToFitTakesAFreeFrameForAnOversizedObject(t *testing.T) {
	xd, a := threeFrames(5, 5, 5)
	got, row := xd.advanceToFit(a, 5, 15, 1)
	if got.currentFrame != 1 || row != 1 {
		t.Errorf("frame %d row %d, want frame 1 row 1", got.currentFrame, row)
	}
}

// TestNextRowMeasuresEachFrame: moving into a narrower frame is judged by that
// frame's width, not the width of the frame being left.
func TestNextRowMeasuresEachFrame(t *testing.T) {
	xd, a := threeFrames(10, 4, 4)
	// Content placed beside the area, right of the narrow middle frame.
	fill(xd.currentGrid, 17, 21, 1, 6)
	a.SetCurrentRow(6)
	got := xd.currentGrid.nextRow(a)
	if got.currentFrame != 1 || got.CurrentRow() != 1 {
		t.Errorf("frame %d row %d, want frame 1 row 1", got.currentFrame, got.CurrentRow())
	}
}

// TestNextRowInASlate: a slate has no page to break to; running out of rows
// there overflows as before instead of reaching for a page.
func TestNextRowInASlate(t *testing.T) {
	a := &area{name: "cols", frame: []*gridRect{{row: 1, col: 1, width: 5, height: 3, currentRow: 3, currentCol: 1}}}
	g := &grid{allocatedBlocks: make(allocationMatrix), areas: map[string]*area{"cols": a}, inSlate: true}
	if got := g.nextRow(a); got != a || got.CurrentRow() != 1 {
		t.Errorf("area %v row %d, want the same area at row 1", got, got.CurrentRow())
	}
}

// TestFindSuitableRowTakesTheLastRow: an object one row high fits into the
// last row of a frame, and one that would reach below the frame does not.
func TestFindSuitableRowTakesTheLastRow(t *testing.T) {
	xd, a := threeFrames(5)
	g := xd.currentGrid
	a.SetCurrentRow(6)
	if row := g.findSuitableRow(5, 1, 1, a); row != 6 {
		t.Errorf("one row from row 6 of 6: got row %d, want 6", row)
	}
	if row := g.findSuitableRow(5, 2, 1, a); row != -1 {
		t.Errorf("two rows from row 6 of 6: got row %d, want -1", row)
	}
	a.SetCurrentRow(1)
	fill(g, 1, 5, 1, 5)
	if row := g.findSuitableRow(5, 1, 1, a); row != 6 {
		t.Errorf("one row below five full ones: got row %d, want 6", row)
	}
}

// TestFitsInRowChecksEveryColumn: an object fits into a row only when none of
// the columns it covers is allocated, its last column and a single column
// included.
func TestFitsInRowChecksEveryColumn(t *testing.T) {
	xd, a := threeFrames(5)
	g := xd.currentGrid
	fill(g, 4, 4, 1, 1)
	if g.fitsInRow(1, 1, 4, a) {
		t.Error("columns 1 to 4 fit although column 4 is allocated")
	}
	if !g.fitsInRow(1, 1, 3, a) {
		t.Error("columns 1 to 3 do not fit although they are free")
	}
	if g.fitsInRow(4, 1, 1, a) {
		t.Error("column 4 alone fits although it is allocated")
	}
	if !g.fitsInRow(5, 1, 1, a) {
		t.Error("column 5 alone does not fit although it is free")
	}
}

// TestFrameOutlinesEachAtItsPlace: the grid trace draws every frame of an
// area where it is, not every one at the place of the area's current frame,
// and in the order of the area names, so that the output does not depend on
// the order of a map.
func TestFrameOutlinesEachAtItsPlace(t *testing.T) {
	xd, a := threeFrames(5, 5)
	g := xd.currentGrid
	g.gridWidth, g.gridHeight = 10, 20
	g.areas["a"] = &area{name: "a", frame: []*gridRect{{row: 3, col: 2, width: 1, height: 1}}}
	a.currentFrame = 1
	got := g.frameOutlines()
	want := []frameOutline{
		{left: 10, top: 40, width: 10, height: 20},
		{left: 0, top: 0, width: 50, height: 120},
		{left: 60, top: 0, width: 50, height: 120},
	}
	if len(got) != len(want) {
		t.Fatalf("got %d outlines, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("outline %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// TestPositionsCountTheGaps: with gaps between the cells, column c starts
// c-1 cells and c-1 gaps right of the margin, row r as many below it, also
// in a frame that does not start in the first column or row.
func TestPositionsCountTheGaps(t *testing.T) {
	g := &grid{marginLeft: 100, marginTop: 200, gridWidth: 10, gridHeight: 20, gridGapX: 5, gridGapY: 2}
	page := &gridRect{col: 1, row: 1, width: 10, height: 10}
	frame := &gridRect{col: 3, row: 4, width: 2, height: 2}
	for _, tc := range []struct {
		f        *gridRect
		col, row coord
		x, y     bag.ScaledPoint
	}{
		{page, 1, 1, 100, 200},
		{page, 2, 2, 115, 222},
		{page, 3, 3, 130, 244},
		{frame, 1, 1, 130, 266},
		{frame, 2, 2, 145, 288},
	} {
		if x := g.frameX(tc.col, tc.f); x != tc.x {
			t.Errorf("column %d of frame at column %d: x = %d, want %d", tc.col, tc.f.col, x, tc.x)
		}
		if y := g.frameY(tc.row, tc.f); y != tc.y {
			t.Errorf("row %d of frame at row %d: y = %d, want %d", tc.row, tc.f.row, y, tc.y)
		}
	}
}
