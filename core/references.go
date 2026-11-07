package core

import (
	"encoding/xml"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
	"strconv"

	"github.com/boxesandglue/boxesandglue/backend/node"
	"github.com/speedata/goxpath"
)

// getNumDest returns a new start stop node with a ActionDest action and a
// distinct numeric Value
func (xd *xtsDocument) getNumDest() *node.StartStop {
	dest := node.NewStartStop()
	dest.Action = node.ActionDest
	dest.Value = xd.nextDestination
	xd.nextDestination++
	return dest
}

// getNameDest returns a new start stop node with a ActionDest action and a
// title set to the name.
func getNameDest(name string) *node.StartStop {
	dest := node.NewStartStop()
	dest.Action = node.ActionDest
	dest.Value = name
	return dest
}

type marker struct {
	name       string
	pdftarget  bool
	pagenumber int
	id         int // a per page uid
}

type mapmarker map[string]marker

func (mm *mapmarker) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		tok, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			m := marker{}
			for _, attr := range t.Attr {
				switch attr.Name.Local {
				case "name":
					m.name = attr.Value
				case "id":
					m.id, err = strconv.Atoi(attr.Value)
					if err != nil {
						return err
					}
				case "page":
					m.pagenumber, err = strconv.Atoi(attr.Value)
					if err != nil {
						return err
					}
				case "pdftarget":
					m.pdftarget, err = strconv.ParseBool(attr.Value)
					if err != nil {
						return err
					}
				}
			}
			(*mm)[m.name] = m
		}
	}
	return nil
}

// MarshalXML writes the marks sorted by name, so that two runs with the same
// result write the same file.
func (mm mapmarker) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeToken(start)
	for _, name := range slices.Sorted(maps.Keys(mm)) {
		v := mm[name]
		mStart := xml.StartElement{Name: xml.Name{Local: "mark"}}
		mStart.Attr = append(mStart.Attr, xml.Attr{Name: xml.Name{Local: "name"}, Value: v.name})
		mStart.Attr = append(mStart.Attr, xml.Attr{Name: xml.Name{Local: "page"}, Value: fmt.Sprintf("%d", v.pagenumber)})
		mStart.Attr = append(mStart.Attr, xml.Attr{Name: xml.Name{Local: "id"}, Value: fmt.Sprintf("%d", v.id)})
		mStart.Attr = append(mStart.Attr, xml.Attr{Name: xml.Name{Local: "pdftarget"}, Value: fmt.Sprintf("%t", v.pdftarget)})
		e.EncodeToken(mStart)
		e.EncodeToken(mStart.End())
	}
	return e.EncodeToken(start.End())
}

func (m marker) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(m.name, start)
}

type auxfile struct {
	Marker   mapmarker `xml:"marker"`
	LastPage int       `xml:"lastpage"`
}

func (d *xtsDocument) writeAuxXML() error {
	var aux auxfile
	aux.LastPage = d.currentPagenumber

	f, err := os.Create(d.cfg.Jobname + "-aux.xml")
	if err != nil {
		return err
	}
	defer f.Close()
	aux.Marker = d.marker

	data, err := xml.MarshalIndent(aux, "", "  ")
	if err != nil {
		return err
	}

	if _, err = f.Write(data); err != nil {
		return err
	}

	return nil
}

func (d *xtsDocument) readAuxFile() (*auxfile, error) {
	data, err := os.ReadFile(d.cfg.Jobname + "-aux.xml")
	if err != nil {
		// OK if file not found, just return
		return nil, nil
	}

	var aux auxfile
	aux.Marker = make(mapmarker)
	if err = xml.Unmarshal(data, &aux); err != nil {
		return nil, err
	}
	d.data.SetVariable("_lastpage", goxpath.Sequence{aux.LastPage})
	return &aux, nil
}

// getMarker returns the mark of this run or else the one of the previous run.
// It notes the answer for outOfDate.
func (d *xtsDocument) getMarker(name string) (marker, bool) {
	m, ok := d.marker[name]
	if !ok && d.aux != nil {
		m, ok = d.aux.Marker[name]
	}
	d.marksRead[markRead{name, m.pagenumber}] = true
	return m, ok
}

// markRead is a page number that the layout got for a mark.
type markRead struct {
	name string
	page int
}

// outOfDate returns the references that told the layout something else than
// the finished run knows, sorted. The document may be wrong in these places, so
// it needs another run. A mark that nobody read does not count: it moves
// without any effect. Neither does a mark that does not exist in the end, the
// layout gets the same 0 in any run.
func (d *xtsDocument) outOfDate() []string {
	var stale []string
	for r := range d.marksRead {
		if d.marker[r.name].pagenumber != r.page {
			stale = append(stale, r.name)
		}
	}
	if d.lastPageRead && d.lastPageOfPreviousRun() != d.currentPagenumber {
		stale = append(stale, "sd:last-page-number()")
	}
	slices.Sort(stale)
	return slices.Compact(stale)
}

// lastPageOfPreviousRun returns the number of pages of the previous run and 0
// before the first one.
func (d *xtsDocument) lastPageOfPreviousRun() int {
	if d.aux == nil {
		return 0
	}
	return d.aux.LastPage
}
