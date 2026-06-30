// Package genschema creates Relax NG and XSD schema files for English and German
package genschema

import (
	"os"
	"path/filepath"

	"github.com/speedata/xts/helper/config"
)

const (
	// SDNAMESPACE is the speedata XTS layout rules namespace
	SDNAMESPACE string = "urn:speedata.de/2021/xts/en"
	// FNNAMESPACE is the speedata XTS function namespace
	FNNAMESPACE string = "urn:speedata.de/2021/xtsfunctions/en"
	// XHTMLNAMESPACE is the XHTML namespace for literal HTML elements in XTS layouts
	XHTMLNAMESPACE string = "http://www.w3.org/1999/xhtml"
)

// DoThings creates Relax NG and XSD schema files for »en« and »de«. Both
// schemas are generated programmatically from commands.xml; no external
// converter (formerly trang.jar) is required.
func DoThings(cfg *config.Config) error {
	basedir := cfg.Basedir()
	c, err := readCommandsFile(basedir)
	if err != nil {
		return err
	}

	schemas := []struct {
		path string
		gen  func() ([]byte, error)
	}{
		{filepath.Join(basedir, "schema", "layoutschema-en.rng"), func() ([]byte, error) { return genRelaxNGSchema(c, "en", true) }},
		{filepath.Join(basedir, "schema", "layoutschema-de.rng"), func() ([]byte, error) { return genRelaxNGSchema(c, "de", true) }},
		{filepath.Join(basedir, "schema", "layoutschema-en.xsd"), func() ([]byte, error) { return genXSDSchema(c, "en") }},
		{filepath.Join(basedir, "schema", "layoutschema-de.xsd"), func() ([]byte, error) { return genXSDSchema(c, "de") }},
	}

	for _, s := range schemas {
		buf, err := s.gen()
		if err != nil {
			return err
		}
		if err = os.WriteFile(s.path, buf, 0o644); err != nil {
			return err
		}
	}
	return nil
}
