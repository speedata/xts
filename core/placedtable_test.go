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

// TestPlacedTable checks what PlaceObject treats as a table: a table, however
// deeply CreateVlist has wrapped it, but not one among other content.
func TestPlacedTable(t *testing.T) {
	table := tableVList()

	if !placedTable(table) {
		t.Error("bare table: not a table")
	}
	if !placedTable(wrapVList(table, 2)) {
		t.Error("wrapped table: not a table")
	}
	if placedTable(wrapVList(node.NewVList(), 2)) {
		t.Error("wrapper without a table: a table")
	}

	// A table without a <TableHead> is one too.
	headerless := node.NewVList()
	headerless.Attributes = node.H{"origin": "table"}
	if !placedTable(wrapVList(headerless, 2)) {
		t.Error("headerless table: not a table")
	}

	// A table next to other content, e.g. a paragraph in mixed <HTML>
	// content, is not.
	sibling := node.NewVList()
	tblWithSibling := tableVList()
	var head node.Node
	head = node.InsertAfter(head, nil, sibling)
	head = node.InsertAfter(head, sibling, tblWithSibling)
	body := node.NewVList()
	body.List = head
	if placedTable(wrapVList(body, 1)) {
		t.Error("table with sibling content: a table")
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
	if !placedTable(wrapper) {
		t.Error("table behind glue: not a table")
	}

	// Spacing with a width is the wrapper's padding or margin.
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
		if placedTable(padded) {
			t.Errorf("table behind %s with a width: a table", name)
		}
	}
}
