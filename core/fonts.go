package core

import (
	"github.com/speedata/boxesandglue/frontend"
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

	ff.AddMember(xd.fontsources["TeXGyreHeros-Regular"], frontend.FontWeight400, frontend.FontStyleNormal)
	ff.AddMember(xd.fontsources["TeXGyreHeros-Italic"], frontend.FontWeight400, frontend.FontStyleItalic)
	ff.AddMember(xd.fontsources["TeXGyreHeros-Bold"], frontend.FontWeight700, frontend.FontStyleNormal)
	ff.AddMember(xd.fontsources["TeXGyreHeros-BoldItalic"], frontend.FontWeight700, frontend.FontStyleItalic)
	return nil
}
