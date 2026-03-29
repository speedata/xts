package core

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/speedata/goxml"
	xpath "github.com/speedata/goxpath"
	"golang.org/x/net/html"
)

// findRecordByName finds a Record matching the element name with default mode
// and no predicate conditions (used for root element dispatch).
func findRecordByName(elemName string) *goxml.Element {
	for i := len(dataRecords) - 1; i >= 0; i-- {
		rec := dataRecords[i]
		if rec.elemName == elemName && rec.mode == "" && rec.pred == "" {
			return rec.layout
		}
	}
	return nil
}

// parseMatch splits a match expression like "foo[not(@bar='baz')]" into
// element name ("foo") and predicate ("not(@bar='baz')").
func parseMatch(match string) (string, string) {
	idx := strings.Index(match, "[")
	if idx == -1 {
		return match, ""
	}
	elemName := match[:idx]
	pred := match[idx+1 : len(match)-1] // strip [ and ]
	return elemName, pred
}

// expandTextValueTemplates evaluates {expr} expressions in text content (XSLT 3.0 style).
// Escaped braces {{ and }} are converted to literal { and }.
func expandTextValueTemplates(xd *xtsDocument, layoutelt *goxml.Element, str string) string {
	// First, replace escaped braces with placeholders
	const leftPlaceholder = "\x00LEFT_BRACE\x00"
	const rightPlaceholder = "\x00RIGHT_BRACE\x00"

	str = strings.ReplaceAll(str, "{{", leftPlaceholder)
	str = strings.ReplaceAll(str, "}}", rightPlaceholder)

	// Evaluate {expr} expressions
	str = attributeValueRE.ReplaceAllStringFunc(str, func(match string) string {
		// Strip curly braces to get the expression
		expr := match[1 : len(match)-1]
		seq, err := evaluateXPath(xd, layoutelt.Namespaces, expr)
		if err != nil {
			slog.Error(fmt.Sprintf("HTML expand-text (line %d): error evaluating expression {%s}: %s", layoutelt.Line, expr, err))
			return ""
		}
		return seq.Stringvalue()
	})

	// Replace placeholders with literal braces
	str = strings.ReplaceAll(str, leftPlaceholder, "{")
	str = strings.ReplaceAll(str, rightPlaceholder, "}")

	return str
}

// hasMixedXHTMLContent checks if any child element uses the XHTML namespace.
func hasMixedXHTMLContent(elt *goxml.Element) bool {
	for _, cld := range elt.Children() {
		if child, ok := cld.(*goxml.Element); ok {
			if ns := child.Namespaces[child.Prefix]; ns == XHTMLNAMESPACE {
				return true
			}
		}
	}
	return false
}

// convertXHTMLElement recursively converts a goxml Element in the XHTML
// namespace to an *html.Node tree. XTS-namespace children are dispatched
// and their results inserted into the tree.
func (xd *xtsDocument) convertXHTMLElement(elt *goxml.Element) (*html.Node, error) {
	n := &html.Node{
		Type: html.ElementNode,
		Data: strings.ToLower(elt.Name),
	}
	for _, attr := range elt.Attributes() {
		if strings.HasPrefix(attr.Name, "xmlns") {
			continue
		}
		n.Attr = append(n.Attr, html.Attribute{Key: attr.Name, Val: attr.Value})
	}
	for _, cld := range elt.Children() {
		switch t := cld.(type) {
		case goxml.CharData:
			n.AppendChild(&html.Node{Type: html.TextNode, Data: t.Contents})
		case *goxml.Element:
			if ns := t.Namespaces[t.Prefix]; ns == XHTMLNAMESPACE {
				child, err := xd.convertXHTMLElement(t)
				if err != nil {
					return nil, err
				}
				n.AppendChild(child)
			} else {
				// XTS command: dispatch and insert results
				if f, ok := dispatchTable[t.Name]; ok {
					seq, err := f(xd, t)
					if err != nil {
						return nil, err
					}
					for _, itm := range seq {
						switch v := itm.(type) {
						case *html.Node:
							n.AppendChild(v)
						case *goxml.Element:
							n.AppendChild(&html.Node{Type: html.TextNode, Data: v.Stringvalue()})
						case string:
							n.AppendChild(&html.Node{Type: html.TextNode, Data: v})
						}
					}
				} else {
					return nil, fmt.Errorf("HTML (line %d): unknown element %q", t.Line, t.Name)
				}
			}
		}
	}
	return n, nil
}

// buildHTMLFromMixedContent processes <HTML> with mixed XHTML literals
// and XTS commands, returning a sequence of *html.Node.
func (xd *xtsDocument) buildHTMLFromMixedContent(layoutelt *goxml.Element) (xpath.Sequence, error) {
	var result xpath.Sequence
	for _, cld := range layoutelt.Children() {
		switch t := cld.(type) {
		case goxml.CharData:
			s := strings.TrimSpace(t.Contents)
			if s != "" {
				result = append(result, &html.Node{Type: html.TextNode, Data: t.Contents})
			}
		case *goxml.Element:
			if ns := t.Namespaces[t.Prefix]; ns == XHTMLNAMESPACE {
				n, err := xd.convertXHTMLElement(t)
				if err != nil {
					return nil, err
				}
				result = append(result, n)
			} else {
				// XTS command
				if f, ok := dispatchTable[t.Name]; ok {
					seq, err := f(xd, t)
					if err != nil {
						return nil, err
					}
					result = append(result, seq...)
				} else {
					return nil, fmt.Errorf("HTML (line %d): unknown element %q", t.Line, t.Name)
				}
			}
		}
	}
	return result, nil
}
