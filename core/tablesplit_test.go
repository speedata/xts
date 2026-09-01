package core

import (
	"testing"

	"github.com/boxesandglue/boxesandglue/backend/node"
)

// wrapVList nests vl inside n plain VLists, the way CSSBuilder.CreateVlist
// wraps a table in html/body boxes before it reaches PlaceObject.
func wrapVList(vl *node.VList, n int) *node.VList {
	for i := 0; i < n; i++ {
		outer := node.NewVList()
		outer.List = vl
		vl = outer
	}
	return vl
}

func tableVList() *node.VList {
	vl := node.NewVList()
	vl.Attributes = node.H{
		"origin":        "table",
		"_headerCount":  1,
		"_buildHeaders": func() ([]*node.HList, error) { return nil, nil },
	}
	return vl
}

// TestSplittableTable checks what PlaceObject treats as a splittable table: a
// table that carries the header-repeat closure, however deeply CreateVlist has
// wrapped it. Anything else keeps the single-placement path.
func TestSplittableTable(t *testing.T) {
	table := tableVList()

	if got := splittableTable(table); got != table {
		t.Errorf("bare table: got %v, want the table itself", got)
	}
	if got := splittableTable(wrapVList(table, 2)); got != table {
		t.Errorf("wrapped table: got %v, want the table", got)
	}
	if got := splittableTable(wrapVList(node.NewVList(), 2)); got != nil {
		t.Errorf("table without a TableHead: got %v, want nil", got)
	}

	// A headerless table is placed as one object, as before: without header
	// rows there is nothing to repeat and no reason to prefer fragments.
	headerless := node.NewVList()
	headerless.Attributes = node.H{"origin": "table"}
	if got := splittableTable(headerless); got != nil {
		t.Errorf("headerless table: got %v, want nil", got)
	}
}

// TestFrameBottom checks the bottom edge of a frame, which is where the
// splitter decides a row will not fit. Row 1 of a frame starts at the top
// margin, so a frame h rows tall ends h grid heights below it.
func TestFrameBottom(t *testing.T) {
	g := &grid{
		gridHeight: 10 * 65536,
		marginTop:  20 * 65536,
	}
	a := &area{frame: []*gridRect{{row: 1, col: 1, width: 5, height: 4}}}

	// posY(4) is the top of the fourth row; the frame ends one row lower.
	want := g.posY(4, a) + g.gridHeight
	if got := g.frameBottom(a); got != want {
		t.Errorf("frameBottom = %s, want %s", got, want)
	}
	if got, want := g.frameBottom(a), g.posY(1, a)+4*g.gridHeight; got != want {
		t.Errorf("frameBottom = %s, want four rows below the first: %s", got, want)
	}
}
