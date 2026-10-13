package core

import (
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// flowBalanceLayout pours n paragraphs into an area of frames side by side,
// 27 rows of 12pt high, after objects, with column-fill fill, and logs the
// flow's bottom.
func flowBalanceLayout(frames, objects string, n int, fill string) string {
	return layoutHead + `
  <PageFormat width="105mm" height="148mm"/>
  <SetGrid nx="21" height="12pt"/>
  <DefineMasterPage name="page" test="true()" margin="10mm">
    <PositioningArea name="cols">` + frames + `</PositioningArea>
  </DefineMasterPage>
  <StyleSheet>p { margin: 0 0 4pt 0; font-size: 9pt; line-height: 11pt }</StyleSheet>
  <Record match="data">` + objects + `
    <Flow area="cols" column-fill="` + fill + `" bottom="b">
      <Loop select="` + strconv.Itoa(n) + `" variable="i">
        <Paragraph><Value select="concat('Paragraph ', $i, ': ')"/><Value>lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor.</Value></Paragraph>
      </Loop>
    </Flow>
    <Message select="concat('bottom=', $b, ' page=', sd:current-page())"/>
  </Record>
</Layout>`
}

var flowBottomRE = regexp.MustCompile(`bottom=([0-9.]+) page=([0-9]+)`)

// flowBottom runs layout and returns the flow's bottom and the page it ends
// on.
func flowBottom(t *testing.T, layout string) (float64, int) {
	t.Helper()
	m := flowBottomRE.FindStringSubmatch(runLayoutLog(t, layout))
	if m == nil {
		t.Fatal("no bottom in the log")
	}
	b, _ := strconv.ParseFloat(m[1], 64)
	p, _ := strconv.Atoi(m[2])
	return b, p
}

// The frames of the area, objects that take rows of them, and a Flow before
// it that the area's flow follows directly.
const (
	twoCols     = `<PositioningFrame column="1" row="1" width="10" height="27"/><PositioningFrame column="12" row="1" width="10" height="27"/>`
	threeCols   = `<PositioningFrame column="1" row="1" width="6" height="27"/><PositioningFrame column="8" row="1" width="6" height="27"/><PositioningFrame column="15" row="1" width="6" height="27"/>`
	introFlow   = `<Flow><Paragraph><Value>An introduction across the width of the page.</Value></Paragraph></Flow>`
	splitFrame2 = `<PlaceObject column="12" row="6"><Box width="10" height="3"/></PlaceObject>`
)

// With column-fill="balance", the end of the flow is spread over the frames
// of its last page, and its bottom is the end of the tallest one: after a
// Flow across the page, in a frame that an object parts in two bands, on
// the last of several pages, with two and three frames. Remaining has to
// name the regions Next hands out, or the balance gives up or comes out
// higher.
func TestFlowColumnFillBalance(t *testing.T) {
	for _, c := range []struct {
		name, frames, objects string
		n                     int
		// The bottom and page of the flow filled and balanced.
		auto, balance float64
		page          int
	}{
		{"after a flow", twoCols, introFlow, 4, 159, 85, 1},
		{"a frame in two bands", twoCols, splitFrame2, 4, 144, 107, 1},
		{"the last of three pages", twoCols, "", 30, 92, 44, 3},
		{"three frames", threeCols, introFlow, 6, 317, 151, 1},
		{"three frames, the last of four pages", threeCols, introFlow, 40, 287, 103, 4},
		{"three frames, one in two bands", threeCols, introFlow + splitFrame2, 6, 324, 188, 1},
	} {
		t.Run(c.name, func(t *testing.T) {
			for _, run := range []struct {
				fill string
				want float64
			}{{"auto", c.auto}, {"balance", c.balance}} {
				b, p := flowBottom(t, flowBalanceLayout(c.frames, c.objects, c.n, run.fill))
				if !near(b, run.want) || p != c.page {
					t.Errorf("column-fill=%s: bottom %.2f on page %d, want %.2f on page %d", run.fill, b, p, run.want, c.page)
				}
			}
		})
	}
}

// A flow in a slate has one region, so there is nothing to balance.
func TestFlowColumnFillInASlate(t *testing.T) {
	_, _, _, ht, _ := dumpBox(t, renderDump(t, flowSlateLayout(`<Flow column-fill="balance">`+threeParagraphs+`</Flow>`)), "s")
	if !near(ht, 39) {
		t.Errorf("slate height %.2f, want 39", ht)
	}
}

func TestFlowColumnFillValue(t *testing.T) {
	_, err := runDump(t, flowBalanceLayout(twoCols, "", 1, "both"))
	if want := `Flow (line 9): column-fill must be auto or balance, not "both"`; err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("error %v, want one with %q", err, want)
	}
}
