package core

import (
	"encoding/xml"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// runMarks typesets layout and returns the pages of its marks and of the
// boxes with an id, from the aux file and the dump.
func runMarks(t *testing.T, layout string) (marks, boxes map[string]int) {
	t.Helper()
	dir := t.TempDir()
	jobname := filepath.Join(dir, "out")
	out, err := os.Create(jobname + ".pdf")
	if err != nil {
		t.Fatal(err)
	}
	var dump strings.Builder
	if err := RunXTS(&XTSConfig{
		Datafile:    strings.NewReader("<data/>"),
		Layoutfile:  strings.NewReader(layout),
		DumpFile:    &dump,
		FindFile:    FindFile,
		Jobname:     jobname,
		Outfile:     out,
		OutFilename: out.Name(),
	}); err != nil {
		t.Fatal(err)
	}
	aux, err := os.ReadFile(jobname + "-aux.xml")
	if err != nil {
		t.Fatal(err)
	}
	var a struct {
		Marks []struct {
			Name string `xml:"name,attr"`
			Page string `xml:"page,attr"`
		} `xml:"marker>mark"`
	}
	if err := xml.Unmarshal(aux, &a); err != nil {
		t.Fatal(err)
	}
	marks = map[string]int{}
	for _, m := range a.Marks {
		marks[m.Name], _ = strconv.Atoi(m.Page)
	}
	var root dumpNode
	if err := xml.Unmarshal([]byte(dump.String()), &root); err != nil {
		t.Fatal(err)
	}
	boxes = map[string]int{}
	for i, pg := range root.Children {
		var walk func(n dumpNode)
		walk = func(n dumpNode) {
			if id := n.attr("attr-id"); id != "" {
				if _, ok := boxes[id]; !ok {
					boxes[id] = i + 1
				}
			}
			for _, c := range n.Children {
				walk(c)
			}
		}
		walk(pg)
	}
	return marks, boxes
}

// lines is the children of a Paragraph of n lines.
func lines(n int) string {
	s := []string{"<Value>Line</Value>"}
	for i := 1; i < n; i++ {
		s = append(s, "<Br/><Value>Line</Value>")
	}
	return strings.Join(s, "")
}

// flowMarksLayout is a page with room for nine lines, the flow's first
// paragraph taking seven of them.
func flowMarksLayout(flow string) string {
	return layoutHead + `
  <PageFormat width="100mm" height="60mm"/>
  <SetGrid width="5mm" height="12pt"/>
  <StyleSheet>p { margin: 0; font-size: 10pt; line-height: 12pt }</StyleSheet>
  <Record match="data">
    <Flow>
      <Paragraph id="a">` + lines(7) + `</Paragraph>
      ` + flow + `
    </Flow>
  </Record>
</Layout>`
}

// A Mark in a flow marks the page of the block that follows it, wherever that
// block lands: here on page 2, where it moves on for each of the reasons a
// block moves on.
func TestFlowMarkTakesThePageOfTheNextBlock(t *testing.T) {
	cases := []struct{ name, flow string }{
		{"fit", `<Mark select="'m'"/><Paragraph id="x" style="line-height: 40pt">` + lines(1) + `</Paragraph>`},
		{"orphans", `<Mark select="'m'"/><Paragraph id="x" style="orphans: 4">` + lines(6) + `</Paragraph>`},
		{"break-after: avoid", `<Mark select="'m'"/><Paragraph id="x" style="break-after: avoid">` + lines(1) + `</Paragraph><Paragraph style="orphans: 3">` + lines(6) + `</Paragraph>`},
		{"forced break", `<Mark select="'m'"/><Paragraph id="x" style="break-before: page">` + lines(1) + `</Paragraph>`},
		{"after loose text", `<HTML>Loose <b>text</b></HTML><Mark select="'m'"/><HTML><p id="x" style="break-before: page">Line</p></HTML>`},
		{"in an Action", `<Action><Mark select="'m'"/></Action><Paragraph id="x" style="break-before: page">` + lines(1) + `</Paragraph>`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			marks, boxes := runMarks(t, flowMarksLayout(c.flow))
			if boxes["x"] != 2 {
				t.Fatalf("the block is on page %d, want it moved on to page 2", boxes["x"])
			}
			if marks["m"] != 2 {
				t.Errorf("the mark is on page %d, want 2", marks["m"])
			}
		})
	}
}

