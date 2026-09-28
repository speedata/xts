package core

import (
	"strings"
	"testing"

	xpath "github.com/speedata/goxpath"
)

// runForAll runs the ForAll snippet over five records and returns $out.
func runForAll(t *testing.T, snippet string) string {
	t.Helper()
	parser, err := xpath.NewParser(strings.NewReader("<root><rec/><rec/><rec/><rec/><rec/></root>"))
	if err != nil {
		t.Fatal(err)
	}
	xd := &xtsDocument{data: parser}
	parser.Ctx.Store = map[any]any{"xd": xd}
	xd.data.SetVariable("out", xpath.Sequence{""})
	if _, err := parser.Evaluate("/root"); err != nil {
		t.Fatal(err)
	}
	if _, err := cmdForall(xd, parseLayoutElt(t, snippet)); err != nil {
		t.Fatalf("cmdForall: %v", err)
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
