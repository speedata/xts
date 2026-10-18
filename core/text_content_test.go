package core

import (
	"strings"
	"testing"
)

// Text directly in a command counts as a <Value> with that text, also inside
// Switch and ForAll; text of whitespace only between commands is dropped,
// a lone space included.
func TestTextContentActsAsValue(t *testing.T) {
	log := runLayoutLog(t, layoutHead+`
  <Record match="data">
    <Message>direct text</Message>
    <Message>
      <Value>a</Value>
      <Value>b</Value>
    </Message>
    <Message><Switch><Case test="true()">case text</Case></Switch></Message>
    <Message><ForAll select="(1, 2)">n<Value select="."/>; </ForAll></Message>
    <Message><ForAll select="(1, 2)"><Value select="."/> </ForAll></Message>
    <PlaceObject><TextBlock><Paragraph>page</Paragraph></TextBlock></PlaceObject>
  </Record>
</Layout>`)
	for _, want := range []string{
		"Message (line 3): direct text",
		"Message (line 4): ab",
		"Message (line 8): case text",
		"Message (line 9): n1; n2;",
		"Message (line 10): 12",
	} {
		if !strings.Contains(log, want) {
			t.Errorf("log lacks %q:\n%s", want, log)
		}
	}
}
