package core

import (
	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/frontend"
)

type textformat struct {
	halignment frontend.HorizontalAlignment
}

func (xd *xtsDocument) defaultTextformats() {
	xd.defineTextformat("text", textformat{halignment: frontend.HAlignJustified})
	xd.defineTextformat("left", textformat{halignment: frontend.HAlignLeft})
	xd.defineTextformat("right", textformat{halignment: frontend.HAlignRight})
	xd.defineTextformat("center", textformat{halignment: frontend.HAlignCenter})
}

func (xd *xtsDocument) defineTextformat(name string, tf textformat) {
	if _, ok := xd.textformats[name]; ok {
		bag.Logger.Infof("Redefine text format %q", name)
	} else {
		bag.Logger.Infof("Define text format %q", name)
	}
	xd.textformats[name] = tf
}
