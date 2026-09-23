package core

import "testing"

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

	got, row := xd.advanceToFit(a, 5, 3, 3, 1)
	if got != a {
		t.Fatalf("got another area %v, want the same one", got)
	}
	if got.currentFrame != 2 || row != 1 {
		t.Errorf("frame %d row %d, want frame 2 row 1: the middle frame has only one free row", got.currentFrame, row)
	}

	a.currentFrame = 0
	got, row = xd.advanceToFit(a, 5, 1, 1, 1)
	if got.currentFrame != 1 || row != 5 {
		t.Errorf("frame %d row %d, want frame 1 row 5, the first free row below what is there", got.currentFrame, row)
	}

	// A table the splitter continues needs one free row to start in.
	a.currentFrame = 0
	got, row = xd.advanceToFit(a, 5, 3, 1, 1)
	if got.currentFrame != 1 || row != 5 {
		t.Errorf("splittable: frame %d row %d, want frame 1 row 5", got.currentFrame, row)
	}
}

// TestAdvanceToFitTakesAFreeFrameForAnOversizedObject: an object taller than
// any frame fits none of them, and must still land in the next free frame
// rather than skip them all and leave the page empty.
func TestAdvanceToFitTakesAFreeFrameForAnOversizedObject(t *testing.T) {
	xd, a := threeFrames(5, 5, 5)
	got, row := xd.advanceToFit(a, 5, 15, 15, 1)
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
