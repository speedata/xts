package core

import (
	"strings"
	"testing"

	"github.com/speedata/goxml"
	"golang.org/x/net/html"
)

// parseValueElt parses a single <Value> layout snippet into its goxml element.
func parseValueElt(t *testing.T, snippet string) *goxml.Element {
	t.Helper()
	doc, err := goxml.Parse(strings.NewReader(snippet))
	if err != nil {
		t.Fatalf("parse %q: %v", snippet, err)
	}
	root, err := doc.Root()
	if err != nil {
		t.Fatalf("root %q: %v", snippet, err)
	}
	return root
}

// TestValuePreservesInlineBr verifies that <Value>foo<br/>bar</Value> keeps the
// <br/> as an inline node instead of flattening the subtree to "foobar".
// Regression: cmdValue used Stringvalue() unconditionally, silently dropping
// inline markup (the mailmerge company block lost every <br/>).
func TestValuePreservesInlineBr(t *testing.T) {
	xd := &xtsDocument{}
	seq, err := cmdValue(xd, parseValueElt(t, `<Value>foo<br/>bar</Value>`))
	if err != nil {
		t.Fatalf("cmdValue: %v", err)
	}
	if len(seq) != 3 {
		t.Fatalf("got %d items, want 3 (\"foo\", <br>, \"bar\"): %#v", len(seq), seq)
	}
	if s, ok := seq[0].(string); !ok || s != "foo" {
		t.Errorf("item 0 = %#v, want string \"foo\"", seq[0])
	}
	if n, ok := seq[1].(*html.Node); !ok || n.Data != "br" {
		t.Errorf("item 1 = %#v, want *html.Node <br>", seq[1])
	}
	if s, ok := seq[2].(string); !ok || s != "bar" {
		t.Errorf("item 2 = %#v, want string \"bar\"", seq[2])
	}
}

// TestValuePureTextStaysSingleString verifies the no-element fast path is
// unchanged: a text-only Value still returns one concatenated string.
func TestValuePureTextStaysSingleString(t *testing.T) {
	xd := &xtsDocument{}
	seq, err := cmdValue(xd, parseValueElt(t, `<Value>hello world</Value>`))
	if err != nil {
		t.Fatalf("cmdValue: %v", err)
	}
	if len(seq) != 1 {
		t.Fatalf("got %d items, want 1: %#v", len(seq), seq)
	}
	if s, ok := seq[0].(string); !ok || s != "hello world" {
		t.Errorf("item 0 = %#v, want string \"hello world\"", seq[0])
	}
}
