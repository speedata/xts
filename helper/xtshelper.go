package main

import (
	"fmt"
	"log"
	"os"

	"github.com/speedata/optionparser"
	"github.com/speedata/xts/helper/changelog"
	"github.com/speedata/xts/helper/config"
	"github.com/speedata/xts/helper/genmarkdown"
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
	op.Command("doc", "Generate xts documentation (standalone)")
	op.Command("changelog", "changelog check VERSION verifies the changelog is ready for a release, changelog notes VERSION prints the release notes")
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
		if err = genmarkdown.DoThings(cfg); err != nil {
			log.Fatal(err)
		}
	case "genschema":
		err = genschema.DoThings(cfg)
		if err != nil {
			log.Fatal(err)
		}
	case "changelog":
		if len(op.Extra) != 3 {
			log.Fatal("usage: xtshelper changelog check|notes VERSION")
		}
		cl, err := changelog.Read(cfg)
		if err != nil {
			log.Fatal(err)
		}
		switch op.Extra[1] {
		case "check":
			if err = cl.Check(op.Extra[2]); err != nil {
				log.Fatal(err)
			}
		case "notes":
			notes, err := cl.ReleaseNotes(op.Extra[2])
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(notes)
		default:
			log.Fatal("usage: xtshelper changelog check|notes VERSION")
		}
	default:
		op.Help()
		os.Exit(-1)
	}
}
