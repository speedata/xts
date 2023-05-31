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
	xd.fontaliases = map[string]string{
		"sans-regular":         "TeXGyreHeros-Regular",
		"sans-italic":          "TeXGyreHeros-Italic",
		"sans-bold":            "TeXGyreHeros-Bold",
		"sans-bolditalic":      "TeXGyreHeros-BoldItalic",
		"serif-regular":        "CrimsonPro-Regular",
		"serif-bold":           "CrimsonPro-Bold",
		"serif-bolditalic":     "CrimsonPro-BoldItalic",
		"serif-italic":         "CrimsonPro-Italic",
		"monospace-bold":       "CamingoCode-Bold",
		"monospace-bolditalic": "CamingoCode-BoldItalic",
		"monospace-italic":     "CamingoCode-Italic",
		"monospace-regular":    "CamingoCode-Regular",
	}
	xd.fontsources = make(map[string]*frontend.FontSource)
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
		xd.fontsources[font.fontname] = &fs
	}
	var err error
	ff := xd.document.NewFontFamily("text")
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

	ff = xd.document.NewFontFamily("sans")
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

	ff = xd.document.NewFontFamily("serif")
	if err = ff.AddMember(xd.fontsources["CrimsonPro-Regular"], frontend.FontWeight400, frontend.FontStyleNormal); err != nil {
		return err
	}
	if err = ff.AddMember(xd.fontsources["CrimsonPro-Italic"], frontend.FontWeight400, frontend.FontStyleItalic); err != nil {
		return err
	}
	if err = ff.AddMember(xd.fontsources["CrimsonPro-Bold"], frontend.FontWeight700, frontend.FontStyleNormal); err != nil {
		return err
	}
	if err = ff.AddMember(xd.fontsources["CrimsonPro-BoldItalic"], frontend.FontWeight700, frontend.FontStyleItalic); err != nil {
		return err
	}

	ff = xd.document.NewFontFamily("monospace")
	if err = ff.AddMember(xd.fontsources["CamingoCode-Regular"], frontend.FontWeight400, frontend.FontStyleNormal); err != nil {
		return err
	}
	if err = ff.AddMember(xd.fontsources["CamingoCode-Italic"], frontend.FontWeight400, frontend.FontStyleItalic); err != nil {
		return err
	}
	if err = ff.AddMember(xd.fontsources["CamingoCode-Bold"], frontend.FontWeight700, frontend.FontStyleNormal); err != nil {
		return err
	}
	if err = ff.AddMember(xd.fontsources["CamingoCode-BoldItalic"], frontend.FontWeight700, frontend.FontStyleItalic); err != nil {
		return err
	}

	xd.fontsizes["text"] = [2]bag.ScaledPoint{tenpoint, twelvepoint}

	return nil
}
