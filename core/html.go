package core

import (
	"fmt"
	"strings"

	"github.com/speedata/bagme/document"
	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/node"
	"github.com/speedata/boxesandglue/frontend"
	"github.com/speedata/goxml"
)

func (xd *xtsDocument) parseHTML(elt *goxml.Element) (frontend.FormatToVList, error) {
	ftv := func(wd bag.ScaledPoint, opts ...frontend.TypesettingOption) (*node.VList, error) {
		str := elt.ToXML()
		css := []string{}
		options := &frontend.Options{}
		for _, opt := range opts {
			opt(options)
		}
		if options.Fontfamily != nil {
			css = append(css, fmt.Sprintf(`body { font-family: %s; }`, options.Fontfamily.Name))
		}
		if fs := options.Fontsize; fs != 0 {
			css = append(css, fmt.Sprintf(`body { font-size: %spt; }`, fs))
		}
		if fs := options.Leading; fs != 0 {
			css = append(css, fmt.Sprintf(`body { line-height: %spt; }`, fs))
		}
		d := document.NewWithFrontend(xd.document, xd.datacss)
		d.AddCSS(strings.Join(css, " "))
		te, err := d.ParseHTML(str)
		if err != nil {
			return nil, err
		}
		vl, err := d.CreateVlist(te, wd)
		if err != nil {
			return nil, err
		}

		return vl, err
	}
	return ftv, nil
}
