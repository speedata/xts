package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/speedata/optionparser"
	"github.com/speedata/xts/helper/config"
	"github.com/speedata/xts/helper/db2html"
	"github.com/speedata/xts/helper/genadoc"
	"github.com/speedata/xts/helper/genschema"
)

var (
	version string
	basedir string
)

// sitedoc: true: needs webserver, false: standalone HTML files
func makedoc(cfg *config.Config, sitedoc bool) error {
	var err error
	err = os.RemoveAll(filepath.Join(cfg.Builddir, "manual"))
	if err != nil {
		return err
	}

	for _, lang := range []string{"en"} {
		err = genadoc.DoThings(cfg, lang)
		if err != nil {
			return err
		}
		var manualfile string
		switch lang {
		case "en":
			manualfile = "xts"
		case "de":
			manualfile = "publisherhandbuch"
		}
		err = db2html.DoThings(cfg, manualfile, sitedoc)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	cfg := config.NewConfig(basedir, version)

	op := optionparser.NewOptionParser()
	op.Command("genschema", "Generate schema (in language de, en and schema xsd and rng)")
	op.Command("doc", "Generate speedata Publisher documentation (standalone)")
	err := op.Parse()
	if err != nil {
		log.Fatal(err)
	}

	var command string
	if len(op.Extra) > 0 {
		command = op.Extra[0]
	} else {
		op.Help()
		os.Exit(-1)
	}
	switch command {
	case "doc":
		err = makedoc(cfg, false)
		if err != nil {
			log.Fatal(err)
		}
	case "genschema":
		err = genschema.DoThings(cfg)
		if err != nil {
			log.Fatal(err)
		}

	default:
		op.Help()
		os.Exit(-1)
	}
}
