package core

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// runLayoutLog runs layout on <data/> and returns what it logged.
func runLayoutLog(t *testing.T, layout string) string {
	t.Helper()
	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})))
	defer slog.SetDefault(old)
	jobname := filepath.Join(t.TempDir(), "out")
	out, err := os.Create(jobname + ".pdf")
	if err != nil {
		t.Fatal(err)
	}
	err = RunXTS(&XTSConfig{
		Datafile:    strings.NewReader("<data/>"),
		Layoutfile:  strings.NewReader(layout),
		FindFile:    FindFile,
		Jobname:     jobname,
		Outfile:     out,
		OutFilename: out.Name(),
	})
	if err != nil {
		t.Fatal(err)
	}
	return buf.String()
}

const layoutHead = `<Layout xmlns="urn:speedata.de/2021/xts/en" xmlns:sd="urn:speedata.de/2021/xtsfunctions/en">`

// sd:current-row(), sd:number-of-columns() and sd:number-of-rows() before
// the first page set it up, as sd:current-page() does, and report that page:
// row 1, and the same grid the page has once an object is on it.
func TestGridFunctionsBeforeFirstPage(t *testing.T) {
	log := runLayoutLog(t, layoutHead+`
  <Record match="data">
    <Message select="concat('before: ', sd:current-row(), ' ', sd:number-of-columns(), ' ', sd:number-of-rows())"/>
    <PlaceObject><TextBlock><Paragraph><Value>Hello</Value></Paragraph></TextBlock></PlaceObject>
    <Message select="concat('after: ', sd:current-row(), ' ', sd:number-of-columns(), ' ', sd:number-of-rows(), ' page ', sd:current-page())"/>
  </Record>
</Layout>`)
	before := regexp.MustCompile(`before: (\d+) (\d+ \d+)`).FindStringSubmatch(log)
	after := regexp.MustCompile(`after: (\d+) (\d+ \d+) page (\d+)`).FindStringSubmatch(log)
	if before == nil || after == nil {
		t.Fatalf("messages missing in\n%s", log)
	}
	if before[1] != "1" {
		t.Errorf("sd:current-row() before the first page = %s, want 1", before[1])
	}
	if before[2] != after[2] {
		t.Errorf("columns and rows before the first page = %s, on the page = %s", before[2], after[2])
	}
	if after[3] != "1" {
		t.Errorf("the object went to page %s, want 1: asking created a page of its own", after[3])
	}
}

// A page type whose test asks for the grid is evaluated while the first page
// is set up, when there is no grid yet: the test fails with an error and the
// page gets another type, instead of XTS stopping with a panic.
func TestGridFunctionInPageTypeTest(t *testing.T) {
	log := runLayoutLog(t, layoutHead+`
  <DefineMasterPage name="first" test="sd:current-row() = 1" margin="1cm"/>
  <Record match="data">
    <PlaceObject><TextBlock><Paragraph><Value>Hello</Value></Paragraph></TextBlock></PlaceObject>
  </Record>
</Layout>`)
	if !strings.Contains(log, "sd:current-row(): no page yet") {
		t.Errorf("no error for the page type test in\n%s", log)
	}
	if !strings.Contains(log, `type="default page"`) {
		t.Errorf("the page did not fall back to the default page type in\n%s", log)
	}
}

// A page type test that does not give true or false is reported, and the
// page gets another type.
func TestPageTypeTestNotBoolean(t *testing.T) {
	log := runLayoutLog(t, layoutHead+`
  <DefineMasterPage name="first" test="'yes'" margin="1cm"/>
  <Record match="data">
    <PlaceObject><TextBlock><Paragraph><Value>Hello</Value></Paragraph></TextBlock></PlaceObject>
  </Record>
</Layout>`)
	if !strings.Contains(log, "does not give true or false") {
		t.Errorf("no error for the page type test in\n%s", log)
	}
}
