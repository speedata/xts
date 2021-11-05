package core

import (
	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/node"
)

type paragraph struct {
	text string
}

type formatOptions struct {
	ff    *fontfamily
	width bag.ScaledPoint
}

func (p *paragraph) format(opts formatOptions) (*node.VList, error) {
	fnt, err := opts.ff.getFont(fontWeight400, fontStyleNormal, bag.MustSp("12pt"))
	if err != nil {
		logger.Errorf(err.Error())
		return nil, err
	}
	glyphs := fnt.Shape(p.text)

	var cur node.Node
	l, err := doc.LoadPatternFile("hyphenationpatterns/hyph-en-us.pat.txt")
	if err != nil {
		return nil, err
	}
	l.Name = "en-US"
	head := node.NewLangWithContents(&node.Lang{Lang: l})
	cur = head

	var lastglue node.Node
	for _, r := range glyphs {
		if r.Glyph == 32 {
			if lastglue == nil {
				g := node.NewGlue()
				g.Width = fnt.Space
				node.InsertAfter(head, cur, g)
				cur = g
				lastglue = g
			}
		} else {
			n := node.NewGlyph()
			n.Hyphenate = r.Hyphenate
			n.Codepoint = r.Codepoint
			n.Components = r.Components
			n.Font = fnt
			n.Width = r.Advance
			node.InsertAfter(head, cur, n)
			cur = n
			lastglue = nil
		}
	}
	if lastglue != nil && lastglue.Prev() != nil {
		p := lastglue.Prev()
		p.SetNext(nil)
		lastglue.SetPrev(nil)
	}

	settings := node.LinebreakSettings{
		HSize:      opts.width,
		LineHeight: 12 * bag.Factor,
	}

	hlist := node.Hpack(head)

	vlist := node.SimpleLinebreak(hlist, settings)

	return vlist, nil
}
