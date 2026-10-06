package core

import (
	"regexp"
	"strings"
	"testing"
)

// flowSlateLayout composes the slate s from body, with three-line flows of
// 13pt lines in a grid of 12pt rows, and places it with the id s.
func flowSlateLayout(body string) string {
	return layoutHead + `
  <SetGrid width="5mm" height="12pt"/>
  <StyleSheet>p { margin: 0; font-size: 10pt; line-height: 13pt }</StyleSheet>
  <Record match="data">
    <Slate name="s">
      <Contents>` + body + `</Contents>
    </Slate>
    <PlaceObject id="s" slate="s"/>
  </Record>
</Layout>`
}

const threeParagraphs = `<Paragraph><Value>Quarterly report</Value></Paragraph>
        <Paragraph><Value>Department of synthetic examples</Value></Paragraph>
        <Paragraph><Value>Draft, not for circulation</Value></Paragraph>`

// A Flow in a slate stacks its blocks as one object at their exact heights:
// three lines of 13pt make a slate of 39pt, not three objects of whole rows.
func TestFlowInASlateIsExact(t *testing.T) {
	_, _, _, ht, _ := dumpBox(t, renderDump(t, flowSlateLayout(`<Flow>`+threeParagraphs+`</Flow>`)), "s")
	if !near(ht, 39) {
		t.Errorf("slate height %.2f, want 39", ht)
	}
}

// The margins between the blocks of a Flow in a slate collapse as on a page:
// 4pt above the first block, then the larger 8pt between the blocks.
func TestFlowInASlateCollapsesMargins(t *testing.T) {
	_, _, _, ht, _ := dumpBox(t, renderDump(t, flowSlateLayout(`<Flow>
        <Paragraph style="margin: 4pt 0 8pt"><Value>One</Value></Paragraph>
        <Paragraph style="margin: 4pt 0 8pt"><Value>Two</Value></Paragraph>
        <Paragraph style="margin: 4pt 0 8pt"><Value>Three</Value></Paragraph>
      </Flow>`)), "s")
	if want := 4.0 + 13 + 8 + 13 + 8 + 13; !near(ht, want) {
		t.Errorf("slate height %.2f, want %.2f", ht, want)
	}
}

// sd:slate-height and bottom after a Flow in a slate give the flow's exact
// end.
func TestFlowInASlateHeightAndBottom(t *testing.T) {
	log := runLayoutLog(t, flowSlateLayout(`<Flow bottom="b">`+threeParagraphs+`</Flow>
        <Message select="concat('height=', sd:slate-height('s', 'pt'), ' bottom=', $b)"/>`))
	if !strings.Contains(log, "height=39 bottom=39") {
		t.Errorf("want height=39 bottom=39 in the log:\n%s", log)
	}
}

// A PlaceObject in the slate after the Flow starts on the next whole row
// below it: the flow ends at 39pt, in the fourth row of 12pt, so the box
// starts at 48pt.
func TestFlowInASlateThenPlaceObject(t *testing.T) {
	root := renderDump(t, flowSlateLayout(`<Flow>`+threeParagraphs+`</Flow>
        <PlaceObject id="after"><Box width="5mm" height="12pt"/></PlaceObject>`))
	_, sy, _, ht, _ := dumpBox(t, root, "s")
	_, ay, _, _, parent := dumpBox(t, root, "after")
	if parent != "s" {
		t.Fatalf("the box is in %q, want it in the slate", parent)
	}
	if !near(sy-ay, 48) {
		t.Errorf("the box starts %.2fpt below the slate's top, want 48", sy-ay)
	}
	if !near(ht, 60) {
		t.Errorf("slate height %.2f, want 60", ht)
	}
}

// A forced break in a Flow in a slate has nothing to break to: the blocks
// stack as without it, and the flow warns once, however many breaks it has.
func TestFlowInASlateWarnsOnAForcedBreak(t *testing.T) {
	layout := flowSlateLayout(`<Flow>
        <Paragraph><Value>One</Value></Paragraph>
        <Paragraph style="break-before: page"><Value>Two</Value></Paragraph>
        <Paragraph style="break-after: column"><Value>Three</Value></Paragraph>
        <Paragraph style="break-before: right"><Value>Four</Value></Paragraph>
      </Flow>`)
	if _, _, _, ht, _ := dumpBox(t, renderDump(t, layout), "s"); !near(ht, 52) {
		t.Errorf("slate height %.2f, want 52", ht)
	}
	log := runLayoutLog(t, layout)
	if n := len(regexp.MustCompile(`a forced break \(\w+\) in a slate has nothing to break to`).FindAllString(log, -1)); n != 1 {
		t.Errorf("%d warnings about the forced break, want 1:\n%s", n, log)
	}
}

// A Flow in a slate takes no area, as a PlaceObject there does not: the
// slate has none of its own.
func TestFlowInASlateRejectsAnArea(t *testing.T) {
	_, err := runDump(t, layoutHead+`
  <DefineMasterPage name="page" test="true()" margin="1cm">
    <PositioningArea name="text"><PositioningFrame width="10" height="10" row="1" column="1"/></PositioningArea>
  </DefineMasterPage>
  <Record match="data">
    <Slate name="s"><Contents><Flow area="text"><Paragraph><Value>x</Value></Paragraph></Flow></Contents></Slate>
  </Record>
</Layout>`)
	if want := "area text not found: a flow in slate s"; err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("error %v, want one with %q", err, want)
	}
}

// A Mark in a Flow in a slate, before the first block, between two and after
// the last, takes the page the slate is placed on, not the one it was
// composed on.
func TestFlowInASlateMarkTakesThePageItIsPlacedOn(t *testing.T) {
	marks, _ := runMarks(t, layoutHead+`
  <Record match="data">
    <PlaceObject><Box width="2" height="1"/></PlaceObject>
    <Slate name="s">
      <Contents>
        <Flow>
          <Mark select="'first'"/>
          <Paragraph><Value>One</Value></Paragraph>
          <Mark select="'between'"/>
          <Paragraph><Value>Two</Value></Paragraph>
          <Mark select="'last'"/>
        </Flow>
      </Contents>
    </Slate>
    <ClearPage/>
    <PlaceObject><Box width="2" height="1"/></PlaceObject>
    <ClearPage/>
    <PlaceObject slate="s"/>
  </Record>
</Layout>`)
	for _, name := range []string{"first", "between", "last"} {
		if got, ok := marks[name]; !ok || got != 3 {
			t.Errorf("mark %s on page %d (%t), want 3", name, got, ok)
		}
	}
}
