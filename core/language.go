package core

import (
	"github.com/speedata/boxesandglue/backend/lang"
	"github.com/speedata/boxesandglue/frontend"
)

var languageMapping = map[string]string{
	"Ancient Greek":           "grc",
	"Armenian":                "hy",
	"Bahasa Indonesia":        "id",
	"Basque":                  "eu",
	"Bulgarian":               "bg",
	"Catalan":                 "ca",
	"Chinese":                 "zh",
	"Croatian":                "hr",
	"Czech":                   "cs",
	"Danish":                  "da",
	"Dutch":                   "nl",
	"English":                 "en_GB",
	"English (Great Britain)": "en_GB",
	"English (USA)":           "en_US",
	"Esperanto":               "eo",
	"Estonian":                "et",
	"Finnish":                 "fi",
	"French":                  "fr",
	"Galician":                "gl",
	"German":                  "de",
	"Greek":                   "el",
	"Gujarati":                "gu",
	"Hindi":                   "hi",
	"Hungarian":               "hu",
	"Icelandic":               "is",
	"Irish":                   "ga",
	"Italian":                 "it",
	"Kannada":                 "kn",
	"Kurmanji":                "ku",
	"Latvian":                 "lv",
	"Lithuanian":              "lt",
	"Malayalam":               "ml",
	"Norwegian Bokmål":        "nb",
	"Norwegian Nynorsk":       "nn",
	"Other":                   "--",
	"Polish":                  "pl",
	"Portuguese":              "pt",
	"Romanian":                "ro",
	"Russian":                 "ru",
	"Serbian":                 "sr",
	"Serbian (cyrillic)":      "sc",
	"Slovak":                  "sk",
	"Slovenian":               "sl",
	"Spanish":                 "es",
	"Swedish":                 "sv",
	"Turkish":                 "tr",
	"Ukrainian":               "uk",
	"Welsh":                   "cy",
}

// getLanguage returns a language object. Name can be a short name or a long name.
func (xd *xtsDocument) getLanguage(name string) (*lang.Lang, error) {
	if ln, ok := languageMapping[name]; ok {
		name = ln
	}
	return frontend.GetLanguage(name)
}
