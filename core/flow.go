package core

import (
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"strings"

	"github.com/boxesandglue/boxesandglue/backend/bag"
	"github.com/boxesandglue/boxesandglue/backend/node"
	"github.com/boxesandglue/boxesandglue/frontend"
	"github.com/boxesandglue/htmlbag"
	"github.com/speedata/goxml"
	xpath "github.com/speedata/goxpath"
	"golang.org/x/net/html"
)

// flowEnd is where the last Flow on a page ended, with the margin below its
// last block still open, for a Flow that follows it directly.
type flowEnd struct {
	y, marginAfter bag.ScaledPoint
}

// flowRegions hands htmlbag the free bands of an area's frames, page after
// page, and places what it fills.
type flowRegions struct {
	xd   *xtsDocument
	name string
	area *area
	// start is where the previous Flow ended, when this one follows it
	// directly; startPage is the page it ended on.
	start     *flowEnd
	startPage *page
	// next is the frame row the search for the next band starts at.
	next    coord
	started bool
	// The current region: its frame row, top and left from the page's top
	// left corner, and size.
	row           coord
	top, left     bag.ScaledPoint
	width, height bag.ScaledPoint
	// filled counts the regions handed back.
	filled int
	// end is the lowest bottom of the regions on endPage, the page the flow
	// ends on, and endMargin the margin below it.
	endPage        *page
	end, endMargin bag.ScaledPoint
	// marks holds the Mark and Bookmark nodes waiting for the flow child of
	// that Fragment.Index to begin.
	marks map[int][]*node.StartStop
}

// placeMarks ships the mark and bookmark nodes ss out at y, measured from the
// page top, and x on the current page, so they take its number and position.
func (r *flowRegions) placeMarks(ss []*node.StartStop, x, y bag.ScaledPoint) {
	for _, s := range ss {
		if s.Attributes != nil {
			if _, ok := s.Attributes["page"]; ok {
				s.Attributes["page"] = r.xd.currentPage
			}
		}
		r.xd.currentPage.outputAbsolute(x, y, node.Vpack(s))
	}
}

// absRow is the page row of the frame row row.
func absRow(a *area, row coord) coord {
	return row + a.frame[a.currentFrame].row - 1
}

// rowFree reports whether no cell of the current frame's row row is
// allocated.
func (g *grid) rowFree(a *area, row coord) bool {
	f := a.frame[a.currentFrame]
	for c := f.col; c < f.col+f.width; c++ {
		if g.allocatedBlocks.allocValue(c, absRow(a, row)) > 0 {
			return false
		}
	}
	return true
}

// band returns the first run of free rows of the current frame from row on,
// first to last, or ok false when there is none.
func (g *grid) band(a *area, row coord) (first, last coord, ok bool) {
	f := a.frame[a.currentFrame]
	for first = max(row, 1); first <= f.height; first++ {
		if g.rowFree(a, first) {
			break
		}
	}
	if first > f.height {
		return 0, 0, false
	}
	for last = first; last < f.height && g.rowFree(a, last+1); last++ {
	}
	return first, last, true
}

// nextFrame moves on to the area's next frame, or to the first frame of a new
// page after the last one.
func (r *flowRegions) nextFrame() error {
	if r.area.currentFrame+1 < len(r.area.frame) {
		r.area.currentFrame++
		r.next = 1
		return nil
	}
	return r.nextPage()
}

// nextPage starts a new page and goes on in the first frame of its area of
// the same name.
func (r *flowRegions) nextPage() error {
	xd := r.xd
	clearPage(xd)
	xd.setupPage()
	a, ok := xd.currentGrid.areas[r.name]
	if !ok {
		return fmt.Errorf("area %s not found on page %d", r.name, xd.currentPage.pagenumber)
	}
	r.area = a
	r.area.currentFrame = 0
	r.next = 1
	return nil
}

