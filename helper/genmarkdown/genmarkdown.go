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
	refdir = filepath.Join(cfg.Basedir(), "doc", "manual", "ref")
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

	referencedir := filepath.Join(cfg.Basedir(), "doc", "manual", "content", "reference", "commands")
	if err = os.MkdirAll(referencedir, 0o755); err != nil {
		return err
	}
	if err = os.MkdirAll(refdir, 0o755); err != nil {
		return err
	}
	for _, v := range c.Commands() {
		p := filepath.Join(referencedir, v.MDlink())
		w, err := os.Create(p)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "---\ntype: docs\nlinktitle: %s\n---\n", v.Name)
		fmt.Fprintf(w, "{{%% include \"%s\" %%}}\n\n\n## See also\n", v.MDlink())
		w.Close()

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
		"xpath":                  `[XPath expressions](/manual/data-processing/xpath)`,
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
