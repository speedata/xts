package core

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/speedata/goxml"
	xpath "github.com/speedata/goxpath"
)

func attr(local, value string) xml.Attr {
	return xml.Attr{Name: xml.Name{Local: local}, Value: value}
}

// TestColumnIsQueryable verifies that cmdColumn returns a queryable
// goxml.Element (not an opaque html node) and that the consume-boundary
// conversion produces the <col data-width="…"> shape htmlbag expects.
func TestColumnIsQueryable(t *testing.T) {
	xd := &xtsDocument{}
	layoutelt := goxml.NewElement()
	layoutelt.Name = "Column"
	layoutelt.SetAttribute(xml.Attr{Name: xml.Name{Local: "width"}, Value: "3cm"})

	seq, err := cmdColumn(xd, layoutelt)
	if err != nil {
		t.Fatalf("cmdColumn returned error: %v", err)
	}
	if len(seq) != 1 {
		t.Fatalf("cmdColumn returned %d items, want 1", len(seq))
	}
	col, ok := seq[0].(*goxml.Element)
	if !ok {
		t.Fatalf("cmdColumn returned %T, want *goxml.Element", seq[0])
	}
	if col.Name != "Column" {
		t.Errorf("element name = %q, want Column", col.Name)
	}
	attrs := col.Attributes()
	if len(attrs) != 1 || attrs[0].Name != "width" || attrs[0].Value != "3cm" {
		t.Errorf("attributes = %v, want width=3cm", attrs)
	}

	hn := goxmlToHTMLNode(col)
	if hn.Data != "col" {
		t.Errorf("converted node data = %q, want col", hn.Data)
	}
	if len(hn.Attr) != 1 || hn.Attr[0].Key != "data-width" || hn.Attr[0].Val != "3cm" {
		t.Errorf("converted node attrs = %v, want data-width=3cm", hn.Attr)
	}
}

// TestDispatchValueContextRejectsActions verifies that an action command is
// rejected when dispatched in a value context (e.g. inside a bound variable).
func TestDispatchValueContextRejectsActions(t *testing.T) {
	xd := &xtsDocument{}
	parent := goxml.NewElement()
	parent.Name = "SetVariable"
	child := goxml.NewElement()
	child.Name = "PlaceObject"
	parent.Append(child)

	if _, err := dispatchValueContext(xd, parent); err == nil {
		t.Error("dispatchValueContext accepted an action, want error")
	}
	if xd.valueContext {
		t.Error("valueContext flag was not restored after dispatch")
	}
}

// TestValueContextRecursesThroughConstructors verifies that the value-context
// restriction reaches actions nested inside a constructor (here <Element>).
func TestValueContextRecursesThroughConstructors(t *testing.T) {
	xd := &xtsDocument{}
	parent := goxml.NewElement()
	parent.Name = "SetVariable"
	elt := goxml.NewElement()
	elt.Name = "Element"
	elt.SetAttribute(xml.Attr{Name: xml.Name{Local: "name"}, Value: "x"})
	action := goxml.NewElement()
	action.Name = "PlaceObject"
	elt.Append(action)
	parent.Append(elt)

	if _, err := dispatchValueContext(xd, parent); err == nil {
		t.Error("nested action inside a constructor was accepted, want error")
	}
}

// TestSetVariableAsRejectsAction verifies that <SetVariable as="…"> with an
// action in its body is rejected.
func TestSetVariableAsRejectsAction(t *testing.T) {
	xd := &xtsDocument{}
	sv := goxml.NewElement()
	sv.Name = "SetVariable"
	sv.SetAttribute(xml.Attr{Name: xml.Name{Local: "variable"}, Value: "x"})
	sv.SetAttribute(xml.Attr{Name: xml.Name{Local: "as"}, Value: "element()*"})
	action := goxml.NewElement()
	action.Name = "PlaceObject"
	sv.Append(action)

	if _, err := cmdSetVariable(xd, sv); err == nil {
		t.Error("SetVariable with as= accepted an action body, want error")
	}
}

