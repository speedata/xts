package core

import "github.com/speedata/boxesandglue/backend/node"

type group struct {
	name     string
	grid     *grid
	contents *node.VList
}

func (xd *xtsDocument) newGroup(groupname string) *group {
	g := newGrid(xd)
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
