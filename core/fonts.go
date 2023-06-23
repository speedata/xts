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
	fontsources := make(map[string]*frontend.FontSource)
	data := []struct {
		fontname string
		filename string
	}{
		{"TeXGyreHeros-Regular", "texgyreheros-regular.otf"},
		{"TeXGyreHeros-Italic", "texgyreheros-italic.otf"},
		{"TeXGyreHeros-Bold", "texgyreheros-bold.otf"},
		{"TeXGyreHeros-BoldItalic", "texgyreheros-bolditalic.otf"},
		{"CrimsonPro-Regular", "CrimsonPro-Regular.ttf"},
		{"CrimsonPro-Bold", "CrimsonPro-Bold.ttf"},
		{"CrimsonPro-BoldItalic", "CrimsonPro-BoldItalic.ttf"},
		{"CrimsonPro-Italic", "CrimsonPro-Italic.ttf"},
		{"CamingoCode-Bold", "CamingoCode-Bold.ttf"},
		{"CamingoCode-BoldItalic", "CamingoCode-BoldItalic.ttf"},
		{"CamingoCode-Italic", "CamingoCode-Italic.ttf"},
		{"CamingoCode-Regular", "CamingoCode-Regular.ttf"},
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
		fontsources[font.fontname] = &fs
	}
	var err error
	ff := xd.document.NewFontFamily("text")
	if err = ff.AddMember(fontsources["TeXGyreHeros-Regular"], frontend.FontWeight400, frontend.FontStyleNormal); err != nil {
		return err
	}
	if err = ff.AddMember(fontsources["TeXGyreHeros-Italic"], frontend.FontWeight400, frontend.FontStyleItalic); err != nil {
		return err
	}
	if err = ff.AddMember(fontsources["TeXGyreHeros-Bold"], frontend.FontWeight700, frontend.FontStyleNormal); err != nil {
		return err
	}
	if err = ff.AddMember(fontsources["TeXGyreHeros-BoldItalic"], frontend.FontWeight700, frontend.FontStyleItalic); err != nil {
		return err
	}

	ff = xd.document.NewFontFamily("sans")
	if err = ff.AddMember(fontsources["TeXGyreHeros-Regular"], frontend.FontWeight400, frontend.FontStyleNormal); err != nil {
		return err
	}
	if err = ff.AddMember(fontsources["TeXGyreHeros-Italic"], frontend.FontWeight400, frontend.FontStyleItalic); err != nil {
		return err
	}
	if err = ff.AddMember(fontsources["TeXGyreHeros-Bold"], frontend.FontWeight700, frontend.FontStyleNormal); err != nil {
		return err
	}
	if err = ff.AddMember(fontsources["TeXGyreHeros-BoldItalic"], frontend.FontWeight700, frontend.FontStyleItalic); err != nil {
		return err
	}

	ff = xd.document.NewFontFamily("serif")
	if err = ff.AddMember(fontsources["CrimsonPro-Regular"], frontend.FontWeight400, frontend.FontStyleNormal); err != nil {
		return err
	}
	if err = ff.AddMember(fontsources["CrimsonPro-Italic"], frontend.FontWeight400, frontend.FontStyleItalic); err != nil {
		return err
	}
	if err = ff.AddMember(fontsources["CrimsonPro-Bold"], frontend.FontWeight700, frontend.FontStyleNormal); err != nil {
		return err
	}
	if err = ff.AddMember(fontsources["CrimsonPro-BoldItalic"], frontend.FontWeight700, frontend.FontStyleItalic); err != nil {
		return err
	}

	ff = xd.document.NewFontFamily("monospace")
	if err = ff.AddMember(fontsources["CamingoCode-Regular"], frontend.FontWeight400, frontend.FontStyleNormal); err != nil {
		return err
	}
	if err = ff.AddMember(fontsources["CamingoCode-Italic"], frontend.FontWeight400, frontend.FontStyleItalic); err != nil {
		return err
	}
	if err = ff.AddMember(fontsources["CamingoCode-Bold"], frontend.FontWeight700, frontend.FontStyleNormal); err != nil {
		return err
	}
	if err = ff.AddMember(fontsources["CamingoCode-BoldItalic"], frontend.FontWeight700, frontend.FontStyleItalic); err != nil {
		return err
	}

	return nil
}
