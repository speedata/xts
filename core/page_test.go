package core

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

const markerLayout = `<Layout xmlns="urn:speedata.de/2021/xts/en"
    xmlns:sd="urn:speedata.de/2021/xtsfunctions/en">
  <Record match="data">
    <Loop select="10">
      <PlaceObject>
        <TextBlock>
          <Action><Mark select="concat('page', sd:current-page())"/></Action>
          <Paragraph><Value>Text</Value></Paragraph>
        </TextBlock>
      </PlaceObject>
      <ClearPage/>
    </Loop>
  </Record>
</Layout>`

func renderMarkerLayout(t *testing.T) {
	t.Helper()
	jobname := filepath.Join(t.TempDir(), "out")
	out, err := os.Create(jobname + ".pdf")
	if err != nil {
		t.Fatal(err)
	}
	err = RunXTS(&XTSConfig{
		Datafile:    strings.NewReader("<data/>"),
		Layoutfile:  strings.NewReader(markerLayout),
		FindFile:    FindFile,
		Jobname:     jobname,
		Outfile:     out,
		OutFilename: out.Name(),
	})
	if err != nil {
		t.Fatal(err)
	}
}

// A render leaves no goroutine behind. Each page used to start one that fed
// it marker ids and never ended, and the page it held kept the whole document
// alive.
func TestRenderLeavesNoGoroutine(t *testing.T) {
	renderMarkerLayout(t)
	before := runtime.NumGoroutine()
	renderMarkerLayout(t)
	deadline := time.Now().Add(time.Second)
	for runtime.NumGoroutine() > before && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if n := runtime.NumGoroutine(); n > before {
		t.Errorf("%d goroutines after rendering 10 pages, %d before", n, before)
	}
}