// Next returns the next free band: below the last region in its frame, in
// the next frame, or on the next page. A forced break to a column takes the
// next frame, any other the next page (left and right one of that parity).
func (r *flowRegions) Next(brk string) (htmlbag.Region, error) {
	xd := r.xd
	if r.started {
		var err error
		switch brk {
		case "":
		case "column":
			err = r.nextFrame()
		default:
			if err = r.nextPage(); err == nil {
				if odd := xd.currentPage.pagenumber%2 == 1; (brk == "left" && odd) || (brk == "right" && !odd) {
					err = r.nextPage()
				}
			}
		}
		if err != nil {
			return htmlbag.Region{}, err
		}
	}
	r.started = true
	for tries := 0; ; tries++ {
		g := xd.currentGrid
		first, last, ok := g.band(r.area, r.next)
		if !ok {
			if tries > 2*len(r.area.frame)+2 {
				return htmlbag.Region{}, fmt.Errorf("area %s has no free row", r.name)
			}
			if err := r.nextFrame(); err != nil {
				return htmlbag.Region{}, err
			}
			continue
		}
		// A band starts at the exact bottom of what ends in the row above
		// it, not on the next whole row. Where that is the end of a Flow
		// this one follows directly, the first region collapses the margin
		// still open with the first block's; the others start below it,
		// level with the first block.
		top, continues := r.bandTop(first)
		bottom := g.posY(last, r.area) + g.gridHeight
		if bottom <= top {
			r.next = last + 1
			continue
		}
		r.row, r.top = first, top
		r.left = g.posX(1, r.area)
		r.width = g.width(r.area.frame[r.area.currentFrame].width)
		r.height = bottom - top
		r.next = last + 1
		reg := htmlbag.Region{
			Width:   r.width,
			Height:  r.height,
			PageNum: xd.currentPage.pagenumber,
			Left:    r.left,
			Top:     xd.currentPage.pageHeight - top,
		}
		if continues && r.filled == 0 {
			reg.MarginBefore = r.start.marginAfter
		}
		return reg, nil
	}
}

// bandTop is the top of a band that starts in the frame row first, and
// whether it continues the Flow this one follows directly.
func (r *flowRegions) bandTop(first coord) (top bag.ScaledPoint, continues bool) {
	g := r.xd.currentGrid
	top = g.posY(first, r.area)
	if y, ok := g.exactTop(r.area, first); ok {
		top = y
		if s := r.start; s != nil && r.xd.currentPage == r.startPage && y == s.y {
			continues = true
			if r.filled > 0 {
				top += s.marginAfter
			}
		}
	}
	return top, continues
}

// startTop is where a flow that has asked for no region would have started:
// the top of the first band of the current frame, or of its next row when
// the frame has none.
func (r *flowRegions) startTop() bag.ScaledPoint {
	if first, _, ok := r.xd.currentGrid.band(r.area, r.next); ok {
		top, _ := r.bandTop(first)
		return top
	}
	return r.xd.currentGrid.posY(r.next, r.area)
}

// Filled places the region's box and allocates the rows it reaches into.
func (r *flowRegions) Filled(f htmlbag.Filled) error {
	xd := r.xd
	g := xd.currentGrid
	for _, fr := range f.Fragments {
		if ss, ok := r.marks[fr.Index]; ok && !fr.Continued {
			r.placeMarks(ss, r.left, r.top+fr.Top)
			delete(r.marks, fr.Index)
		}
	}
	if f.Box != nil && f.Used > 0 {
		xd.currentPage.outputAbsolute(r.left, r.top, f.Box)
		ht := r.top - g.posY(r.row, r.area) + f.Used
		g.allocate(1, r.row, r.area, r.width, ht)
		r.area.SetCurrentRow(r.row + g.heightToRows(ht))
		r.area.SetCurrentCol(1)
	}
	r.filled++
	if r.endPage != xd.currentPage {
		r.endPage, r.end = xd.currentPage, 0
	}
	if bottom := r.top + f.Used; bottom >= r.end {
		r.end, r.endMargin = bottom, f.MarginAfter
	} else {
		r.endMargin = 0
	}
	return nil
}

