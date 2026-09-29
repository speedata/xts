package core

import (
	"testing"

	"github.com/boxesandglue/boxesandglue/backend/node"
	"github.com/boxesandglue/htmlbag"
	xpath "github.com/speedata/goxpath"
	"golang.org/x/net/html"
)

// TestCellKeepsNestedTablesInPlace checks that nested tables stay where they
// are among a cell's paragraphs and that a second one does not replace the
// first.
func TestCellKeepsNestedTablesInPlace(t *testing.T) {
	xd := &xtsDocument{cssbuilder: &htmlbag.CSSBuilder{PendingVLists: map[string]*node.VList{}}}
	pA := &html.Node{Type: html.ElementNode, Data: "p"}
	pB := &html.Node{Type: html.ElementNode, Data: "p"}
	vl1, vl2 := node.NewVList(), node.NewVList()
	td := &html.Node{Type: html.ElementNode, Data: "td"}
	appendCellContents(xd, td, xpath.Sequence{pA, vl1, pB, vl2})

	for _, a := range td.Attr {
		if a.Key == "data-vlist-id" {
			t.Errorf("td has data-vlist-id %q, want it on a child", a.Val)
		}
	}
	var got []any
	for c := td.FirstChild; c != nil; c = c.NextSibling {
		var vlid string
		for _, a := range c.Attr {
			if a.Key == "data-vlist-id" {
				vlid = a.Val
			}
		}
		if vlid == "" {
			got = append(got, c)
			continue
		}
		vl, ok := xd.cssbuilder.PendingVLists[vlid]
		if !ok {
			t.Fatalf("child %q is not a pending VList", vlid)
		}
		got = append(got, vl)
	}
	want := []any{pA, vl1, pB, vl2}
	if len(got) != len(want) {
		t.Fatalf("td has %d children, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("child %d = %p, want %p", i, got[i], want[i])
		}
	}
}