// TestConstructorAllowedInValueContext verifies that pure constructors are not
// rejected in a value context.
func TestConstructorAllowedInValueContext(t *testing.T) {
	xd := &xtsDocument{}
	parent := goxml.NewElement()
	parent.Name = "SetVariable"
	col := goxml.NewElement()
	col.Name = "Column"
	col.SetAttribute(xml.Attr{Name: xml.Name{Local: "width"}, Value: "1cm"})
	parent.Append(col)

	seq, err := dispatchValueContext(xd, parent)
	if err != nil {
		t.Fatalf("constructor rejected in value context: %v", err)
	}
	if len(seq) != 1 {
		t.Fatalf("got %d items, want 1", len(seq))
	}
}

// TestCallTemplateParamScoping verifies that <CallTemplate> binds parameters in
// the template body and restores the previous binding afterwards.
func TestCallTemplateParamScoping(t *testing.T) {
	parser, err := xpath.NewParser(strings.NewReader("<root/>"))
	if err != nil {
		t.Fatal(err)
	}
	xd := &xtsDocument{data: parser, templates: map[string]*goxml.Element{}}

	// Template "t": its body captures the current value of $p into $captured.
	tmpl := goxml.NewElement()
	tmpl.Name = "Template"
	tmpl.SetAttribute(attr("name", "t"))
	body := goxml.NewElement()
	body.Name = "SetVariable"
	body.SetAttribute(attr("variable", "captured"))
	body.SetAttribute(attr("select", "$p"))
	tmpl.Append(body)
	if _, err := cmdTemplate(xd, tmpl); err != nil {
		t.Fatalf("cmdTemplate: %v", err)
	}

	// Outer binding of $p.
	xd.data.SetVariable("p", xpath.Sequence{"outer"})

	// CallTemplate t with Param p = 'inner'.
	call := goxml.NewElement()
	call.Name = "CallTemplate"
	call.SetAttribute(attr("name", "t"))
	param := goxml.NewElement()
	param.Name = "Param"
	param.SetAttribute(attr("name", "p"))
	param.SetAttribute(attr("select", "'inner'"))
	call.Append(param)

	if _, err := cmdCallTemplate(xd, call); err != nil {
		t.Fatalf("cmdCallTemplate: %v", err)
	}

	if got, _ := xd.data.GetVariable("captured"); got.Stringvalue() != "inner" {
		t.Errorf("captured = %q, want inner (param bound inside body)", got.Stringvalue())
	}
	if got, _ := xd.data.GetVariable("p"); got.Stringvalue() != "outer" {
		t.Errorf("p after call = %q, want outer (binding restored)", got.Stringvalue())
	}
}

// TestCallTemplateUnknown verifies that calling an undefined template errors.
func TestCallTemplateUnknown(t *testing.T) {
	xd := &xtsDocument{templates: map[string]*goxml.Element{}}
	call := goxml.NewElement()
	call.Name = "CallTemplate"
	call.SetAttribute(attr("name", "missing"))
	if _, err := cmdCallTemplate(xd, call); err == nil {
		t.Error("calling an unknown template did not error")
	}
}

// TestCallTemplateRejectedInValueContext verifies CallTemplate is an action and
// thus not allowed in a value context.
func TestCallTemplateRejectedInValueContext(t *testing.T) {
	xd := &xtsDocument{templates: map[string]*goxml.Element{}}
	sv := goxml.NewElement()
	sv.Name = "SetVariable"
	sv.SetAttribute(attr("variable", "x"))
	sv.SetAttribute(attr("as", "item()*"))
	call := goxml.NewElement()
	call.Name = "CallTemplate"
	call.SetAttribute(attr("name", "t"))
	sv.Append(call)
	if _, err := cmdSetVariable(xd, sv); err == nil {
		t.Error("CallTemplate accepted in value context, want error")
	}
}
