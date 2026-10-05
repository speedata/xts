package core

import (
	"strings"
	"testing"

	xpath "github.com/speedata/goxpath"
)

// runForAll runs the ForAll snippet over five records and returns $out.
func runForAll(t *testing.T, snippet string) string {
	t.Helper()
	return runCommand(t, cmdForall, snippet)
}

// commandDocument defines the records over five records, with the root as
// the context and $out empty.
func commandDocument(t *testing.T, records ...string) *xtsDocument {
	t.Helper()
	parser, err := xpath.NewParser(strings.NewReader("<root><rec/><rec/><rec/><rec/><rec/></root>"))
	if err != nil {
		t.Fatal(err)
	}
	xd := &xtsDocument{data: parser}
	parser.Ctx.Store = map[any]any{"xd": xd}
	xd.data.SetVariable("out", xpath.Sequence{""})
	for _, rec := range records {
		if _, err := cmdRecord(xd, parseLayoutElt(t, rec)); err != nil {
			t.Fatalf("cmdRecord: %v", err)
		}
	}
	if _, err := parser.Evaluate("/root"); err != nil {
		t.Fatal(err)
	}
	return xd
}

// runCommand defines the records, runs the command snippet over five records
// and returns $out.
func runCommand(t *testing.T, cmd commandFunc, snippet string, records ...string) string {
	t.Helper()
	xd := commandDocument(t, records...)
	if _, err := cmd(xd, parseLayoutElt(t, snippet)); err != nil {
		t.Fatalf("%s: %v", snippet, err)
	}
	got, _ := xd.data.GetVariable("out")
	return got.Stringvalue()
}

// TestForAllPositionAfterPath checks that a path of two or more steps in the
// body of a ForAll leaves position() and last() of the loop alone.
func TestForAllPositionAfterPath(t *testing.T) {
	const snippet = `<ForAll select="rec">
		<SetVariable variable="n" select="count(../rec)"/>
		<SetVariable variable="out" select="concat($out, position(), '/', last(), ' ')"/>
	</ForAll>`
	if got, want := runForAll(t, snippet), "1/5 2/5 3/5 4/5 5/5 "; got != want {
		t.Errorf("position()/last() = %q, want %q", got, want)
	}
}

// TestForAllRestoresPosition checks that an inner ForAll gives the outer one
// its position() and last() back.
func TestForAllRestoresPosition(t *testing.T) {
	const snippet = `<ForAll select="rec">
		<ForAll select="../rec[position() &lt; 3]"/>
		<SetVariable variable="out" select="concat($out, position(), '/', last(), ' ')"/>
	</ForAll>`
	if got, want := runForAll(t, snippet), "1/5 2/5 3/5 4/5 5/5 "; got != want {
		t.Errorf("position()/last() = %q, want %q", got, want)
	}
}

// TestProcessNodePositionAndLast checks that position() and last() in a
// Record called by ProcessNode count the selected nodes.
func TestProcessNodePositionAndLast(t *testing.T) {
	const record = `<Record match="rec">
		<SetVariable variable="out" select="concat($out, position(), '/', last(), ' ')"/>
	</Record>`
	if got, want := runCommand(t, cmdProcessNode, `<ProcessNode select="rec"/>`, record), "1/5 2/5 3/5 4/5 5/5 "; got != want {
		t.Errorf("position()/last() = %q, want %q", got, want)
	}
}

// TestProcessNodeRestoresPosition checks that a ProcessNode gives the ForAll
// around it its position() and last() back.
func TestProcessNodeRestoresPosition(t *testing.T) {
	const snippet = `<ForAll select="rec">
		<ProcessNode select="../rec[position() &lt; 3]" mode="inner"/>
		<SetVariable variable="out" select="concat($out, position(), '/', last(), ' ')"/>
	</ForAll>`
	const record = `<Record match="rec" mode="inner"/>`
	if got, want := runCommand(t, cmdForall, snippet, record), "1/5 2/5 3/5 4/5 5/5 "; got != want {
		t.Errorf("position()/last() = %q, want %q", got, want)
	}
}

// TestProcessNodeReturnsWhatItsRecordsReturn checks that ProcessNode gives
// its caller what the records return, in order, as a ForAll does.
func TestProcessNodeReturnsWhatItsRecordsReturn(t *testing.T) {
	const snippet = `<SetVariable variable="out"><ProcessNode select="rec"/></SetVariable>`
	const record = `<Record match="rec"><Value select="position()"/></Record>`
	if got, want := runCommand(t, cmdSetVariable, snippet, record), "12345"; got != want {
		t.Errorf("$out = %q, want %q", got, want)
	}
}

// TestProcessNodeStopsAtAnError checks that an error in a record ends the
// ProcessNode at once and reaches its caller, with position(), last() and the
// context given back.
func TestProcessNodeStopsAtAnError(t *testing.T) {
	const record = `<Record match="rec">
		<SetVariable variable="out" select="concat($out, position(), ' ')"/>
		<ForAll select="nosuchfunction()"/>
	</Record>`
	xd := commandDocument(t, record)
	ctx := xd.data.Ctx
	oldContext := ctx.GetContextSequence()
	ctx.Pos = 2
	ctx.SetSize(3)
	if _, err := cmdProcessNode(xd, parseLayoutElt(t, `<ProcessNode select="rec"/>`)); err == nil {
		t.Fatal("no error from the record")
	}
	if got, _ := xd.data.GetVariable("out"); got.Stringvalue() != "1 " {
		t.Errorf("$out = %q, want the first record only", got.Stringvalue())
	}
	if ctx.Pos != 2 || ctx.Size() != 3 {
		t.Errorf("position()/last() = %d/%d after the error, want 2/3", ctx.Pos, ctx.Size())
	}
	if got := ctx.GetContextSequence(); len(got) != 1 || len(oldContext) != 1 || got[0] != oldContext[0] {
		t.Errorf("context after the error = %v, want %v", got, oldContext)
	}
}

// TestEvaluateXPathRestoresNamespaces checks that an expression evaluated in
// the middle of another one, as when sd:current-page() sets up the page, gives
// the outer expression its namespaces back.
func TestEvaluateXPathRestoresNamespaces(t *testing.T) {
	parser, err := xpath.NewParser(strings.NewReader("<root/>"))
	if err != nil {
		t.Fatal(err)
	}
	xd := &xtsDocument{data: parser}
	outer := map[string]string{"sd": "urn:speedata.de/2021/xtsfunctions/en"}
	parser.Ctx.Namespaces = outer
	if _, err := evaluateXPath(xd, nil, "true()"); err != nil {
		t.Fatal(err)
	}
	if got := parser.Ctx.Namespaces["sd"]; got != outer["sd"] {
		t.Errorf("namespace for sd after evaluateXPath = %q, want %q", got, outer["sd"])
	}
}
