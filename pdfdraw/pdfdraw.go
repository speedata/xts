package pdfdraw

import (
	"fmt"
	"strings"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/document"
)

// Object represents a set of PDF instructions to draw a PDF graphic
type Object struct {
	pdfstring []string
}

// New creates a new PDF object.
func New() *Object {
	return &Object{}
}

// Color sets the stroking and nonstroking color
func (pd *Object) Color(col document.Color) *Object {
	pd.ColorStroking(col).ColorNonstroking(col)
	return pd
}

// ColorStroking sets the stroking color
func (pd *Object) ColorStroking(col document.Color) *Object {
	if col.Space != document.ColorNone {
		pd.pdfstring = append(pd.pdfstring, col.PDFStringFG())
	}
	return pd
}

// ColorNonstroking sets the non stroking color. If the color is the color
// “none”, then no color will be set.
func (pd *Object) ColorNonstroking(col document.Color) *Object {
	if col.Space != document.ColorNone {
		pd.pdfstring = append(pd.pdfstring, col.PDFStringFG())
	}
	return pd
}

// Curveto appends a bezier cuve from the current point to point 3 controlled by
// points 1 and 2.
func (pd *Object) Curveto(x1, y1, x2, y2, x3, y3 bag.ScaledPoint) *Object {
	pd.pdfstring = append(pd.pdfstring, fmt.Sprintf("%s %s %s %s %s %s c", x1, y1, x2, y2, x3, y3))
	return pd
}

// Moveto moves the cursor relative to the current point.
func (pd *Object) Moveto(x, y bag.ScaledPoint) *Object {
	pd.pdfstring = append(pd.pdfstring, fmt.Sprintf("%s %s m", x, y))
	return pd
}

// Lineto draws a straight line from the current point to the point given at x and y.
func (pd *Object) Lineto(x, y bag.ScaledPoint) *Object {
	pd.pdfstring = append(pd.pdfstring, fmt.Sprintf("%s %s l", x, y))
	return pd
}

// Circle draws a circle. TODO: document where it starts/ends etc.
func (pd *Object) Circle(x, y, radiusX, radiusY bag.ScaledPoint) *Object {
	circleBezier := 0.551915024494

	shiftDown, shiftRight := -2*radiusY, bag.ScaledPoint(0)
	dx := bag.ScaledPointFromFloat(radiusX.ToPT() * (1 - circleBezier))
	dy := bag.ScaledPointFromFloat(radiusY.ToPT() * (1 - circleBezier))

	x1 := shiftRight
	y1 := shiftDown + radiusY
	x2 := x1
	y2 := shiftDown + radiusY*2 - dy
	x3 := shiftRight + dx
	y3 := shiftDown + radiusY*2
	x4 := shiftRight + radiusX
	y4 := shiftDown + radiusY*2
	x5 := shiftRight + radiusX*2 - dx
	y5 := y3
	x6 := shiftRight + radiusX*2
	y6 := y2
	x7 := x6
	y7 := y1
	x8 := x6
	y8 := shiftDown + dy
	x9 := x5
	y9 := shiftDown
	x10 := x4
	y10 := y9
	x11 := x3
	y11 := y9
	x12 := x1
	y12 := y8
	pd.Moveto(x1, y1)
	pd.Curveto(x2, y2, x3, y3, x4, y4)
	pd.Curveto(x5, y5, x6, y6, x7, y7)
	pd.Curveto(x8, y8, x9, y9, x10, y10)
	pd.Curveto(x11, y11, x12, y12, x1, y1)
	return pd
}

// Rect draws a rectangle
func (pd *Object) Rect(x, y, wd, ht bag.ScaledPoint) *Object {
	pd.pdfstring = append(pd.pdfstring, fmt.Sprintf("%s %s %s %s re", x, y, wd, ht))
	return pd
}

// Fill fills the current object
func (pd *Object) Fill() *Object {
	pd.pdfstring = append(pd.pdfstring, "f")
	return pd
}

func (pd *Object) String() string {
	return strings.Join(pd.pdfstring, " ")
}
