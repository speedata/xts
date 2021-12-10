package main

import (
	"log"
	"os"

	"github.com/speedata/optionparser"
	"github.com/speedata/xts/helper/config"
	"github.com/speedata/xts/helper/genschema"
)

var (
	version string
	basedir string
)

func main() {
	cfg := config.NewConfig(basedir, version)

	op := optionparser.NewOptionParser()
	op.Command("genschema", "Generate schema (in language de, en and schema xsd and rng)")
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
