package core

import (
	"fmt"
	"strings"

	"github.com/speedata/goxpath"
)

const fnNS = "urn:speedata.de/2021/xtsfunctions/en"

func init() {
	goxpath.RegisterFunction(&goxpath.Function{Name: "dummytext", Namespace: fnNS, F: fnDummytext, MinArg: 0, MaxArg: 1})
	goxpath.RegisterFunction(&goxpath.Function{Name: "even", Namespace: fnNS, F: fnEven, MinArg: 1, MaxArg: 1})
	goxpath.RegisterFunction(&goxpath.Function{Name: "group-height", Namespace: fnNS, F: fnGroupheight, MinArg: 1, MaxArg: 2})
	goxpath.RegisterFunction(&goxpath.Function{Name: "group-width", Namespace: fnNS, F: fnGroupwidth, MinArg: 1, MaxArg: 2})
	goxpath.RegisterFunction(&goxpath.Function{Name: "file-exists", Namespace: fnNS, F: fnFileExists, MinArg: 1, MaxArg: 1})
	goxpath.RegisterFunction(&goxpath.Function{Name: "odd", Namespace: fnNS, F: fnOdd, MinArg: 1, MaxArg: 1})
}

func fnDummytext(ctx *goxpath.Context, args []goxpath.Sequence) (goxpath.Sequence, error) {
	str := `Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod
	tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim
	veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea
	commodo consequat. Duis aute irure dolor in reprehenderit in voluptate
	velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint
	occaecat cupidatat non proident, sunt in culpa qui officia deserunt
	mollit anim id est laborum.`
	ret := strings.Join(strings.Fields(str), " ")
	return goxpath.Sequence{ret}, nil
}

func fnEven(ctx *goxpath.Context, args []goxpath.Sequence) (goxpath.Sequence, error) {
	nv, err := goxpath.NumberValue(args[0])
	if err != nil {
		return nil, err
	}
	return goxpath.Sequence{int(nv)%2 == 0}, nil
}

func fnGroupheight(ctx *goxpath.Context, args []goxpath.Sequence) (goxpath.Sequence, error) {
	groupname := args[0][0].(string)
	xd := ctx.Store["xd"].(*xtsDocument)
	xd.setupPage()
	if grp, ok := xd.groups[groupname]; ok {
		groupheight := grp.contents.Height
		if len(args) == 1 {
			return goxpath.Sequence{xd.currentGrid.heightToRows(groupheight)}, nil
		}
		unit := args[1][0].(string)
		val, err := groupheight.ToUnit(unit)
		if err != nil {
			return nil, err
		}
		return goxpath.Sequence{val}, nil

	}
	return nil, fmt.Errorf("sd:group-height() group %q not found", groupname)
}

func fnGroupwidth(ctx *goxpath.Context, args []goxpath.Sequence) (goxpath.Sequence, error) {
	groupname := args[0][0].(string)
	xd := ctx.Store["xd"].(*xtsDocument)
	xd.setupPage()
	if grp, ok := xd.groups[groupname]; ok {
		groupwidth := grp.contents.Width
		if len(args) == 1 {
			return goxpath.Sequence{xd.currentGrid.widthToColumns(groupwidth)}, nil
		}
		unit := args[1][0].(string)
		val, err := groupwidth.ToUnit(unit)
		if err != nil {
			return nil, err
		}
		return goxpath.Sequence{val}, nil

	}
	return nil, fmt.Errorf("sd:group-height() group %q not found", groupname)

}

func fnFileExists(ctx *goxpath.Context, args []goxpath.Sequence) (goxpath.Sequence, error) {
	seq := args[0]
	return goxpath.Sequence{fileexists(seq.Stringvalue())}, nil
}

func fnOdd(ctx *goxpath.Context, args []goxpath.Sequence) (goxpath.Sequence, error) {
	nv, err := goxpath.NumberValue(args[0])
	if err != nil {
		return nil, err
	}
	return goxpath.Sequence{int(nv)%2 == 1}, nil
}