// A mark first in the flow takes the page of the first block, also with
// break-before: page on it. A mark after the last block takes the page the
// flow ends on.
func TestFlowMarkFirstAndLast(t *testing.T) {
	marks, boxes := runMarks(t, layoutHead+`
  <PageFormat width="100mm" height="60mm"/>
  <SetGrid width="5mm" height="12pt"/>
  <StyleSheet>p { margin: 0; font-size: 10pt; line-height: 12pt }</StyleSheet>
  <Record match="data">
    <PlaceObject><Box width="5" height="1"/></PlaceObject>
    <Flow>
      <Mark select="'first'"/>
      <Paragraph id="x" style="break-before: page">`+lines(12)+`</Paragraph>
      <Mark select="'last'"/>
    </Flow>
  </Record>
</Layout>`)
	if marks["first"] != boxes["x"] {
		t.Errorf("the first mark is on page %d, the block starts on page %d", marks["first"], boxes["x"])
	}
	if marks["last"] != boxes["x"]+1 {
		t.Errorf("the last mark is on page %d, want %d where the split block ends", marks["last"], boxes["x"]+1)
	}
}

// A Mark inside a block of a flow has no page of its own yet: an error.
func TestFlowMarkInsideAParagraph(t *testing.T) {
	_, err := runDump(t, layoutHead+`
  <Record match="data">
    <Flow>
      <Paragraph><Action><Mark select="'m'"/></Action><Value>A</Value></Paragraph>
    </Flow>
  </Record>
</Layout>`)
	if want := "Mark (line 4): not allowed inside a block of a Flow"; err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("error %v, want one with %q", err, want)
	}
}

// pdfDests reads the page numbers and y positions of the named destinations
// and outline entries in an uncompressed PDF, by the order of the pages in
// /Kids.
func pdfDests(t *testing.T, pdf string) (named, outlines map[string][2]float64) {
	t.Helper()
	kids := regexp.MustCompile(`/Kids \[([^\]]*)\]`).FindStringSubmatch(pdf)
	if kids == nil {
		t.Fatal("no /Kids in the PDF")
	}
	pageNum := map[string]int{}
	for i, ref := range regexp.MustCompile(`(\d+) 0 R`).FindAllStringSubmatch(kids[1], -1) {
		pageNum[ref[1]] = i + 1
	}
	at := func(page, y string) [2]float64 {
		f, _ := strconv.ParseFloat(y, 64)
		return [2]float64{float64(pageNum[page]), f}
	}
	obj := func(num string) string {
		m := regexp.MustCompile(`(?s)\n` + num + ` 0 obj(.*?)endobj`).FindStringSubmatch(pdf)
		if m == nil {
			t.Fatalf("no object %s", num)
		}
		return m[1]
	}
	named = map[string][2]float64{}
	if names := regexp.MustCompile(`(?s)/Dests <<.*?/Names \[([^\]]*)\]`).FindStringSubmatch(pdf); names != nil {
		for _, m := range regexp.MustCompile(`\(([^)]*)\) (\d+) 0 R`).FindAllStringSubmatch(names[1], -1) {
			if d := regexp.MustCompile(`/D \[(\d+) 0 R /XYZ \S+ (\S+)`).FindStringSubmatch(obj(m[2])); d != nil {
				named[m[1]] = at(d[1], d[2])
			}
		}
	}
	outlines = map[string][2]float64{}
	for _, m := range regexp.MustCompile(`/Dest \[ (\d+) 0 R /XYZ \S+ (\S+) 0\][^>]*?/Title \(([^)]*)\)`).FindAllStringSubmatch(pdf, -1) {
		outlines[m[3]] = at(m[1], m[2])
	}
	return named, outlines
}

// renderPDF typesets layout and returns the PDF as a string.
func renderPDF(t *testing.T, layout string) string {
	t.Helper()
	jobname := filepath.Join(t.TempDir(), "out")
	out, err := os.Create(jobname + ".pdf")
	if err != nil {
		t.Fatal(err)
	}
	if err := RunXTS(&XTSConfig{
		Datafile:    strings.NewReader("<data/>"),
		Layoutfile:  strings.NewReader(layout),
		FindFile:    FindFile,
		Jobname:     jobname,
		Outfile:     out,
		OutFilename: out.Name(),
	}); err != nil {
		t.Fatal(err)
	}
	out.Close()
	b, err := os.ReadFile(jobname + ".pdf")
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// A Mark with pdftarget is a named destination, and a Bookmark an outline
// entry, at the top of the block that follows it, here on page 2.
func TestFlowMarkPDFTargetAndBookmark(t *testing.T) {
	pdf := renderPDF(t, flowMarksLayout(`<Bookmark select="'Chapter'" level="1"/><Mark select="'m'" pdftarget="yes"/><Paragraph style="orphans: 4">`+lines(6)+`</Paragraph>`))
	named, outlines := pdfDests(t, pdf)
	// The grid's top on page 2: 60mm less the 10mm margin.
	top := 50 / 25.4 * 72
	for name, got := range map[string][2]float64{"named destination m": named["m"], "bookmark Chapter": outlines["Chapter"]} {
		if got[0] != 2 || math.Abs(got[1]-top) > 0.01 {
			t.Errorf("%s at page %v, y %.2f; want page 2, y %.2f", name, got[0], got[1], top)
		}
	}
}
