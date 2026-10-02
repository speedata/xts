package core

import (
	"strings"
	"testing"
)

// tableRows is a one-column table of n rows with an id and a TableHead.
func tableRows(id string, n int) string {
	var b strings.Builder
	b.WriteString(`<Table id="` + id + `" width="5"><TableHead><Tr><Td><Paragraph><Value>Head</Value></Paragraph></Td></Tr></TableHead>`)
	for i := 0; i < n; i++ {
		b.WriteString(`<Tr><Td><Paragraph><Value>row</Value></Paragraph></Td></Tr>`)
	}
	b.WriteString(`</Table>`)
	return b.String()
}

// pagesWith returns the pages, counted from 1, that hold a box with the id.
func pagesWith(root dumpNode, id string) []int {
	var pages []int
	for i, pg := range root.Children {
		if _, _, ok := pg.find(id, ""); ok {
			pages = append(pages, i+1)
		}
	}
	return pages
}

// PlaceObject keeps a table whole: one that does not fit below what is on
// the page moves on to the next page as one object, head included.
func TestPlaceObjectKeepsATableWhole(t *testing.T) {
	root := renderDump(t, layoutHead+`
  <PageFormat width="105mm" height="148mm"/>
  <SetGrid width="5mm" height="12pt"/>
  <Record match="data">
    <PlaceObject><Box width="17" height="20" backgroundcolor="gray"/></PlaceObject>
    <PlaceObject>`+tableRows("t", 15)+`</PlaceObject>
  </Record>
</Layout>`)
	if got := pagesWith(root, "t"); len(got) != 1 || got[0] != 2 {
		t.Fatalf("the table is on pages %v, want page 2 only", got)
	}
	if _, y, _, _, _ := dumpBox(t, root.Children[1], "t"); !near(y, 148/25.4*72-28.35) {
		t.Errorf("the table starts at y %.2f, want the top of page 2's grid", y)
	}
}

// A table taller than any frame of the area cannot be placed whole: that is
// an error naming PlaceObject and its line, pointing to Flow. Placed at an
// absolute position, it is placed there as before.
func TestPlaceObjectTableTallerThanTheFrame(t *testing.T) {
	_, err := runDump(t, layoutHead+`
  <PageFormat width="105mm" height="148mm"/>
  <Record match="data">
    <PlaceObject>`+tableRows("t", 60)+`</PlaceObject>
  </Record>
</Layout>`)
	if want := "PlaceObject (line 4): the table is"; err == nil || !strings.Contains(err.Error(), want) || !strings.Contains(err.Error(), "Flow") {
		t.Errorf("error %v, want one with %q that mentions Flow", err, want)
	}

	root := renderDump(t, layoutHead+`
  <PageFormat width="105mm" height="148mm"/>
  <Record match="data">
    <PlaceObject column="1cm" row="1cm">`+tableRows("t", 60)+`</PlaceObject>
  </Record>
</Layout>`)
	if got := pagesWith(root, "t"); len(got) != 1 || got[0] != 1 {
		t.Errorf("the absolutely placed table is on pages %v, want page 1 only", got)
	}
}

// A table that fits a frame of the area waits for one: it does not go into a
// lower frame and run past its bottom. Here frame 2 is 10 rows high, the table
// about 20, and the rest of frame 1 too short, so it starts page 2.
func TestPlaceObjectTableSkipsALowerFrame(t *testing.T) {
	root := renderDump(t, layoutHead+`
  <PageFormat width="105mm" height="148mm"/>
  <SetGrid nx="21" ny="30"/>
  <DefineMasterPage name="page" test="true()" margin="10mm">
    <PositioningArea name="cols">
      <PositioningFrame column="1" row="1" width="10" height="30"/>
      <PositioningFrame column="12" row="1" width="10" height="10"/>
    </PositioningArea>
  </DefineMasterPage>
  <Record match="data">
    <PlaceObject area="cols"><Box width="10" height="16"/></PlaceObject>
    <PlaceObject area="cols">`+tableRows("t", 18)+`</PlaceObject>
  </Record>
</Layout>`)
	if got := pagesWith(root, "t"); len(got) != 1 || got[0] != 2 {
		t.Errorf("the table is on pages %v, want page 2 only", got)
	}
}
