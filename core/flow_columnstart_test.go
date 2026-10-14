package core

import (
	"strconv"
	"strings"
	"testing"
)

// columnStartLayout pours paragraphs of three lines into an area of frames
// side by side, after objects, with column-fill fill, and logs the flow's
// bottom. Each entry of classes is the class of one paragraph; "space" gives
// it a margin-top of one line.
func columnStartLayout(frames, objects, fill string, classes ...string) string {
	var paras strings.Builder
	for i, class := range classes {
		paras.WriteString(`<Paragraph class="` + class + `"><Value>Paragraph ` + strconv.Itoa(i+1) + `: lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor.</Value></Paragraph>`)
	}
	return layoutHead + `
  <PageFormat width="105mm" height="148mm"/>
  <SetGrid nx="21" height="12pt"/>
  <DefineMasterPage name="page" test="true()" margin="10mm">
    <PositioningArea name="cols">` + frames + `</PositioningArea>
  </DefineMasterPage>
  <StyleSheet>
    p { margin: 0; font-size: 9pt; line-height: 11pt }
    p.space { margin-top: 11pt }
  </StyleSheet>
  <Record match="data">` + objects + `
    <Flow area="cols" column-fill="` + fill + `" bottom="b">` + paras.String() + `</Flow>
    <Message select="concat('bottom=', $b, ' page=', sd:current-page())"/>
  </Record>
</Layout>`
}

// twoFrames is two frames side by side, h1 and h2 rows high.
func twoFrames(h1, h2 int) string {
	return `<PositioningFrame column="1" row="1" width="10" height="` + strconv.Itoa(h1) + `"/>` +
		`<PositioningFrame column="12" row="1" width="10" height="` + strconv.Itoa(h2) + `"/>`
}

// A block that starts the second frame below a Flow across the page loses
// its margin-top at the automatic break, as the second column below a
// spanner does, so the flow ends where it ends without the margin (xts#74).
// The columns are balanced, filled, or so short that the margin would push
// the block to the next page.
func TestFlowColumnStartTruncatesMargin(t *testing.T) {
	for _, c := range []struct {
		name, frames, fill string
	}{
		{"balanced", twoCols, "balance"},
		{"filled", twoFrames(5, 5), "auto"},
		{"no room for the margin", twoFrames(4, 4), "auto"},
	} {
		t.Run(c.name, func(t *testing.T) {
			bottom, page := flowBottom(t, columnStartLayout(c.frames, introFlow, c.fill, "", "space"))
			wantBottom, wantPage := flowBottom(t, columnStartLayout(c.frames, introFlow, c.fill, "", ""))
			if bottom != wantBottom || page != wantPage {
				t.Errorf("the flow ends at %v on page %d, want %v on page %d as without the margin", bottom, page, wantBottom, wantPage)
			}
			if page != 1 {
				t.Errorf("the flow ends on page %d, want 1", page)
			}
		})
	}
}

// A block that moves on to the second band of a frame, below an object in
// the column, keeps its margin-top: only the frame's first band is the top
// of a column.
func TestFlowLaterBandKeepsMargin(t *testing.T) {
	// The first frame holds paragraph 1, the first band of the second one
	// paragraph 2, and paragraph 3 moves on below the object, where the band
	// is tall enough for it with its margin.
	frames := twoFrames(4, 9)
	objects := introFlow + `<PlaceObject column="12" row="5"><Box width="10" height="1"/></PlaceObject>`
	bottom, page := flowBottom(t, columnStartLayout(frames, objects, "auto", "", "space", "space"))
	without, _ := flowBottom(t, columnStartLayout(frames, objects, "auto", "", "space", ""))
	if page != 1 {
		t.Fatalf("the flow ends on page %d, want 1", page)
	}
	if bottom-without != 11 {
		t.Errorf("the flow ends at %v, %vpt below the flow without the margin, want 11pt", bottom, bottom-without)
	}
}
