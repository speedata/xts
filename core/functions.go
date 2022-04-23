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

func fnFileExists(ctx *goxpath.Context, args []goxpath.Sequence) (goxpath.Sequence, error) {
	seq := args[0]
	fmt.Println(seq)
	return goxpath.Sequence{true}, nil
}

func fnOdd(ctx *goxpath.Context, args []goxpath.Sequence) (goxpath.Sequence, error) {
	nv, err := goxpath.NumberValue(args[0])
	if err != nil {
		return nil, err
	}
	return goxpath.Sequence{int(nv)%2 == 1}, nil
}
