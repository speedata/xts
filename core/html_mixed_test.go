package core

import (
	"strings"
	"testing"

	"github.com/speedata/goxml"
	xpath "github.com/speedata/goxpath"
	"golang.org/x/net/html"
)

func parseLayoutElt(t *testing.T, snippet string) *goxml.Element {
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

func TestMixedContentExpandsText(t *testing.T) {
	const snippet = `<HTML xmlns:h="http://www.w3.org/1999/xhtml" expand-text="yes"><h:div style="width:{2+3}pt">{1+1}</h:div></HTML>`
	parser, err := xpath.NewParser(strings.NewReader("<root/>"))
	if err != nil {
		t.Fatal(err)
	}
	xd := &xtsDocument{data: parser}
	seq, err := cmdHTML(xd, parseLayoutElt(t, snippet))
	if err != nil {
		t.Fatalf("cmdHTML: %v", err)
	}
	if len(seq) != 1 {
		t.Fatalf("got %d items, want 1: %#v", len(seq), seq)
	}
	n, ok := seq[0].(*html.Node)
	if !ok {
		t.Fatalf("item 0 = %#v, want *html.Node", seq[0])
	}
	var style string
	for _, a := range n.Attr {
		if a.Key == "style" {
			style = a.Val
		}
	}
	if style != "width:5pt" {
		t.Errorf("style = %q, want %q", style, "width:5pt")
	}
	if n.FirstChild == nil || n.FirstChild.Data != "2" {
		t.Errorf("text = %#v, want \"2\"", n.FirstChild)
	}
}

// TestMixedContentWithoutExpandText pins the other side: without the attribute
// the braces are literal, as they are in the non-mixed path.
func TestMixedContentWithoutExpandText(t *testing.T) {
	const snippet = `<HTML xmlns:h="http://www.w3.org/1999/xhtml"><h:div style="width:{2+3}pt">{1+1}</h:div></HTML>`
	parser, err := xpath.NewParser(strings.NewReader("<root/>"))
	if err != nil {
		t.Fatal(err)
	}
	xd := &xtsDocument{data: parser}
	seq, err := cmdHTML(xd, parseLayoutElt(t, snippet))
	if err != nil {
		t.Fatalf("cmdHTML: %v", err)
	}
	n, ok := seq[0].(*html.Node)
	if !ok {
		t.Fatalf("item 0 = %#v, want *html.Node", seq[0])
	}
	for _, a := range n.Attr {
		if a.Key == "style" && a.Val != "width:{2+3}pt" {
			t.Errorf("style = %q, want the braces left alone", a.Val)
		}
	}
	if n.FirstChild == nil || n.FirstChild.Data != "{1+1}" {
		t.Errorf("text = %#v, want \"{1+1}\"", n.FirstChild)
	}
}
