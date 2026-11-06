package core

import (
	"strings"
	"testing"
)

const overprintLayout = `<Layout xmlns="urn:speedata.de/2021/xts/en">
	<DefineColor name="ink" model="cmyk" c="0" m="0" y="0" k="100" overprint="%s"/>
	<Record match="data">
		<PlaceObject><Box width="1" height="1" background-color="ink"/></PlaceObject>
	</Record>
</Layout>`

// The attribute overprint reaches the PDF as a graphics state, and only then.
func TestDefineColorOverprint(t *testing.T) {
	for _, tc := range []struct {
		value string
		want  bool
	}{{"yes", true}, {"no", false}} {
		pdf := renderPDF(t, strings.ReplaceAll(overprintLayout, "%s", tc.value))
		if got := strings.Contains(pdf, "/OPM 1"); got != tc.want {
			t.Errorf("overprint=%q: /OPM 1 in the PDF: %v", tc.value, got)
		}
	}
}

func TestDefineColorOverprintWarning(t *testing.T) {
	log := runLayoutLog(t, `<Layout xmlns="urn:speedata.de/2021/xts/en">
		<DefineColor name="ink" value="black" overprint="yes"/>
		<Record match="data">
			<PlaceObject><Box width="1" height="1" background-color="ink"/></PlaceObject>
		</Record>
	</Layout>`)
	if !strings.Contains(log, "knocks out") {
		t.Errorf("no warning for gray overprint in %q", log)
	}
}
