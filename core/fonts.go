package core

import (
	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/frontend"
)

var (
	tenpoint    = bag.MustSp("10pt")
	twelvepoint = bag.MustSp("12pt")
)

func (xd *xtsDocument) defaultfont() error {
	xd.fontsources = make(map[string]*frontend.FontSource)
	ff := xd.document.NewFontFamily("text")
	data := []struct {
		fontname string
		filename string
	}{
		{"TeXGyreHeros-Regular", "texgyreheros-regular.otf"},
		{"TeXGyreHeros-Italic", "texgyreheros-italic.otf"},
		{"TeXGyreHeros-Bold", "texgyreheros-bold.otf"},
		{"TeXGyreHeros-BoldItalic", "texgyreheros-bolditalic.otf"},
	}
	for _, font := range data {
		fn, err := xd.cfg.FindFile(font.filename)
		if err != nil {
			return err
		}
		fs := frontend.FontSource{
			Name:   font.fontname,
			Source: fn,
		}
		xd.fontsources[font.fontname] = &fs
	}
	var err error
	if err = ff.AddMember(xd.fontsources["TeXGyreHeros-Regular"], frontend.FontWeight400, frontend.FontStyleNormal); err != nil {
		return err
	}
	if err = ff.AddMember(xd.fontsources["TeXGyreHeros-Italic"], frontend.FontWeight400, frontend.FontStyleItalic); err != nil {
		return err
	}
	if err = ff.AddMember(xd.fontsources["TeXGyreHeros-Bold"], frontend.FontWeight700, frontend.FontStyleNormal); err != nil {
		return err
	}
	if err = ff.AddMember(xd.fontsources["TeXGyreHeros-BoldItalic"], frontend.FontWeight700, frontend.FontStyleItalic); err != nil {
		return err
	}

	xd.fontsizes["text"] = [2]bag.ScaledPoint{tenpoint, twelvepoint}

	return nil
}
