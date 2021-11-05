package core

import (
	"fmt"
	"strings"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/node"
	"github.com/speedata/goxml"
	"github.com/speedata/goxpath/xpath"
)

type commandFunc func(*goxml.Element, *xpath.Parser) (xpath.Sequence, error)

var (
	dataDispatcher map[string]map[string]*goxml.Element
	dispatchTable  map[string]commandFunc
)

func init() {
	dispatchTable = map[string]commandFunc{
		"B":                cmdB,
		"DefineFontfamily": cmdDefineFontfamily,
		"LoadFontfile":     cmdLoadFontfile,
		"Paragraph":        cmdParagraph,
		"PlaceObject":      cmdPlaceObject,
		"Record":           cmdRecord,
		"Textblock":        cmdTextblock,
		"Value":            cmdValue,
	}
}

func dispatch(layoutelement *goxml.Element, data *xpath.Parser) (xpath.Sequence, error) {
	var retSequence xpath.Sequence
	for _, cld := range layoutelement.Children() {
		if elt, ok := cld.(*goxml.Element); ok {
			if f, ok := dispatchTable[elt.Name]; ok {
				logger.Debugf("Call %s (line %d)", elt.Name, elt.Line)
				seq, err := f(elt, data)
				if err != nil {
					return nil, err
				}
				retSequence = append(retSequence, seq...)
			} else {
				fmt.Printf("Element %q unknown\n", elt.Name)
			}
		}
	}
	return retSequence, nil
}

func init() {
	dataDispatcher = make(map[string]map[string]*goxml.Element)
}

func cmdB(layoutelt *goxml.Element, dataelt *xpath.Parser) (xpath.Sequence, error) {
	abc, err := dispatch(layoutelt, dataelt)
	return abc, err
}

func cmdDefineFontfamily(layoutelt *goxml.Element, dataelt *xpath.Parser) (xpath.Sequence, error) {
	var size, leading bag.ScaledPoint
	var err error
	if size, err = getAttributeSize("fontsize", layoutelt, true, true, "", dataelt); err != nil {
		return nil, err
	}
	if leading, err = getAttributeSize("leading", layoutelt, true, true, "", dataelt); err != nil {
		return nil, err
	}
	var name, fontface string
	if name, err = getAttributeString("name", layoutelt, true, true, "", dataelt); err != nil {
		return nil, err
	}

	ff := fontfamily{
		name:    name,
		size:    size,
		leading: leading,
	}

	for _, cld := range layoutelt.Children() {
		if c, ok := cld.(*goxml.Element); ok {
			if fontface, err = getAttributeString("fontface", c, true, true, "", dataelt); err != nil {
				return nil, err
			}
			switch c.Name {
			case "Regular":
				if ff.weightStyleFontname == nil {
					ff.weightStyleFontname = make(map[fontWeight]map[fontStyle]string)
				}
				if ff.weightStyleFontname[fontWeight400] == nil {
					ff.weightStyleFontname[fontWeight400] = make(map[fontStyle]string)
				}
				ff.weightStyleFontname[fontWeight400][fontStyleNormal] = fontface
			case "Italic":
				if ff.weightStyleFontname == nil {
					ff.weightStyleFontname = make(map[fontWeight]map[fontStyle]string)
				}
				if ff.weightStyleFontname[fontWeight400] == nil {
					ff.weightStyleFontname[fontWeight400] = make(map[fontStyle]string)
				}
				ff.weightStyleFontname[fontWeight400][fontStyleItalic] = fontface
			case "Bold":
				if ff.weightStyleFontname == nil {
					ff.weightStyleFontname = make(map[fontWeight]map[fontStyle]string)
				}
				if ff.weightStyleFontname[fontWeight700] == nil {
					ff.weightStyleFontname[fontWeight700] = make(map[fontStyle]string)
				}
				ff.weightStyleFontname[fontWeight700][fontStyleNormal] = fontface
			case "BoldItalic":
				if ff.weightStyleFontname == nil {
					ff.weightStyleFontname = make(map[fontWeight]map[fontStyle]string)
				}
				if ff.weightStyleFontname[fontWeight700] == nil {
					ff.weightStyleFontname[fontWeight700] = make(map[fontStyle]string)
				}
				ff.weightStyleFontname[fontWeight700][fontStyleItalic] = fontface
			}
		}
	}
	definefontfamily(&ff)
	return nil, nil
}

func cmdLoadFontfile(layoutelt *goxml.Element, dataelt *xpath.Parser) (xpath.Sequence, error) {
	var filename, name string
	var err error
	if filename, err = getAttributeString("filename", layoutelt, true, true, "", dataelt); err != nil {
		return nil, err
	}
	if name, err = getAttributeString("name", layoutelt, true, true, "", dataelt); err != nil {
		return nil, err
	}

	fontNameFontFile[name] = &fontfile{
		filename: filename,
	}
	return nil, nil
}

func cmdRecord(layoutelt *goxml.Element, dataelt *xpath.Parser) (xpath.Sequence, error) {
	var elt, mode string
	var err error
	if elt, err = getAttributeString("element", layoutelt, true, false, "", dataelt); err != nil {
		return nil, err
	}
	if mode, err = getAttributeString("mode", layoutelt, false, false, "", dataelt); err != nil {
		return nil, err
	}
	dp := dataDispatcher[elt]
	if dp == nil {
		dataDispatcher[elt] = make(map[string]*goxml.Element)
	}
	dataDispatcher[elt][mode] = layoutelt
	return nil, nil
}

func cmdParagraph(layoutelt *goxml.Element, dataelt *xpath.Parser) (xpath.Sequence, error) {
	abc, err := dispatch(layoutelt, dataelt)
	logger.Debug("paragraph after dispatch")
	if err != nil {
		return nil, err
	}
	var txt []string
	for _, itm := range abc {
		switch t := itm.(type) {
		case *goxml.Element:
			txt = append(txt, t.Stringvalue())
		}
	}
	return xpath.Sequence{paragraph{strings.Join(txt, "")}}, nil
}

func cmdPlaceObject(layoutelt *goxml.Element, dataelt *xpath.Parser) (xpath.Sequence, error) {
	seq, err := dispatch(layoutelt, dataelt)
	if err != nil {
		return nil, err
	}
	vl := seq[0].(*node.VList)
	doc.OutputAt(bag.MustSp("2cm"), bag.MustSp("20cm"), vl)
	return seq, nil
}

func cmdTextblock(layoutelt *goxml.Element, dataelt *xpath.Parser) (xpath.Sequence, error) {
	seq, err := dispatch(layoutelt, dataelt)
	if err != nil {
		return nil, err
	}

	opts := formatOptions{
		width: bag.MustSp("160pt"),
		ff:    fontfamilies["text"],
	}
	var vlist *node.VList
	for _, itm := range seq {
		switch t := itm.(type) {
		case paragraph:
			vlist, err = t.format(opts)
			if err != nil {
				return nil, err
			}

		}
	}
	return xpath.Sequence{vlist}, nil
}

func cmdValue(layoutelt *goxml.Element, dataelt *xpath.Parser) (xpath.Sequence, error) {
	var selection string
	var err error

	if selection, err = getAttributeString("select", layoutelt, false, false, "", dataelt); err != nil {
		return nil, err
	}
	if selection != "" {
		eval, err := dataelt.Evaluate(selection)
		if err != nil {
			return nil, err
		}
		return eval, nil
	}

	fmt.Println(layoutelt.Children())
	return nil, nil
}
