package core

import (
	"bytes"
	"encoding/xml"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

type dumpNode struct {
	Attrs    []xml.Attr `xml:",any,attr"`
	Children []dumpNode `xml:",any"`
}

func (n dumpNode) attr(name string) string {
	for _, a := range n.Attrs {
		if a.Name.Local == name {
			return a.Value
		}
	}
	return ""
}

// find returns the box with the given PlaceObject id and the id of the
// nearest box around it that has one.
func (n dumpNode) find(id, parent string) (dumpNode, string, bool) {
	for _, c := range n.Children {
		if c.attr("attr-id") == id {
			return c, parent, true
		}
		p := parent
		if cid := c.attr("attr-id"); cid != "" {
			p = cid
		}
		if f, fp, ok := c.find(id, p); ok {
			return f, fp, true
		}
	}
	return dumpNode{}, "", false
}

// renderDump typesets layout and returns its geometry dump.
func renderDump(t *testing.T, layout string) dumpNode {
	t.Helper()
	jobname := filepath.Join(t.TempDir(), "out")
	out, err := os.Create(jobname + ".pdf")
	if err != nil {
		t.Fatal(err)
	}
	var dump bytes.Buffer
	err = RunXTS(&XTSConfig{
		Datafile:    strings.NewReader("<data/>"),
		Layoutfile:  strings.NewReader(layout),
		DumpFile:    &dump,
		FindFile:    FindFile,
		Jobname:     jobname,
		Outfile:     out,
		OutFilename: out.Name(),
	})
	if err != nil {
		t.Fatal(err)
	}
	var root dumpNode
	if err := xml.Unmarshal(dump.Bytes(), &root); err != nil {
		t.Fatal(err)
	}
	return root
}

// dumpBox finds the box with the given id in the dump and returns its position
// and size in points, and the id of the box around it.
func dumpBox(t *testing.T, root dumpNode, id string) (x, y, wd, ht float64, parent string) {
	t.Helper()
	n, parent, ok := root.find(id, "")
	if !ok {
		t.Fatalf("no box %q in the dump", id)
	}
	num := func(a string) float64 {
		f, err := strconv.ParseFloat(n.attr(a), 64)
		if err != nil {
			t.Fatalf("box %q: %s: %v", id, a, err)
		}
		return f
	}
	return num("x"), num("y"), num("width"), num("height") + num("depth"), parent
}

func near(a, b float64) bool { return math.Abs(a-b) < 0.02 }

// An object placed at lengths in a slate goes into the slate, measured from
// its top left corner, like one placed in the slate's grid.
func TestSlateLengthsGoIntoTheSlate(t *testing.T) {
	root := renderDump(t, `<Layout xmlns="urn:speedata.de/2021/xts/en">
  <Record match="data">
    <Slate name="s">
      <Contents>
        <PlaceObject id="text"><TextBlock width="5"><Paragraph><Value>Band</Value></Paragraph></TextBlock></PlaceObject>
        <PlaceObject id="mark" column="10mm" row="5mm"><Box width="5mm" height="30mm"/></PlaceObject>
      </Contents>
    </Slate>
    <PlaceObject id="band" column="3" row="4" slate="s"/>
  </Record>
</Layout>`)
	bx, by, _, bht, _ := dumpBox(t, root, "band")
	mx, my, _, _, parent := dumpBox(t, root, "mark")
	if parent != "band" {
		t.Fatalf("mark is in %q, want it in the slate", parent)
	}
	const mm = 72 / 25.4
	if !near(mx, bx+10*mm) || !near(my, by-5*mm) {
		t.Errorf("mark at %.2f, %.2f, want %.2f, %.2f", mx, my, bx+10*mm, by-5*mm)
	}
	// The mark reaches below the text, so it sizes the slate.
	if !near(bht, 35*mm) {
		t.Errorf("slate height %.2f, want %.2f", bht, 35*mm)
	}
}

// ids lists the PlaceObject ids under n in the order they are drawn.
func (n dumpNode) ids() []string {
	var out []string
	for _, c := range n.Children {
		if id := c.attr("attr-id"); id != "" {
			out = append(out, id)
		}
		out = append(out, c.ids()...)
	}
	return out
}

// With allocate="no" an object in a slate takes no room: at lengths or in
// the grid, negative offsets included, it neither sizes the slate nor moves
// what follows.
func TestSlateAllocateNoTakesNoRoom(t *testing.T) {
	root := renderDump(t, `<Layout xmlns="urn:speedata.de/2021/xts/en">
  <Record match="data">
    <Slate name="s">
      <Contents>
        <PlaceObject id="logo" column="-20mm" row="-5mm" allocate="no"><Box width="100mm" height="40mm"/></PlaceObject>
        <PlaceObject id="text"><TextBlock width="5"><Paragraph><Value>Band</Value></Paragraph></TextBlock></PlaceObject>
        <PlaceObject id="tint" column="1" row="2" allocate="no"><Box width="12" height="6"/></PlaceObject>
        <PlaceObject id="next" column="1"><TextBlock width="5"><Paragraph><Value>Next</Value></Paragraph></TextBlock></PlaceObject>
      </Contents>
    </Slate>
    <PlaceObject id="band" column="3" row="4" slate="s"/>
  </Record>
</Layout>`)
	bx, by, bwd, bht, _ := dumpBox(t, root, "band")
	_, ty, twd, _, _ := dumpBox(t, root, "text")
	_, ny, _, nht, _ := dumpBox(t, root, "next")
	lx, ly, _, _, parent := dumpBox(t, root, "logo")
	if parent != "band" {
		t.Fatalf("logo is in %q, want it in the slate", parent)
	}
	const mm = 72 / 25.4
	if !near(lx, bx-20*mm) || !near(ly, by+5*mm) {
		t.Errorf("logo at %.2f, %.2f, want %.2f, %.2f", lx, ly, bx-20*mm, by+5*mm)
	}
	// The tint covers the rows below the text, and next still starts in
	// the first of them.
	_, iy, _, _, _ := dumpBox(t, root, "tint")
	if !near(ny, iy) || ny >= ty {
		t.Errorf("next at %.2f, want %.2f, the tint's row", ny, iy)
	}
	if !near(bwd, twd) || !near(bht, by-ny+nht) {
		t.Errorf("slate %.2f by %.2f, want %.2f by %.2f", bwd, bht, twd, by-ny+nht)
	}
	got := strings.Join(root.ids(), " ")
	if want := "band logo text tint next"; got != want {
		t.Errorf("drawn in the order %q, want %q", got, want)
	}
}

// An object with allocate="no" in a slate is drawn in the order it was
// placed, as on a page: a background placed first lies under the text.
func TestSlateAllocateNoDrawsInOrder(t *testing.T) {
	root := renderDump(t, `<Layout xmlns="urn:speedata.de/2021/xts/en">
  <Record match="data">
    <Slate name="band">
      <Contents>
        <PlaceObject id="bg" column="1" row="1" allocate="no"><Box width="60mm" height="10mm" backgroundcolor="orange"/></PlaceObject>
        <PlaceObject id="text"><TextBlock width="10"><Paragraph><Value>Text over a background.</Value></Paragraph></TextBlock></PlaceObject>
      </Contents>
    </Slate>
    <PlaceObject id="band" column="1" row="1" slate="band"/>
  </Record>
</Layout>`)
	got := strings.Join(root.ids(), " ")
	if want := "band bg text"; got != want {
		t.Errorf("drawn in the order %q, want %q", got, want)
	}
}
