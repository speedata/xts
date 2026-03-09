package core

import "github.com/boxesandglue/boxesandglue/backend/node"

type group struct {
	name     string
	grid     *grid
	contents *node.VList
}

func (xd *xtsDocument) newGroup(groupname string) *group {
	g := newGrid(xd)
	// Copy grid dimensions and margins from the current page grid so that
	// tables and other width-dependent commands work correctly inside groups.
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
	g.inGroup = true
	gr := group{
		name: groupname,
		grid: g,
	}
	xd.groups[groupname] = &gr
	return &gr
}
