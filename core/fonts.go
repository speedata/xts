package core

// import (
// 	"github.com/speedata/boxesandglue/backend/bag"
// 	"github.com/speedata/boxesandglue/backend/font"
// 	"github.com/speedata/boxesandglue/pdfbackend/pdf"
// )

// var (
// 	fontNameFontFile = make(map[string]*fontfile)
// 	fontfamilies     = make(map[string]*fontfamily)
// )

// type fontfile struct {
// 	filename string
// }

// type fontWeight int
// type fontStyle int

// const (
// 	fontWeight100 fontWeight = iota
// 	fontWeight200
// 	fontWeight300
// 	fontWeight400 // normal
// 	fontWeight500
// 	fontWeight600
// 	fontWeight700 // bold
// 	fontWeight800
// 	fontWeight900
// )

// const (
// 	fontStyleNormal fontStyle = iota
// 	fontStyleItalic
// 	fontStyleOblique
// )

// type fontfamily struct {
// 	xd                  *xtsDocument
// 	id                  int
// 	name                string
// 	weightStyleFontname map[fontWeight]map[fontStyle]string
// 	weightStyleFace     map[fontWeight]map[fontStyle]*pdf.Face
// 	size                bag.ScaledPoint
// 	leading             bag.ScaledPoint
// }

// func resolveFont(name string) string {
// 	logger.Debugf("Resolve font %q", name)
// 	if ff, ok := fontNameFontFile[name]; ok {
// 		return ff.filename
// 	}

// 	return ""
// }

// func (ff *fontfamily) getFont(weight fontWeight, style fontStyle, size bag.ScaledPoint) (*font.Font, error) {
// 	logger.Debugf("Get font weight %d style %d size %s", weight, style, size.String())
// 	if face, ok := ff.weightStyleFace[weight][style]; ok {
// 		return ff.xd.doc.CreateFont(face, size), nil
// 	}
// 	fontfile := resolveFont(ff.weightStyleFontname[weight][style])
// 	fn, err := findFile(fontfile)
// 	if err != nil {
// 		return nil, err
// 	}
// 	face, err := ff.xd.doc.LoadFace(fn, 0)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if ff.weightStyleFace == nil {
// 		ff.weightStyleFace = make(map[fontWeight]map[fontStyle]*pdf.Face)
// 	}

// 	if ff.weightStyleFace[weight] == nil {
// 		ff.weightStyleFace[weight] = make(map[fontStyle]*pdf.Face)
// 	}

// 	ff.weightStyleFace[weight][style] = face

// 	return ff.xd.doc.CreateFont(face, size), nil
// }

// func definefontfamily(ff *fontfamily) {
// 	ff.id = len(fontfamilies)
// 	fontfamilies[ff.name] = ff
// 	logger.Infof("Define fontfamily %s (id %d)", ff.name, ff.id)
// }
