package core

import (
	"reflect"
	"testing"

	"github.com/boxesandglue/boxesandglue/backend/document"
	"github.com/boxesandglue/boxesandglue/backend/node"
	"github.com/boxesandglue/boxesandglue/frontend"
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
// table, however deeply CreateVlist has wrapped it. Anything else keeps the
// single-placement path.
func TestSplittableTable(t *testing.T) {
	table := tableVList()

	if got := splittableTable(table); got != table {
		t.Errorf("bare table: got %v, want the table itself", got)
	}
	if got := splittableTable(wrapVList(table, 2)); got != table {
		t.Errorf("wrapped table: got %v, want the table", got)
	}
	if got := splittableTable(wrapVList(node.NewVList(), 2)); got != nil {
		t.Errorf("wrapper without a table: got %v, want nil", got)
	}

	// A table without a <TableHead> splits too: placed as one object, a
	// table taller than the frame runs off the bottom of the page.
	headerless := node.NewVList()
	headerless.Attributes = node.H{"origin": "table"}
	if got := splittableTable(wrapVList(headerless, 2)); got != headerless {
		t.Errorf("headerless table: got %v, want the table", got)
	}

	// A table next to other content, e.g. a paragraph in mixed <HTML>
	// content, keeps the single-placement path: splitTable lays out the
	// table alone and would drop the sibling.
	sibling := node.NewVList()
	tblWithSibling := tableVList()
	var head node.Node
	head = node.InsertAfter(head, nil, sibling)
	head = node.InsertAfter(head, sibling, tblWithSibling)
	body := node.NewVList()
	body.List = head
	if got := splittableTable(wrapVList(body, 1)); got != nil {
		t.Errorf("table with sibling content: got %v, want nil", got)
	}

	// Glue of no width next to the wrapper carries no content and is
	// tolerated.
	glue := node.NewGlue()
	tblAfterGlue := tableVList()
	head = nil
	head = node.InsertAfter(head, nil, glue)
	head = node.InsertAfter(head, glue, tblAfterGlue)
	wrapper := node.NewVList()
	wrapper.List = head
	if got := splittableTable(wrapper); got != tblAfterGlue {
		t.Errorf("table behind glue: got %v, want the table", got)
	}

	// Spacing with a width is the wrapper's padding or margin, which the
	// fragments would lose.
	for name, space := range map[string]node.Node{
		"glue": func() node.Node { g := node.NewGlue(); g.Width = 4 * 65536; return g }(),
		"kern": func() node.Node { k := node.NewKern(); k.Kern = 4 * 65536; return k }(),
	} {
		tbl := tableVList()
		head = nil
		head = node.InsertAfter(head, nil, space)
		head = node.InsertAfter(head, space, tbl)
		padded := node.NewVList()
		padded.List = head
		if got := splittableTable(padded); got != nil {
			t.Errorf("table behind %s with a width: got %v, want nil", name, got)
		}
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

// TestKeepsWithNext checks which rows the splitter holds to the row after
// them: an HList that frontend.BuildTable marked for a rowspan, and nothing
// else.
func TestKeepsWithNext(t *testing.T) {
	marked := node.NewHList()
	marked.Attributes = node.H{"_keepWithNext": true}
	if !keepsWithNext(marked) {
		t.Error("a row marked _keepWithNext does not keep with the next")
	}
	if keepsWithNext(node.NewHList()) {
		t.Error("an unmarked row keeps with the next")
	}
	unmarked := node.NewHList()
	unmarked.Attributes = node.H{"_keepWithNext": false}
	if keepsWithNext(unmarked) {
		t.Error("a row marked false keeps with the next")
	}
	if keepsWithNext(node.NewGlue()) {
		t.Error("glue between rows keeps with the next")
	}
}

// splitRows splits a table of n rows, each one grid row high, across three
// frames six rows high, with rows first..last joined by a rowspan. It returns
// the row numbers each fragment holds.
func splitRows(t *testing.T, n, first, last int) [][]int {
	t.Helper()
	xd, _ := threeFrames(5, 5, 5)
	g := xd.currentGrid
	g.gridWidth, g.gridHeight = 10*65536, 10*65536
	xd.currentPage.bagPage = &document.Page{}

	table := node.NewVList()
	table.Width = g.width(5)
	var tail node.Node
	for i := 1; i <= n; i++ {
		r := node.NewHList()
		r.Height = g.gridHeight
		r.Attributes = node.H{"row": i, "_keepWithNext": i >= first && i < last}
		table.List = node.InsertAfter(table.List, tail, r)
		tail = r
	}

	var fragments [][]int
	record := func(vl *node.VList) {
		var rows []int
		for r := vl.List; r != nil; r = r.Next() {
			rows = append(rows, r.(*node.HList).Attributes["row"].(int))
		}
		fragments = append(fragments, rows)
	}
	if err := xd.splitTable(table, "cols", 1, 1, "", false, frontend.HAlignLeft, record); err != nil {
		t.Fatal(err)
	}
	return fragments
}

// TestSplitTableRowspan checks that the rows a rowspan joins go to the next
// frame together, and that a group taller than a frame is broken inside
// instead of running past the bottom of the frame.
func TestSplitTableRowspan(t *testing.T) {
	for _, tc := range []struct {
		name              string
		rows, first, last int
		want              [][]int
	}{
		{"group fits a frame", 10, 5, 7, [][]int{{1, 2, 3, 4}, {5, 6, 7, 8, 9, 10}}},
		{"group taller than a frame", 14, 3, 12, [][]int{{1, 2}, {3, 4, 5, 6, 7, 8}, {9, 10, 11, 12, 13, 14}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := splitRows(t, tc.rows, tc.first, tc.last); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("fragments %v, want %v", got, tc.want)
			}
		})
	}
}