func cmdFlow(xd *xtsDocument, layoutelt *goxml.Element) (xpath.Sequence, error) {
	attValues := &struct {
		Area   string
		Bottom string
	}{}
	if err := getXMLAttributes(xd, layoutelt, attValues); err != nil {
		return nil, err
	}
	if attValues.Area == "" {
		attValues.Area = defaultAreaName
	}
	if xd.currentSlate != nil {
		return nil, newTypesettingError("Flow", layoutelt.Line, "a flow is not possible in a slate")
	}
	xd.setupPage()
	area, ok := xd.currentGrid.areas[attValues.Area]
	if !ok {
		return nil, newTypesettingErrorf("Flow", layoutelt.Line, "area %s not found", attValues.Area)
	}
	if xd.inFlow {
		return nil, newTypesettingError("Flow", layoutelt.Line, "a flow inside a flow is not possible")
	}
	// The blocks are laid out by the flow, so the frame's width is the
	// widest a child may be.
	xd.store["maxwidth"] = int(area.frame[area.currentFrame].width)
	xd.inFlow = true
	xd.flowOrigin = map[node.Node]*goxml.Element{}
	seq, err := dispatch(xd, layoutelt)
	origin := xd.flowOrigin
	xd.inFlow, xd.flowOrigin = false, nil
	if err != nil {
		return nil, err
	}
	body := &html.Node{Data: "body", Type: html.ElementNode}
	// pending are the marks and bookmarks waiting for the next child, before
	// holds them by the child they wait for.
	var pending []*node.StartStop
	before := map[*html.Node][]*node.StartStop{}
	appendChild := func(n *html.Node) {
		if len(pending) > 0 {
			before[n] = append(before[n], pending...)
			pending = nil
		}
		body.AppendChild(n)
	}
	for _, itm := range seq {
		switch t := itm.(type) {
		case *html.Node:
			appendChild(t)
		case *goxml.Element:
			appendChild(goxmlToHTMLNode(t))
		case goxml.Element:
			appendChild(goxmlToHTMLNode(&t))
		case marker:
			pending = append(pending, xd.markDest(t))
		case *node.StartStop:
			// From an Action with a Mark, or a Bookmark.
			pending = append(pending, t)
		default:
			if n, ok := itm.(node.Node); ok && origin[n] != nil {
				elt := origin[n]
				slog.Warn(fmt.Sprintf("%s (line %d): cannot be part of a flow and is left out", elt.Name, elt.Line))
			} else {
				slog.Warn(fmt.Sprintf("Flow (line %d): a %T cannot be part of a flow and is left out", layoutelt.Line, itm))
			}
		}
	}
	doc := &html.Node{Type: html.DocumentNode}
	root := &html.Node{Data: "html", Type: html.ElementNode}
	root.AppendChild(&html.Node{Data: "head", Type: html.ElementNode})
	root.AppendChild(body)
	doc.AppendChild(root)
	trailing := pending
	te, err := xd.cssbuilder.ParseHTMLFromNode(doc)
	if err != nil {
		return nil, newTypesettingError("Flow", layoutelt.Line, err.Error())
	}

	r := &flowRegions{xd: xd, name: attValues.Area, area: area, next: area.CurrentRow()}
	if len(before) > 0 {
		r.marks = map[int][]*node.StartStop{}
		index := flowChildIndices(body, te)
		for c := body.FirstChild; c != nil; c = c.NextSibling {
			if ss, ok := before[c]; ok {
				if i, ok := index[c]; ok {
					r.marks[i] = append(r.marks[i], ss...)
				} else {
					// A child that makes no block, such as an empty one: the
					// marks wait for the next one.
					if n := c.NextSibling; n != nil {
						before[n] = append(ss, before[n]...)
					} else {
						trailing = append(ss, trailing...)
					}
				}
			}
		}
	}
	if area.CurrentCol() != 1 {
		// The row is taken in part already.
		r.next++
	}
	if fe := xd.currentGrid.flowEnd; fe != nil {
		r.start, r.startPage = fe, xd.currentPage
	}
	if err := xd.cssbuilder.FlowText(te, r); err != nil {
		return nil, newTypesettingError("Flow", layoutelt.Line, err.Error())
	}

	// Marks after the last child, or before a child that never began, take
	// the page where the flow ends.
	var rest []*node.StartStop
	for _, i := range slices.Sorted(maps.Keys(r.marks)) {
		rest = append(rest, r.marks[i]...)
	}
	rest = append(rest, trailing...)
	if len(rest) > 0 {
		end := r.end
		if !r.started {
			end = r.startTop()
		}
		r.placeMarks(rest, xd.currentGrid.posX(1, r.area), end)
	}

	if r.endPage == xd.currentPage {
		g := xd.currentGrid
		g.flowEnd = &flowEnd{y: r.end, marginAfter: r.endMargin}
		// An object placed on the page after the flow starts on the next
		// whole row, below the longest frame.
		pa := g.areas[pageAreaName]
		if row := max(1, g.heightToRows(r.end-g.marginTop)); pa.CurrentRow() <= row {
			pa.SetCurrentRow(row + 1)
			pa.SetCurrentCol(1)
		}
	}
	if attValues.Bottom != "" {
		end := r.end
		if !r.started {
			// FlowText asks for no region when there is no block.
			end = r.startTop()
		}
		// From the top of the page's grid, as (sd:current-row() - 1) times
		// the row height measures a row.
		xd.data.SetVariable(attValues.Bottom, xpath.Sequence{(end - xd.currentGrid.marginTop).ToPT()})
	}
	return nil, nil
}

