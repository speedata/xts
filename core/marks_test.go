package core

import "testing"

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
