package core

import (
	"testing"
	"time"

	"github.com/boxesandglue/boxesandglue/backend/bag"
	"github.com/boxesandglue/boxesandglue/backend/node"
)

// A slate built while it is composed, as sd:slate-height does, and built
// again after another object holds every object once. Before, the second
// build kept links from the first: one below the other the list became a
// cycle, beside each other the second object came twice.
func TestSlateBuildsAgainAfterAnotherObject(t *testing.T) {
	for _, c := range []struct {
		name    string
		secondX bag.ScaledPoint
	}{
		{"one below the other", 0},
		{"beside each other", bag.MustSP("4cm")},
	} {
		t.Run(c.name, func(t *testing.T) {
			obj := func() *node.VList {
				vl := node.NewVList()
				vl.Width, vl.Height = bag.MustSP("3cm"), bag.MustSP("12pt")
				return vl
			}
			first, second, third := obj(), obj(), obj()
			s := &slate{}
			s.appendItem(slateItem{vl: first})
			secondY := first.Height
			if c.secondX != 0 {
				secondY = 0
			}
			s.appendItem(slateItem{x: c.secondX, y: secondY, vl: second})
			s.buildContents()
			s.appendItem(slateItem{y: secondY + second.Height, vl: third})

			done := make(chan map[*node.VList]int, 1)
			go func() {
				seen := map[*node.VList]int{}
				var count func(n node.Node)
				count = func(n node.Node) {
					for steps := 0; n != nil && steps < 100; n, steps = n.Next(), steps+1 {
						switch v := n.(type) {
						case *node.VList:
							if v == first || v == second || v == third {
								seen[v]++
							} else {
								count(v.List)
							}
						case *node.HList:
							count(v.List)
						}
					}
				}
				count(s.buildContents().List)
				done <- seen
			}()
			select {
			case seen := <-done:
				for i, vl := range []*node.VList{first, second, third} {
					if seen[vl] != 1 {
						t.Errorf("object %d is in the slate %d times, want once", i+1, seen[vl])
					}
				}
			case <-time.After(5 * time.Second):
				t.Fatal("the second build did not end")
			}
		})
	}
}
