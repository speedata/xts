package core

import (
	"testing"

	"github.com/boxesandglue/boxesandglue/backend/bag"
	"github.com/boxesandglue/boxesandglue/backend/document"
	"github.com/boxesandglue/boxesandglue/backend/node"
)

func box(wd, ht string) *node.VList {
	vl := node.NewVList()
	vl.Width, vl.Height = bag.MustSP(wd), bag.MustSP(ht)
	return vl
}

// drawOrder lists the given objects in the order the slate draws them.
func drawOrder(s *slate, objs ...*node.VList) []*node.VList {
	var got []*node.VList
	var walk func(n node.Node)
	walk = func(n node.Node) {
		for ; n != nil; n = n.Next() {
			switch t := n.(type) {
			case *node.VList:
				found := false
				for _, o := range objs {
					if t == o {
						got = append(got, t)
						found = true
					}
				}
				if !found {
					walk(t.List)
				}
			case *node.HList:
				walk(t.List)
			}
		}
	}
	walk(s.buildContents().List)
	return got
}

// An object placed behind in a slate is drawn under the slate's other
// objects and still takes its room.
func TestSlateLayerBehind(t *testing.T) {
	s := &slate{}
	a, b, panel := box("100pt", "20pt"), box("50pt", "10pt"), box("120pt", "40pt")
	s.appendItem(slateItem{vl: a})
	s.appendItem(slateItem{y: bag.MustSP("20pt"), vl: b})
	s.appendItem(slateItem{x: bag.MustSP("-5pt"), y: bag.MustSP("5pt"), vl: panel, behind: true})
	want := []*node.VList{panel, a, b}
	got := drawOrder(s, a, b, panel)
	if len(got) != len(want) {
		t.Fatalf("slate draws %d objects, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("object %d is %p, want %p", i, got[i], want[i])
		}
	}
	vl := s.buildContents()
	if want := bag.MustSP("115pt"); vl.Width != want {
		t.Errorf("slate width %s, want %s", vl.Width, want)
	}
	if want := bag.MustSP("45pt"); vl.Height+vl.Depth != want {
		t.Errorf("slate height %s, want %s", vl.Height+vl.Depth, want)
	}
}

// layer="behind" goes under everything placed on the page but what its
// AtPageCreation drew, in the order such objects are placed.
func TestOutputBehind(t *testing.T) {
	a, b, c, d := box("1pt", "1pt"), box("1pt", "1pt"), box("1pt", "1pt"), box("1pt", "1pt")
	p := &page{bagPage: &document.Page{}, underlay: 1}
	p.bagPage.Objects = []document.Object{{Vlist: a}, {Vlist: b}}
	p.outputBehind(0, 0, c)
	p.outputBehind(0, 0, d)
	want := []*node.VList{a, c, d, b}
	for i, o := range p.bagPage.Objects {
		if o.Vlist != want[i] {
			t.Errorf("object %d is %p, want %p", i, o.Vlist, want[i])
		}
	}
}
