package core

import (
	"strings"
	"testing"

	xpath "github.com/speedata/goxpath"
)

// TestLoopReturnsItsChildren checks that Loop passes on what its children
// return, so that rows or paragraphs made in a Loop reach the Table or
// Paragraph around it.
func TestLoopReturnsItsChildren(t *testing.T) {
	parser, err := xpath.NewParser(strings.NewReader("<root/>"))
	if err != nil {
		t.Fatal(err)
	}
	xd := &xtsDocument{data: parser}
	parser.Ctx.Store = map[any]any{"xd": xd}
	seq, err := cmdLoop(xd, parseLayoutElt(t, `<Loop select="3"><Value select="$_loopcounter"/></Loop>`))
	if err != nil {
		t.Fatal(err)
	}
	if len(seq) != 3 || seq.Stringvalue() != "123" {
		t.Errorf("Loop returned %d items %q, want 3 items \"123\"", len(seq), seq.Stringvalue())
	}
}
