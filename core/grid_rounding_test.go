package core

import (
	"testing"

	"github.com/boxesandglue/boxesandglue/backend/bag"
)

// An object's height takes the next whole number of rows, as it always has;
// with rounding="nearest" it takes the nearest.
func TestHeightToRowsRounding(t *testing.T) {
	g := &grid{gridHeight: bag.Factor / 4}
	for _, tc := range []struct {
		height      bag.ScaledPoint
		up, nearest coord
	}{
		{bag.Factor * 143 / 10, 58, 57},
		{bag.Factor * 147 / 10, 59, 59},
		{bag.Factor * 10, 40, 40},
		{bag.Factor*10 + 1, 40, 40},
		{bag.Factor * 1001 / 100, 41, 40},
	} {
		g.nearest = false
		if got := g.heightToRows(tc.height); got != tc.up {
			t.Errorf("%s up: %d rows, want %d", tc.height, got, tc.up)
		}
		g.nearest = true
		if got := g.heightToRows(tc.height); got != tc.nearest {
			t.Errorf("%s nearest: %d rows, want %d", tc.height, got, tc.nearest)
		}
	}
}

func fineGrid() (*grid, *area) {
	g := &grid{gridWidth: bag.MustSP("10pt"), gridHeight: bag.MustSP("0.25pt"), nx: 10, ny: 4000, allocatedBlocks: make(allocationMatrix), areas: map[string]*area{}, nearest: true}
	g.setPage(&page{})
	return g, g.areas[pageAreaName]
}

// Twenty paragraphs of a height that is not a whole number of fine rows end
// where their exact heights add up to, within half a row, instead of each one
// rounding the same way.
func TestFlowRowsCarryTheRounding(t *testing.T) {
	g, a := fineGrid()
	ht := bag.MustSP("14.63pt")
	for i := 0; i < 20; i++ {
		g.allocate(1, a.CurrentRow(), a, bag.MustSP("100pt"), ht)
	}
	got := bag.ScaledPoint(a.CurrentRow()-1) * g.gridHeight
	if d := (got - 20*ht).ToPT(); d > 0.125 || d < -0.125 {
		t.Errorf("twenty 14.63pt objects end at %s, %.3fpt from their exact height", got, d)
	}
}

// NextRow moves the flow by whole rows and keeps its carry.
func TestNextRowKeepsTheCarry(t *testing.T) {
	g, a := fineGrid()
	ht := bag.MustSP("14.63pt")
	for i := 0; i < 10; i++ {
		g.allocate(1, a.CurrentRow(), a, bag.MustSP("100pt"), ht)
		from := a.CurrentRow()
		a.SetCurrentRow(from + 32)
		keepCarry(a, 0, from)
	}
	got := bag.ScaledPoint(a.CurrentRow()-1) * g.gridHeight
	if d := (got - 10*(ht+bag.MustSP("8pt"))).ToPT(); d > 0.125 || d < -0.125 {
		t.Errorf("ten objects with 8pt between end %.3fpt from their exact height", d)
	}
}

// Rounding up, as by default, each object takes the next whole row and
// nothing is carried.
func TestFlowRowsRoundUpByDefault(t *testing.T) {
	g, a := fineGrid()
	g.nearest = false
	ht := bag.MustSP("14.63pt")
	for i := 0; i < 20; i++ {
		g.allocate(1, a.CurrentRow(), a, bag.MustSP("100pt"), ht)
	}
	if got, want := a.CurrentRow(), coord(1+20*59); got != want {
		t.Errorf("twenty 14.63pt objects end at row %d, want %d (59 rows each)", got, want)
	}
}

// SetGrid rounding="nearest" applies to the grids made after it, a SetGrid
// without the attribute keeps it, and any other value is an error.
func TestSetGridRounding(t *testing.T) {
	xd := &xtsDocument{}
	for _, tc := range []struct {
		snippet string
		nearest bool
		err     bool
	}{
		{`<SetGrid rounding="nearest"/>`, true, false},
		{`<SetGrid height="12pt"/>`, true, false},
		{`<SetGrid rounding="up"/>`, false, false},
		{`<SetGrid rounding="down"/>`, false, true},
	} {
		_, err := cmdSetGrid(xd, parseLayoutElt(t, tc.snippet))
		if (err != nil) != tc.err {
			t.Errorf("%s: error %v", tc.snippet, err)
		}
		if xd.defaultGridNearest != tc.nearest {
			t.Errorf("%s: nearest = %t, want %t", tc.snippet, xd.defaultGridNearest, tc.nearest)
		}
		if g := newGrid(xd); g.nearest != tc.nearest {
			t.Errorf("%s: new grid nearest = %t, want %t", tc.snippet, g.nearest, tc.nearest)
		}
	}
}
