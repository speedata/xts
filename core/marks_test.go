package core

import (
	"math"
	"testing"
)

// A Mark in a slate takes the page the slate is placed on, not the page it
// was composed on; placed twice, it keeps the last one, as a second Mark of
// the same name does.
func TestMarkInASlateTakesThePageItIsPlacedOn(t *testing.T) {
	slate := `<Slate name="s">
      <Contents>
        <PlaceObject>
          <TextBlock>
            <Action><Mark select="'m'"/></Action>
            <Paragraph><Value>A</Value></Paragraph>
          </TextBlock>
        </PlaceObject>
      </Contents>
    </Slate>`
	for _, c := range []struct {
		name, body string
		want       int
	}{
		{"composed on page 1", `<PlaceObject><Box width="2" height="1"/></PlaceObject>` + slate + `<ClearPage/><PlaceObject slate="s"/>`, 2},
		{"composed before the first page", slate + `<PlaceObject><Box width="2" height="1"/></PlaceObject><ClearPage/><PlaceObject slate="s"/>`, 2},
		{"placed twice", slate + `<PlaceObject slate="s"/><ClearPage/><PlaceObject slate="s"/>`, 2},
	} {
		t.Run(c.name, func(t *testing.T) {
			marks, _ := runMarks(t, `<Layout xmlns="urn:speedata.de/2021/xts/en"><Record match="data">`+c.body+`</Record></Layout>`)
			if got, ok := marks["m"]; !ok || got != c.want {
				t.Errorf("mark m on page %d (%t), want %d", got, ok, c.want)
			}
		})
	}
}

// A Mark in an Action inside a Paragraph takes the page of the line it is in,
// also inside an inline element and in a paragraph that a flow breaks across
// two pages.
func TestMarkInAParagraphTakesThePageOfItsLine(t *testing.T) {
	mark := `<Action><Mark select="'m'"/></Action>`
	placed := func(par string) string {
		return `<Layout xmlns="urn:speedata.de/2021/xts/en"><Record match="data">
    <PlaceObject><Box width="2" height="1"/></PlaceObject>
    <ClearPage/>
    <PlaceObject><TextBlock>` + par + `</TextBlock></PlaceObject>
  </Record></Layout>`
	}
	// The page has room for nine lines, the first paragraph takes seven, so
	// the second one breaks after its second line.
	flow := func(second string) string {
		return flowMarksLayout(`<Paragraph style="orphans: 1; widows: 1">` + second + `</Paragraph>`)
	}
	for _, c := range []struct {
		name, layout string
		want         int
	}{
		{"in a placed paragraph", placed(`<Paragraph>` + mark + `<Value>A</Value></Paragraph>`), 2},
		{"in an inline element", placed(`<Paragraph><Value>A </Value><B>` + mark + `<Value>B</Value></B></Paragraph>`), 2},
		{"first line in a flow", flow(mark + lines(4)), 1},
		{"third line in a flow", flow(lines(2) + `<Br/>` + mark + lines(2)), 2},
	} {
		t.Run(c.name, func(t *testing.T) {
			marks, _ := runMarks(t, c.layout)
			if got, ok := marks["m"]; !ok || got != c.want {
				t.Errorf("mark m on page %d (%t), want %d", got, ok, c.want)
			}
		})
	}
}

// A Bookmark and a Mark with pdftarget inside a Paragraph point at their
// line: here the third one of a paragraph that a flow breaks after its
// second, so at the top line of page 2.
func TestBookmarkInAParagraphPointsAtItsLine(t *testing.T) {
	pdf := renderPDF(t, flowMarksLayout(`<Paragraph style="orphans: 1; widows: 1">`+lines(2)+`<Br/><Bookmark select="'Chapter'" level="1"/><Action><Mark select="'m'" pdftarget="yes"/></Action>`+lines(2)+`</Paragraph>`))
	named, outlines := pdfDests(t, pdf)
	// The grid's top on page 2: 60mm less the 10mm margin. bag puts a
	// destination in a line its height and depth above the baseline, so
	// within a line's 12pt of the top.
	top := 50 / 25.4 * 72
	for name, got := range map[string][2]float64{"named destination m": named["m"], "bookmark Chapter": outlines["Chapter"]} {
		if got[0] != 2 || math.Abs(got[1]-top) > 12 {
			t.Errorf("%s at page %v, y %.2f; want page 2, within 12pt of y %.2f", name, got[0], got[1], top)
		}
	}
}