// flowChildIndices maps the children of body to the Fragment.Index of the flow
// child each one becomes in te, the Text htmlbag made of the document. An
// element htmlbag sets as a block is a child of its own; text and inline
// elements between blocks make one anonymous child together. A child that
// makes no block, such as an empty paragraph, has no index.
func flowChildIndices(body *html.Node, te *frontend.Text) map[*html.Node]int {
	// The body Text, inside the html one.
	for range 2 {
		if len(te.Items) != 1 {
			break
		}
		inner, ok := te.Items[0].(*frontend.Text)
		if !ok {
			break
		}
		te = inner
	}
	// htmlbag counts the body's Texts, but not the whitespace between blocks.
	type item struct {
		tag   string
		index int
	}
	var items []item
	for _, itm := range te.Items {
		t, ok := itm.(*frontend.Text)
		if !ok {
			continue
		}
		tag, _ := t.Settings[frontend.SettingDebug].(string)
		if tag == "" && whitespaceText(t) {
			continue
		}
		items = append(items, item{tag, len(items)})
	}
	index := map[*html.Node]int{}
	j, inRun := 0, false
	for c := body.FirstChild; c != nil && j <= len(items); c = c.NextSibling {
		if c.Type == html.ElementNode && j < len(items) && items[j].tag == c.Data {
			index[c] = items[j].index
			j++
			inRun = false
			continue
		}
		if c.Type == html.TextNode && strings.TrimSpace(c.Data) == "" {
			continue
		}
		if c.Type == html.ElementNode && isBlockElement(c.Data) {
			// A block htmlbag dropped.
			inRun = false
			continue
		}
		if !inRun {
			if j >= len(items) || items[j].tag != "" {
				continue
			}
			inRun = true
			j++
		}
		index[c] = items[j-1].index
	}
	return index
}

// whitespaceText reports whether t holds nothing but whitespace.
func whitespaceText(t *frontend.Text) bool {
	for _, itm := range t.Items {
		switch v := itm.(type) {
		case string:
			if strings.TrimSpace(v) != "" {
				return false
			}
		case *frontend.Text:
			if !whitespaceText(v) {
				return false
			}
		default:
			return false
		}
	}
	return true
}

// isBlockElement reports whether an element of a Flow's children is a block.
func isBlockElement(tag string) bool {
	switch tag {
	case "p", "div", "table", "ul", "ol", "li", "h1", "h2", "h3", "h4", "h5", "h6", "pre", "blockquote", "section", "article", "hr", "dl", "figure":
		return true
	}
	return false
}
