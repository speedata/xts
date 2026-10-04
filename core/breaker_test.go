package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/boxesandglue/boxesandglue/backend/node"
	"github.com/boxesandglue/htmlbag"
)

func runWithBreakers(t *testing.T, layout string, breakers map[string]htmlbag.BreakerFunc) error {
	t.Helper()
	jobname := filepath.Join(t.TempDir(), "out")
	out, err := os.Create(jobname + ".pdf")
	if err != nil {
		t.Fatal(err)
	}
	return RunXTS(&XTSConfig{
		Datafile:    strings.NewReader("<data/>"),
		Layoutfile:  strings.NewReader(layout),
		FindFile:    FindFile,
		Jobname:     jobname,
		Outfile:     out,
		OutFilename: out.Name(),
		Breakers:    breakers,
	})
}

// recordBreaker breaks at every legal breakpoint and counts the paragraphs it
// breaks.
type recordBreaker struct{ paragraphs *int }

func (b recordBreaker) Breaks(p *node.BreakProblem) []int {
	*b.paragraphs++
	breaks := make([]int, len(p.Candidates))
	for i := range breaks {
		breaks[i] = i
	}
	return breaks
}

func TestBreakersReachTheCSSBuilder(t *testing.T) {
	var got []htmlbag.BreakerStyles
	var paragraphs int
	err := runWithBreakers(t, `<Layout xmlns="urn:speedata.de/2021/xts/en">
  <StyleSheet>p.probe { -bag-line-breaker: Probe; font-size: 10pt }</StyleSheet>
  <Record match="data">
    <PlaceObject><TextBlock><Paragraph class="probe"><Value>one two three</Value></Paragraph></TextBlock></PlaceObject>
  </Record>
</Layout>`, map[string]htmlbag.BreakerFunc{
		"probe": func(s htmlbag.BreakerStyles) node.Breaker {
			got = append(got, s)
			return recordBreaker{&paragraphs}
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) == 0 || paragraphs == 0 {
		t.Fatal("the registered breaker was not used")
	}
	if s := got[0]; s.Name != "probe" || s.FontSize.ToPT() != 10 {
		t.Errorf("styles = %+v, want probe at 10pt", s)
	}
}

func TestBreakersReservedNameIsAnError(t *testing.T) {
	err := runWithBreakers(t, `<Layout xmlns="urn:speedata.de/2021/xts/en"/>`, map[string]htmlbag.BreakerFunc{
		"auto": func(htmlbag.BreakerStyles) node.Breaker { return nil },
	})
	if err == nil || !strings.Contains(err.Error(), "reserved") {
		t.Fatalf("err = %v, want the reserved name refused", err)
	}
}
