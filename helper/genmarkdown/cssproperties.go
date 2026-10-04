package genmarkdown

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/boxesandglue/htmlbag"
	"github.com/speedata/xts/helper/config"
)

// manualLinks adds a pointer into the XTS manual to the value description of a
// property. htmlbag knows nothing about the manual, so the links live here.
var manualLinks = map[string]string{
	"line-height":   "see [Fonts](/manual/text-and-styling/fonts#font-size-and-line-height)",
	"-bag-bookmark": "see [Bookmarks](/manual/advanced/pdf-options#bookmarks)",
}

// noteOverrides replaces htmlbag's note on a property where XTS behaves
// differently, so the reference does not promise what htmlbag's pagination
// does but XTS doesn't.
var noteOverrides = map[string]string{
	"page-break-before":  "In a `Flow`, `column` continues in the next frame of the area, `page` on the next page, `left` and `verso` on the next left (even) page, `right` and `recto` on the next right (odd) one. A forced break is taken at any depth; on a block's first child it applies to the block itself",
	"page-break-after":   "Acts in a `Flow` as `break-before` does, `avoid` keeps a block with the next one, such as a heading with its first paragraph. A forced break is taken at any depth; on a block's last child it applies to the block itself",
	"page-break-inside":  "`avoid` keeps a paragraph or a container in one frame of a `Flow` unless it is taller than the frame. XTS does not split a table row: a row that does not fit moves to the next frame whole",
	"widows":             "Acts in a `Flow`; a paragraph placed with `PlaceObject` is never split",
	"orphans":            "Acts in a `Flow`, as with `widows`",
	"vertical-align":     "",
	"-bag-leading-model": "",
}

// valueOverrides replaces htmlbag's value description where it names Go API
// that XTS does not offer to the layout author.
var valueOverrides = map[string]string{
	"-bag-leading-model": "`half` (CSS line boxes, the default), `trailing` (TeX style), or a name that the Go program running XTS registered in `XTSConfig.LineModels`",
	"-bag-line-breaker":  "`auto` (Knuth-Plass, the default) or a name that the Go program running XTS registered in `XTSConfig.Breakers`",
}

type cssProperty struct {
	Names   string
	Values  string
	Example string
}

type cssGroup struct {
	Title      string
	Properties []cssProperty
}

func cssGroupFor(g htmlbag.PropertyGroup) cssGroup {
	grp := cssGroup{Title: string(g)}
	for _, spec := range htmlbag.PropertiesInGroup(g) {
		names := make([]string, 0, len(spec.Aliases)+1)
		for _, n := range spec.Names() {
			names = append(names, "`"+n+"`")
		}
		values := spec.Values
		if o, ok := valueOverrides[spec.Name]; ok {
			values = o
		}
		if link, ok := manualLinks[spec.Name]; ok {
			values += ", " + link
		}
		note := spec.Note
		if o, ok := noteOverrides[spec.Name]; ok {
			note = o
		}
		if note != "" {
			values += ". " + note
		}
		grp.Properties = append(grp.Properties, cssProperty{
			Names:   strings.Join(names, ", "),
			Values:  values,
			Example: spec.Example,
		})
	}
	return grp
}

// writeCSSProperties renders the CSS reference page of the manual from
// htmlbag.Properties and the template doc/templates/css-properties.txt, which
// holds the parts that are specific to XTS (selectors, at-rule prose).
func writeCSSProperties(cfg *config.Config) error {
	tmpl, err := template.ParseFiles(filepath.Join(cfg.Basedir(), "doc", "templates", "css-properties.txt"))
	if err != nil {
		return err
	}
	data := struct {
		Elements []cssGroup
		Page     cssGroup
		FontFace cssGroup
		Color    cssGroup
	}{}
	for _, g := range htmlbag.PropertyGroups {
		switch g {
		case htmlbag.GroupPage:
			data.Page = cssGroupFor(g)
		case htmlbag.GroupFontFace:
			data.FontFace = cssGroupFor(g)
		case htmlbag.GroupColor:
			data.Color = cssGroupFor(g)
		default:
			data.Elements = append(data.Elements, cssGroupFor(g))
		}
	}
	for _, grp := range append(data.Elements, data.Page, data.FontFace, data.Color) {
		for _, p := range grp.Properties {
			if strings.Contains(p.Values, "|") || strings.Contains(p.Example, "|") {
				return fmt.Errorf("css-properties: %s contains a pipe, which breaks the Markdown table", p.Names)
			}
		}
	}
	out := filepath.Join(cfg.Basedir(), "doc", "manual", "content", "reference", "css-properties.md")
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.ExecuteTemplate(f, "css-properties.txt", data)
}
