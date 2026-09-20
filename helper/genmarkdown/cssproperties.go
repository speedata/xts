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
		if link, ok := manualLinks[spec.Name]; ok {
			values += ", " + link
		}
		if spec.Note != "" {
			values += ". " + spec.Note
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
