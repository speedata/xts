package core

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/node"
	"github.com/speedata/goxpath"
)

// getNumDest returns a new start stop node with a ActionDest action and a
// distinct numeric Value
func getNumDest() *node.StartStop {
	dest := node.NewStartStop()
	dest.Action = node.ActionDest
	dest.Value = <-destinationNumbers
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
	append     bool
	pdftarget  bool
	pagenumber int
	id         int // a per page uid
	shiftup    bag.ScaledPoint
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
				}
			}
			(*mm)[m.name] = m
			// fmt.Println(tok, err)
		}
	}
	return nil
}
func (mm mapmarker) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeToken(start)
	for _, v := range mm {
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

	f, err := os.Create(d.jobname + "-aux.xml")
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
	data, err := os.ReadFile(d.jobname + "-aux.xml")
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

func (d *xtsDocument) getMarker(name string) (marker, bool) {
	if m, ok := d.marker[name]; ok {
		return m, ok
	}
	if d.aux != nil {
		if m, ok := d.aux.Marker[name]; ok {
			return m, ok
		}
	}
	return marker{}, false
}
