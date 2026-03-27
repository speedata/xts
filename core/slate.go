package core

import "github.com/boxesandglue/boxesandglue/backend/node"

type slate struct {
	name     string
	grid     *grid
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
