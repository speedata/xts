package genmarkdown

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/speedata/xts/helper/commandsxml"
	"github.com/speedata/xts/helper/config"
)

var (
	refdir    string
	srcpath   string
	templates *template.Template
)

// DoThings creates the markdown documentation
func DoThings(cfg *config.Config) error {
	refdir = filepath.Join(cfg.Basedir(), "..", "xts-docs", "docs", "reference", ".ref")
	srcpath = filepath.Join(cfg.Basedir(), "doc")

	r, err := os.Open(filepath.Join(cfg.Basedir(), "doc", "commands-xml", "commands.xml"))
	if err != nil {
		return err
	}

	c, err := commandsxml.ReadCommandsFile(r)
	if err != nil {
		return err
	}
	funcMap := template.FuncMap{
		"childelements":  childelements,
		"parentelements": parentelements,
		"atttypeinfo":    atttypeinfo,
	}
	templates, err = template.New("").Funcs(funcMap).ParseFiles(filepath.Join(srcpath, "templates", "command.txt"))
	if err != nil {
		return err
	}

	referencedir := filepath.Join(cfg.Basedir(), "..", "xts-docs", "docs", "reference")
	for _, v := range c.Commands() {
		p := filepath.Join(referencedir, v.MDlink())
		// If the file doesn't exist, create it, or append to the file
		w, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			w.WriteString(fmt.Sprintf("---\ntitle: %s\n---\n", v.Name))
			w.WriteString(fmt.Sprintf("{%% include-markdown \".ref/%s\" %%}\n\n\n## See also\n", v.MDlink()))
			w.Close()
		}

		fullpath := filepath.Join(refdir, v.MDlink())
		builddoc(c, v, fullpath)
	}

	return nil
}

func parentelements(cmd *commandsxml.Command) string {
	var ret []string
	x := cmd.Parents()
	if len(x) == 0 {
		return "(none)"
	}

	for _, v := range x {
		ret = append(ret, fmt.Sprintf("[%s](%s)", v.Name, v.CmdLink()))
	}
	return strings.Join(ret, ", ")
}

func childelements(children []*commandsxml.Command) string {
	if len(children) == 0 {
		return string("(none)")
	}

	var ret []string
	for _, cmd := range children {
		ret = append(ret, fmt.Sprintf("[%s](%s)", cmd.Name, cmd.CmdLink()))
	}
	return strings.Join(ret, ", ")
}

func atttypeinfo(att *commandsxml.Attribute) string {
	atttypesEn := map[string]string{
		"boolean":                "yes or no",
		"xpath":                  `[XPath expressions](../../manual/xpath.md)`,
		"numberorlength":         "number or length",
		"numberlengthorstar":     "Number, length or *-numbers",
		"yesnolength":            "yes, no or length",
		"zerotohundred":          "0 up to 100",
		"zerohundredtwofivefive": "0 to 100 or 0 to 255",
	}
	ret := []string{}
	if att.Type != "" {
		if x, ok := atttypesEn[att.Type]; ok {
			ret = append(ret, x)
		} else {
			ret = append(ret, att.Type)
		}
	}
	if att.Optional {
		ret = append(ret, "optional")
	}
	return string(strings.Join(ret, ", "))
}

func builddoc(c *commandsxml.Commands, v *commandsxml.Command, fullpath string) {
	type sdata struct {
		Commands *commandsxml.Commands
		Command  *commandsxml.Command
	}

	f, err := os.Create(fullpath)
	if err != nil {
		panic(err)
	}
	err = templates.ExecuteTemplate(f, "command.txt", sdata{Commands: c, Command: v})
	if err != nil {
		fmt.Println(err)
	}
	f.Close()
	// wg.Done()
}
