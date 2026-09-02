package core

import (
	"sort"

	"github.com/boxesandglue/boxesandglue/backend/bag"
	"github.com/boxesandglue/boxesandglue/backend/node"
)

// A slateItem is an object placed on the slate. The x and y offsets are
// relative to the top left corner of the slate.
type slateItem struct {
	x  bag.ScaledPoint
	y  bag.ScaledPoint
	vl *node.VList
}

type slate struct {
	name     string
	grid     *grid
	items    []slateItem
	contents *node.VList
}

func (xd *xtsDocument) newSlate(slatename string) *slate {
	g := newGrid(xd)
	// Copy grid dimensions and margins from the current page grid so that
	// tables and other width-dependent commands work correctly inside slates.
	if cg := xd.currentGrid; cg != nil {
		g.nx = cg.nx
		g.ny = cg.ny
		g.marginLeft = cg.marginLeft
		g.marginRight = cg.marginRight
		g.marginTop = cg.marginTop
		g.marginBottom = cg.marginBottom
	}
	g.allocatedBlocks = make(allocationMatrix)
	g.areas[pageAreaName] = &area{
		name:  pageAreaName,
		frame: []*gridRect{{1, 1, coord(g.nx), coord(g.ny), 1, 1}},
	}
	g.inSlate = true
	s := slate{
		name: slatename,
		grid: g,
	}
	xd.slates[slatename] = &s
	return &s
}

// setGrid overrides the grid geometry copied from the page. Zero values keep
// the copied setting. Changing nx or ny resizes the placement area.
func (s *slate) setGrid(width, height, dx, dy bag.ScaledPoint, nx, ny int) {
	g := s.grid
	if width > 0 {
		g.gridWidth = width
	}
	if height > 0 {
		g.gridHeight = height
	}
	if dx > 0 {
		g.gridGapX = dx
	}
	if dy > 0 {
		g.gridGapY = dy
	}
	if nx > 0 {
		g.nx = nx
	}
	if ny > 0 {
		g.ny = ny
	}
	g.areas[pageAreaName].frame = []*gridRect{{1, 1, coord(g.nx), coord(g.ny), 1, 1}}
}

// appendItem puts an object onto the slate and invalidates the cached
// contents.
func (s *slate) appendItem(x, y bag.ScaledPoint, vl *node.VList) {
	s.items = append(s.items, slateItem{x: x, y: y, vl: vl})
	s.contents = nil
}

// buildContents composes the placed objects into a single VList and caches
// the result. The objects are stacked top to bottom with kerns in between, a
// negative kern moves back up for overlapping objects. The horizontal offset
// is a kern in an hbox around the object. An empty slate returns nil.
func (s *slate) buildContents() *node.VList {
	if s.contents != nil {
		return s.contents
	}
	if len(s.items) == 0 {
		return nil
	}
	items := make([]slateItem, len(s.items))
	copy(items, s.items)
	sort.SliceStable(items, func(i, j int) bool { return items[i].y < items[j].y })

	var head, tail node.Node
	appendNode := func(n node.Node) {
		head = node.InsertAfter(head, tail, n)
		tail = n
	}
	var curY, maxX, maxY bag.ScaledPoint
	for _, itm := range items {
		if delta := itm.y - curY; delta != 0 {
			k := node.NewKern()
			k.Kern = delta
			appendNode(k)
		}
		var n node.Node = itm.vl
		if itm.x != 0 {
			k := node.NewKern()
			k.Kern = itm.x
			node.InsertAfter(k, k, itm.vl)
			n = node.Hpack(k)
		}
		appendNode(n)
		curY = itm.y + itm.vl.Height + itm.vl.Depth
		if curY > maxY {
			maxY = curY
		}
		if right := itm.x + itm.vl.Width; right > maxX {
			maxX = right
		}
	}
	if delta := maxY - curY; delta != 0 {
		k := node.NewKern()
		k.Kern = delta
		appendNode(k)
	}
	vl := node.Vpack(head)
	vl.Width = maxX
	vl.Height = maxY
	vl.Depth = 0
	s.contents = vl
	return vl
}
