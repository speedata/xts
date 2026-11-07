package core

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

const auxLayoutHead = `<Layout xmlns="urn:speedata.de/2021/xts/en" xmlns:sd="urn:speedata.de/2021/xtsfunctions/en"><Record match="data">`

// runTwice runs layout twice with the same jobname, so the second run reads the
// aux file of the first one. It returns what each run found out of date and the
// aux file after each run.
func runTwice(t *testing.T, body string) (outOfDate [2][]string, aux [2]string) {
	t.Helper()
	jobname := filepath.Join(t.TempDir(), "out")
	for i := range outOfDate {
		out, err := os.Create(jobname + ".pdf")
		if err != nil {
			t.Fatal(err)
		}
		xc := &XTSConfig{
			Datafile:    strings.NewReader("<data/>"),
			Layoutfile:  strings.NewReader(auxLayoutHead + body + `</Record></Layout>`),
			FindFile:    FindFile,
			Jobname:     jobname,
			Outfile:     out,
			OutFilename: out.Name(),
		}
		if err := RunXTS(xc); err != nil {
			t.Fatal(err)
		}
		outOfDate[i] = xc.Info.OutOfDate
		data, err := os.ReadFile(jobname + "-aux.xml")
		if err != nil {
			t.Fatal(err)
		}
		aux[i] = string(data)
	}
	return outOfDate, aux
}

func text(s string) string {
	return `<PlaceObject><TextBlock><Paragraph>` + s + `</Paragraph></TextBlock></PlaceObject>`
}

func mark(name string) string {
	return text(`<Action><Mark select="'` + name + `'"/></Action><Value>A</Value>`)
}

func TestOutOfDate(t *testing.T) {
	const newPage = `<ClearPage/>`
	for _, c := range []struct {
		name, body  string
		first, then []string
	}{
		{
			name:  "a mark further on is read before it exists",
			body:  text(`<Value select="sd:page-number('m')"/>`) + newPage + mark("m"),
			first: []string{"m"},
		},
		{
			name: "a mark before is known in the same run",
			body: mark("m") + newPage + text(`<Value select="sd:page-number('m')"/>`),
		},
		{
			name: "a mark nobody reads does not count",
			body: mark("m"),
		},
		{
			name: "a mark that does not exist never changes",
			body: text(`<Value select="sd:page-number('nope')"/>`),
		},
		{
			name:  "the number of pages is read before it is known",
			body:  text(`<Value select="sd:last-page-number()"/>`),
			first: []string{"sd:last-page-number()"},
		},
		{
			name:  "several references come sorted and once",
			body:  text(`<Value select="sd:page-number('b')"/><Value select="sd:page-number('a')"/><Value select="sd:page-number('b')"/>`) + newPage + mark("a") + mark("b"),
			first: []string{"a", "b"},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			got, _ := runTwice(t, c.body)
			if !slices.Equal(got[0], c.first) {
				t.Errorf("first run %v, want %v", got[0], c.first)
			}
			if !slices.Equal(got[1], c.then) {
				t.Errorf("second run %v, want %v", got[1], c.then)
			}
		})
	}
}

// A later read of a mark must not hide that an earlier one got an old value.
func TestOutOfDateKeepsTheEarlyRead(t *testing.T) {
	read := text(`<Value select="sd:page-number('m')"/>`)
	got, _ := runTwice(t, read+`<ClearPage/>`+mark("m")+read)
	if !slices.Equal(got[0], []string{"m"}) {
		t.Errorf("first run %v, want m", got[0])
	}
}

// Two old values of one mark are one reference. The second Mark of the same
// name wins, so the mark is on page 2 and both reads are out of date.
func TestOutOfDateNamesAMarkOnce(t *testing.T) {
	read := text(`<Value select="sd:page-number('m')"/>`)
	got, _ := runTwice(t, read+mark("m")+`<ClearPage/>`+read+mark("m"))
	if !slices.Equal(got[0], []string{"m"}) {
		t.Errorf("first run %v, want m once", got[0])
	}
}

// The aux file is the same for the same result, the comparison of runs needs
// that. With a map in random order, ten marks differ all the time.
func TestAuxFileIsSorted(t *testing.T) {
	var body strings.Builder
	for i := range 10 {
		body.WriteString(mark(fmt.Sprintf("m%d", 9-i)))
	}
	_, aux := runTwice(t, body.String())
	if aux[0] != aux[1] {
		t.Errorf("two runs write different files:\n%s\n%s", aux[0], aux[1])
	}
	var a struct {
		Marks []struct {
			Name string `xml:"name,attr"`
		} `xml:"marker>mark"`
	}
	if err := xml.Unmarshal([]byte(aux[0]), &a); err != nil {
		t.Fatal(err)
	}
	if !slices.IsSortedFunc(a.Marks, func(x, y struct {
		Name string `xml:"name,attr"`
	}) int {
		return strings.Compare(x.Name, y.Name)
	}) {
		t.Errorf("marks are not sorted: %v", a.Marks)
	}
}

func TestAuxFileKeepsPDFTarget(t *testing.T) {
	data, err := xml.Marshal(auxfile{Marker: mapmarker{
		"a": {name: "a", pagenumber: 3, id: 2, pdftarget: true},
		"b": {name: "b", pagenumber: 4, id: 1},
	}})
	if err != nil {
		t.Fatal(err)
	}
	got := auxfile{Marker: mapmarker{}}
	if err := xml.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if want := (marker{name: "a", pagenumber: 3, id: 2, pdftarget: true}); got.Marker["a"] != want {
		t.Errorf("a is %+v, want %+v", got.Marker["a"], want)
	}
	if got.Marker["b"].pdftarget {
		t.Error("b came back as a target")
	}
}
