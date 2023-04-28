package core

import (
	"fmt"
	"strings"

	"github.com/speedata/goxpath"
)

const fnNS = "urn:speedata.de/2021/xtsfunctions/en"

func init() {
	goxpath.RegisterFunction(&goxpath.Function{Name: "current-page", Namespace: fnNS, F: fnCurrentPage, MinArg: 0, MaxArg: 0})
	goxpath.RegisterFunction(&goxpath.Function{Name: "current-row", Namespace: fnNS, F: fnCurrentRow, MinArg: 0, MaxArg: 1})
	goxpath.RegisterFunction(&goxpath.Function{Name: "dummytext", Namespace: fnNS, F: fnDummytext, MinArg: 0, MaxArg: 1})
	goxpath.RegisterFunction(&goxpath.Function{Name: "even", Namespace: fnNS, F: fnEven, MinArg: 1, MaxArg: 1})
	goxpath.RegisterFunction(&goxpath.Function{Name: "file-exists", Namespace: fnNS, F: fnFileExists, MinArg: 1, MaxArg: 1})
	goxpath.RegisterFunction(&goxpath.Function{Name: "group-height", Namespace: fnNS, F: fnGroupheight, MinArg: 1, MaxArg: 2})
	goxpath.RegisterFunction(&goxpath.Function{Name: "group-width", Namespace: fnNS, F: fnGroupwidth, MinArg: 1, MaxArg: 2})
	goxpath.RegisterFunction(&goxpath.Function{Name: "last-page-number", Namespace: fnNS, F: fnLastPagenumber, MinArg: 0, MaxArg: 0})
	goxpath.RegisterFunction(&goxpath.Function{Name: "mode", Namespace: fnNS, F: fnMode, MinArg: 1, MaxArg: 1})
	goxpath.RegisterFunction(&goxpath.Function{Name: "number-of-columns", Namespace: fnNS, F: fnNumberOfColumns, MinArg: 0, MaxArg: 1})
	goxpath.RegisterFunction(&goxpath.Function{Name: "number-of-rows", Namespace: fnNS, F: fnNumberOfRows, MinArg: 0, MaxArg: 1})
	goxpath.RegisterFunction(&goxpath.Function{Name: "odd", Namespace: fnNS, F: fnOdd, MinArg: 1, MaxArg: 1})
	goxpath.RegisterFunction(&goxpath.Function{Name: "page-number", Namespace: fnNS, F: fnPagenumber, MinArg: 1, MaxArg: 1})
	goxpath.RegisterFunction(&goxpath.Function{Name: "roman-numeral", Namespace: fnNS, F: fnRomannumeral, MinArg: 1, MaxArg: 1})
}

func fnCurrentPage(ctx *goxpath.Context, args []goxpath.Sequence) (goxpath.Sequence, error) {
	xd := ctx.Store["xd"].(*xtsDocument)
	cp := xd.currentPagenumber
	return goxpath.Sequence{cp}, nil
}

func fnCurrentRow(ctx *goxpath.Context, args []goxpath.Sequence) (goxpath.Sequence, error) {
	areaname := defaultAreaName
	if len(args) > 0 {
		firstArg := args[0]
		areaname = firstArg.Stringvalue()
	}
	var area *area
	var ok bool
	xd := ctx.Store["xd"].(*xtsDocument)
	if area, ok = xd.currentGrid.areas[areaname]; !ok {
		return nil, fmt.Errorf("area %s unknown", areaname)
	}
	return goxpath.Sequence{int(area.CurrentRow())}, nil
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
	nv, err := args[0].IntValue()
	if err != nil {
		return nil, err
	}
	return goxpath.Sequence{nv%2 == 0}, nil
}

func fnLastPagenumber(ctx *goxpath.Context, args []goxpath.Sequence) (goxpath.Sequence, error) {
	xd := ctx.Store["xd"].(*xtsDocument)
	if a := xd.aux; a != nil {
		return goxpath.Sequence{a.LastPage}, nil
	}
	return goxpath.Sequence{0}, nil
}

func fnMode(ctx *goxpath.Context, args []goxpath.Sequence) (goxpath.Sequence, error) {
	xd := ctx.Store["xd"].(*xtsDocument)
	findMode := args[0].Stringvalue()

	for _, mode := range xd.cfg.Mode {
		if mode == findMode {
			return goxpath.Sequence{true}, nil
		}
	}
	return goxpath.Sequence{false}, nil
}

func fnNumberOfColumns(ctx *goxpath.Context, args []goxpath.Sequence) (goxpath.Sequence, error) {
	areaname := defaultAreaName
	if len(args) > 0 {
		firstArg := args[0]
		areaname = firstArg.Stringvalue()
	}
	var area *area
	var ok bool
	xd := ctx.Store["xd"].(*xtsDocument)
	if area, ok = xd.currentGrid.areas[areaname]; !ok {
		return nil, fmt.Errorf("area %s unknown", areaname)
	}
	return goxpath.Sequence{int(area.frame[area.currentFrame].width)}, nil
}

func fnNumberOfRows(ctx *goxpath.Context, args []goxpath.Sequence) (goxpath.Sequence, error) {
	areaname := defaultAreaName
	if len(args) > 0 {
		firstArg := args[0]
		areaname = firstArg.Stringvalue()
	}
	var area *area
	var ok bool
	xd := ctx.Store["xd"].(*xtsDocument)
	if area, ok = xd.currentGrid.areas[areaname]; !ok {
		return nil, fmt.Errorf("area %s unknown", areaname)
	}
	return goxpath.Sequence{int(area.frame[area.currentFrame].height)}, nil
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

func fnPagenumber(ctx *goxpath.Context, args []goxpath.Sequence) (goxpath.Sequence, error) {
	markerName := args[0][0].(string)
	xd := ctx.Store["xd"].(*xtsDocument)
	m, found := xd.getMarker(markerName)
	if !found {
		return goxpath.Sequence{0}, nil
	}
	return goxpath.Sequence{m.pagenumber}, nil
}

func fnRomannumeral(ctx *goxpath.Context, args []goxpath.Sequence) (goxpath.Sequence, error) {
	number, err := args[0].IntValue()
	if err != nil {
		return nil, err
	}
	maxRomanNumber := 3999
	if number > maxRomanNumber || number < 0 {
		return nil, fmt.Errorf("romannumeral: number out of range: %d (0-3999)", number)
	}

	conversions := []struct {
		value int
		digit string
	}{
		{1000, "M"},
		{900, "CM"},
		{500, "D"},
		{400, "CD"},
		{100, "C"},
		{90, "XC"},
		{50, "L"},
		{40, "XL"},
		{10, "X"},
		{9, "IX"},
		{5, "V"},
		{4, "IV"},
		{1, "I"},
	}

	var roman strings.Builder
	for _, conversion := range conversions {
		for number >= conversion.value {
			roman.WriteString(conversion.digit)
			number -= conversion.value
		}
	}

	return goxpath.Sequence{roman.String()}, nil
}

func fnOdd(ctx *goxpath.Context, args []goxpath.Sequence) (goxpath.Sequence, error) {
	nv, err := goxpath.NumberValue(args[0])
	if err != nil {
		return nil, err
	}
	return goxpath.Sequence{int(nv)%2 == 1}, nil
}
