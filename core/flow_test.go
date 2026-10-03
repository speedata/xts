package core

import (
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// The children of a Flow are dispatched before the flow is laid out, so a
// command that acts on the page then is an error there, naming the command
// and its line, also inside a ForAll. CallTemplate still works.
func TestFlowRejectsPageActions(t *testing.T) {
	for _, c := range []struct{ name, child string }{
		{"PlaceObject", `<PlaceObject><Box width="1cm" height="1cm"/></PlaceObject>`},
		{"ClearPage", `<ClearPage/>`},
		{"NextFrame", `<NextFrame/>`},
		{"NextRow", `<NextRow/>`},
		{"PlaceObject", `<ForAll select="."><PlaceObject><Box width="1cm" height="1cm"/></PlaceObject></ForAll>`},
	} {
		t.Run(c.name, func(t *testing.T) {
			_, err := runDump(t, layoutHead+`
  <Record match="data">
    <Flow>
      <Paragraph><Value>First</Value></Paragraph>
      `+c.child+`
      <Paragraph><Value>Second</Value></Paragraph>
    </Flow>
  </Record>
</Layout>`)
			if want := c.name + " (line 5): not allowed inside a Flow"; err == nil || !strings.Contains(err.Error(), want) {
				t.Errorf("error %v, want one with %q", err, want)
			}
		})
	}
	t.Run("CallTemplate", func(t *testing.T) {
		root := renderDump(t, layoutHead+`
  <Template name="set"><SetVariable variable="v" select="'From a template'"/></Template>
  <Record match="data">
    <Flow><CallTemplate name="set"/><Paragraph id="p"><Value select="$v"/></Paragraph></Flow>
  </Record>
</Layout>`)
		if _, _, ok := root.find("p", ""); !ok {
			t.Error("the paragraph after the CallTemplate is not in the flow")
		}
	})
}

// bottom of a Flow that places nothing is where the flow would have started:
// at the exact bottom of a box placed before it, or where a Flow it follows
// directly ended.
func TestFlowBottomWithoutBlocks(t *testing.T) {
	// bottom returns the values of bottom=… the layout logged, each once.
	bottom := func(t *testing.T, body string) []string {
		t.Helper()
		log := runLayoutLog(t, layoutHead+`
  <SetGrid width="5mm" height="12pt"/>
  <Record match="data">`+body+`</Record>
</Layout>`)
		var got []string
		for _, m := range regexp.MustCompile(`bottom=([-0-9.]+)`).FindAllStringSubmatch(log, -1) {
			if !slices.Contains(got, m[1]) {
				got = append(got, m[1])
			}
		}
		return got
	}
	t.Run("below a box", func(t *testing.T) {
		got := bottom(t, `
    <PlaceObject><Box width="5cm" height="30pt"/></PlaceObject>
    <Flow bottom="b"><ForAll select="nothing"><Paragraph><Value>x</Value></Paragraph></ForAll></Flow>
    <Message select="concat('bottom=', $b)"/>`)
		if len(got) != 1 || got[0] != "30" {
			t.Errorf("bottom %v, want [30]", got)
		}
	})
	t.Run("after a flow", func(t *testing.T) {
		got := bottom(t, `
    <Flow bottom="a"><Paragraph style="margin: 0; line-height: 15pt"><Value>x</Value></Paragraph></Flow>
    <Message select="concat('bottom=', $a)"/>
    <Flow bottom="b"><ForAll select="nothing"><Paragraph><Value>x</Value></Paragraph></ForAll></Flow>
    <Message select="concat('bottom=', $b)"/>`)
		if len(got) != 1 {
			t.Errorf("bottom %v, want one value for both flows", got)
		}
	})
}

// A child that cannot be part of a flow is left out with a warning naming
// the command and its own line.
func TestFlowNamesLeftOutChild(t *testing.T) {
	img, err := filepath.Abs("../qa/helloworld/firstpage.png")
	if err != nil {
		t.Fatal(err)
	}
	log := runLayoutLog(t, layoutHead+`
  <Record match="data">
    <Flow>
      <Paragraph><Value>Kept</Value></Paragraph>
      <TextBlock><Paragraph><Value>A text block</Value></Paragraph></TextBlock>
      <Box width="1cm" height="1cm"/>
      <Image href="`+img+`" width="1cm"/>
    </Flow>
  </Record>
</Layout>`)
	for _, want := range []string{"TextBlock (line 5)", "Box (line 6)", "Image (line 7)"} {
		if !strings.Contains(log, want+": cannot be part of a flow and is left out") {
			t.Errorf("no warning for %s in\n%s", want, log)
		}
	}
}

// A table with stretch="max" and a width is that wide in a flow, as it is in
// a PlaceObject, not as wide as the frame.
func TestFlowTableStretchKeepsItsWidth(t *testing.T) {
	root := renderDump(t, layoutHead+`
  <SetGrid width="5mm" height="12pt"/>
  <Record match="data">
    <Flow>
      <Table id="t" stretch="max" width="10"><Tr><Td><Paragraph><Value>A</Value></Paragraph></Td></Tr></Table>
    </Flow>
  </Record>
</Layout>`)
	n, _, ok := root.find("t", "")
	if !ok {
		t.Fatal("no table t in the dump")
	}
	if wd := n.attr("width"); wd != "141.73" {
		t.Errorf("the table is %spt wide, want 141.73 (10 columns of 5mm)", wd)
	}
}
