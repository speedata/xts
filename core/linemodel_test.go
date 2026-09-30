package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/boxesandglue/boxesandglue/backend/node"
	"github.com/boxesandglue/htmlbag"
)

func runWithLineModels(t *testing.T, layout string, models map[string]htmlbag.LineModelFunc) error {
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
		LineModels:  models,
	})
}

func TestLineModelsReachTheCSSBuilder(t *testing.T) {
	var got []htmlbag.LineModelStyles
	err := runWithLineModels(t, `<Layout xmlns="urn:speedata.de/2021/xts/en" xmlns:h="http://www.w3.org/1999/xhtml">
  <Record match="data">
    <PlaceObject><HTML><h:p style="-bag-leading-model: Probe; font-size: 10pt; line-height: 14pt">Text</h:p></HTML></PlaceObject>
  </Record>
</Layout>`, map[string]htmlbag.LineModelFunc{
		"probe": func(s htmlbag.LineModelStyles) node.LineModel {
			got = append(got, s)
			return nil
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) == 0 {
		t.Fatal("the registered line model was not called")
	}
	if s := got[0]; s.Name != "probe" || s.FontSize.ToPT() != 10 || s.LineHeight.ToPT() != 14 {
		t.Errorf("styles = %+v, want probe at 10pt/14pt", s)
	}
}

func TestLineModelsReservedNameIsAnError(t *testing.T) {
	err := runWithLineModels(t, `<Layout xmlns="urn:speedata.de/2021/xts/en"/>`, map[string]htmlbag.LineModelFunc{
		"half": func(htmlbag.LineModelStyles) node.LineModel { return nil },
	})
	if err == nil || !strings.Contains(err.Error(), "reserved") {
		t.Fatalf("err = %v, want the reserved name refused", err)
	}
}
