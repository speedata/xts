package core

import (
	"fmt"
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

// A forced break to a page side goes to the next page of that side: recto
// is a right (odd) page and verso a left (even) one, as right and left (#60).
func TestFlowBreakToPageSide(t *testing.T) {
	for _, tc := range []struct {
		value     string
		startPage int // the page the flow starts on
		want      string
	}{
		{"right", 1, "3"},
		{"recto", 1, "3"},
		{"left", 1, "2"},
		{"verso", 1, "2"},
		{"right", 2, "3"},
		{"recto", 2, "3"},
		{"left", 2, "4"},
		{"verso", 2, "4"},
		{"page", 1, "2"},
	} {
		t.Run(fmt.Sprintf("%s from page %d", tc.value, tc.startPage), func(t *testing.T) {
			clear := ""
			if tc.startPage == 2 {
				clear = `<PlaceObject><TextBlock><Paragraph><Value>1</Value></Paragraph></TextBlock></PlaceObject><ClearPage/>`
			}
			log := runLayoutLog(t, layoutHead+`
  <Record match="data">`+clear+`
    <Flow>
      <Paragraph><Value>Before</Value></Paragraph>
      <Paragraph style="break-before: `+tc.value+`"><Value>After</Value></Paragraph>
    </Flow>
    <Message select="concat('endpage=', sd:current-page())"/>
  </Record>
</Layout>`)
			m := regexp.MustCompile(`endpage=(\d+)`).FindStringSubmatch(log)
			if m == nil || m[1] != tc.want {
				t.Errorf("flow ends on page %v, want %s", m, tc.want)
			}
		})
	}
}

// A band that reaches the grid's last row, in a frame that ends there, runs
// on to the bottom margin. The page area is 113.39pt high, nine rows of 12pt
// and 5.39pt below them, and the lines 11.3pt: ten fit to the margin, nine
// to the last row.
func TestFlowRunsOnToTheBottomMargin(t *testing.T) {
	head := func(master string) string {
		return layoutHead + `
  <PageFormat width="100mm" height="60mm"/>
  <SetGrid width="5mm" height="12pt"/>
  <DefineMasterPage name="page" test="true()" margin="10mm">` + master + `</DefineMasterPage>
  <StyleSheet>p { margin: 0; font-size: 10pt; line-height: 11.3pt }</StyleSheet>
`
	}
	lines := func(area string, n int) string {
		var b strings.Builder
		fmt.Fprintf(&b, `<Flow%s>`, area)
		for i := 1; i <= n; i++ {
			fmt.Fprintf(&b, `<Paragraph id="l%d"><Value>Line %d</Value></Paragraph>`, i, i)
		}
		b.WriteString(`</Flow>`)
		return b.String()
	}
	// onFirstPage counts the lines on the first page that holds any. With a
	// footer from AtPageCreation, the flow starts on page 2: the footer's
	// PlaceObject leaves the current row below it on page 1.
	onFirstPage := func(t *testing.T, layout string) int {
		t.Helper()
		for _, pg := range renderDump(t, layout).Children {
			n := 0
			for _, id := range pg.ids() {
				if strings.HasPrefix(id, "l") {
					n++
				}
			}
			if n > 0 {
				return n
			}
		}
		return 0
	}
	for _, c := range []struct {
		name, master, area string
		want               int
	}{
		{"a frame that reaches the last row", "", "", 10},
		{"a frame that ends above it", `<PositioningArea name="a"><PositioningFrame column="1" row="1" width="16" height="8"/></PositioningArea>`, ` area="a"`, 8},
		{"two frames side by side", `<PositioningArea name="a"><PositioningFrame column="1" row="1" width="7" height="9"/><PositioningFrame column="9" row="1" width="8" height="9"/></PositioningArea>`, ` area="a"`, 20},
		{"a footer in the last row", `<AtPageCreation><PlaceObject column="1" row="9"><Box width="16" height="1"/></PlaceObject></AtPageCreation>`, "", 8},
	} {
		t.Run(c.name, func(t *testing.T) {
			layout := head(c.master) + `<Record match="data">` + lines(c.area, 30) + `</Record></Layout>`
			if got := onFirstPage(t, layout); got != c.want {
				t.Errorf("%d lines on the first page, want %d", got, c.want)
			}
			if log := runLayoutLog(t, layout); strings.Contains(log, "protrudes into the bottom margin") {
				t.Errorf("a warning about the bottom margin:\n%s", log)
			}
		})
	}
	// A block that fits in no empty frame and passes the margin still warns.
	t.Run("a block taller than the frame", func(t *testing.T) {
		layout := head("") + `<Record match="data"><Flow><Paragraph style="line-height: 150pt"><Value>Tall</Value></Paragraph></Flow></Record></Layout>`
		if log := runLayoutLog(t, layout); !strings.Contains(log, "protrudes into the bottom margin") {
			t.Errorf("no warning about the bottom margin:\n%s", log)
		}
	})
	// A Flow that follows one ending in the rest below the last row starts
	// on the next page, not in what is left of the rest.
	t.Run("a flow after one that ends in the rest", func(t *testing.T) {
		layout := head("") + `<Record match="data">` + lines("", 10) + `<Flow><Paragraph id="after"><Value>After</Value></Paragraph></Flow></Record></Layout>`
		root := renderDump(t, layout)
		if got := pagesWith(root, "l10"); !slices.Equal(got, []int{1}) {
			t.Errorf("the first flow's last line is on pages %v, want [1]", got)
		}
		if got := pagesWith(root, "after"); !slices.Equal(got, []int{2}) {
			t.Errorf("the second flow is on pages %v, want [2]", got)
		}
		if n := len(root.Children); n != 2 {
			t.Errorf("%d pages, want 2", n)
		}
	})
}
