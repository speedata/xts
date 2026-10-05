package core

import (
	"strings"
	"testing"
)

// A Flow over ProcessNode gets the blocks of the records, and a Mark in a
// record takes the page of the block that follows it, as in a ForAll.
func TestFlowOverProcessNode(t *testing.T) {
	marks, boxes := runMarks(t, layoutHead+`
  <PageFormat width="100mm" height="60mm"/>
  <SetGrid width="5mm" height="12pt"/>
  <StyleSheet>p { margin: 0; font-size: 10pt; line-height: 12pt }</StyleSheet>
  <Record match="data">
    <Flow>
      <Paragraph id="a">`+lines(2)+`</Paragraph>
      <ProcessNode select="." mode="chapter"/>
    </Flow>
  </Record>
  <Record match="data" mode="chapter">
    <Mark select="'m'"/>
    <Paragraph id="x" style="break-before: page">`+lines(1)+`</Paragraph>
  </Record>
</Layout>`)
	if boxes["a"] != 1 || boxes["x"] != 2 {
		t.Fatalf("the blocks are on pages %d and %d, want 1 and 2", boxes["a"], boxes["x"])
	}
	if marks["m"] != 2 {
		t.Errorf("the mark is on page %d, want 2", marks["m"])
	}
}

// An error in a record called by ProcessNode ends the run, as it does in a
// ForAll.
func TestProcessNodePassesOnTheError(t *testing.T) {
	_, err := runDump(t, layoutHead+`
  <Record match="data">
    <ProcessNode select="." mode="m"/>
  </Record>
  <Record match="data" mode="m">
    <PlaceObject area="nosuch"><TextBlock><Paragraph><Value>A</Value></Paragraph></TextBlock></PlaceObject>
  </Record>
</Layout>`)
	if want := "area nosuch not found"; err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("error %v, want one with %q", err, want)
	}
}
