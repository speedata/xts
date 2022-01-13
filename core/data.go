package core

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/document"
	"github.com/speedata/goxml"
	"github.com/speedata/goxpath/xpath"
)

// Get the values from the child elements of B, Paragraph and its ilk and fill
// the provided typesetting element to get a recursive data structure.
func getTextvalues(te *document.TypesettingElement, seq xpath.Sequence, cmdname string) {
	for _, itm := range seq {
		switch t := itm.(type) {
		case *goxml.Element:
			te.Items = append(te.Items, t.Stringvalue())
		case float64:
			te.Items = append(te.Items, strconv.FormatFloat(t, 'f', -1, 64))
		case goxml.CharData:
			te.Items = append(te.Items, string(t))
		case string:
			te.Items = append(te.Items, t)
		case *document.TypesettingElement:
			te.Items = append(te.Items, t)
		default:
			bag.Logger.DPanicf("%s: unknown type %T", cmdname, t)
		}
	}
}

func getStructTag(f reflect.StructField, tagName string) string {
	return string(f.Tag.Get(tagName))
}

var (
	dummyBool        bool
	dummyStr         string
	boolType         = reflect.TypeOf(true)
	boolPtrType      = reflect.TypeOf(&dummyBool)
	stringType       = reflect.TypeOf("")
	stringPtrType    = reflect.TypeOf(&dummyStr)
	intType          = reflect.TypeOf(0)
	scaledPointsType = reflect.TypeOf(bag.ScaledPoint(0))
)

// getXMLAtttributes fills the struct at v with the attribute values of the current element.
func getXMLAtttributes(xd *xtsDocument, layoutelt *goxml.Element, v interface{}) error {
	attributes := make(map[string]string)
	for _, attrib := range layoutelt.Attributes() {
		attributes[attrib.Name] = attrib.Value
	}

	val := reflect.ValueOf(v)

	// If it's an interface or a pointer, unwrap it.
	if val.Kind() == reflect.Ptr && val.Elem().Kind() == reflect.Struct {
		val = val.Elem()
	} else {
		return fmt.Errorf("s must be a struct")
	}

	var valNumFields = val.NumField()

	var mustexist bool
	var dflt string
	var allowXPath bool
	var attValue string

	for i := 0; i < valNumFields; i++ {
		mustexist = false
		allowXPath = true
		dflt = ""

		field := val.Field(i)
		structField := val.Type().Field(i)
		for _, tag := range strings.Split(getStructTag(structField, "sdxml"), ",") {
			if strings.HasPrefix(tag, "default:") {
				dflt = strings.TrimPrefix(tag, "default:")
			} else if tag == "mustexist" {
				mustexist = true
			} else if tag == "noescape" {
				allowXPath = false
			}
		}
		fieldName := strings.ToLower(structField.Name)
		hasAttribute := false
		if a, ok := attributes[fieldName]; ok {
			hasAttribute = true
			if allowXPath {
				attValue = attributeValueRE.ReplaceAllStringFunc(a, func(a string) string {
					// strip curly braces
					seq, err := xd.data.Evaluate(a[1 : len(a)-1])
					if err != nil {
						bag.Logger.Errorf("Layout line %d: %s", layoutelt.Line, err)
						return ""
					}
					return seq.Stringvalue()
				})

			} else {
				attValue = a
			}
		} else {
			if mustexist {
				bag.Logger.Errorf("Layout line %d: attribute %s on element %s not found", layoutelt.Line, fieldName, layoutelt.Name)
				return fmt.Errorf("line %d: attribute %s on element %s not found", layoutelt.Line, fieldName, layoutelt.Name)
			}
			attValue = dflt
		}
		if hasAttribute {
			switch field.Type() {
			case intType:
				attInt, err := strconv.Atoi(attValue)
				if err != nil {
					return err
				}
				field.SetInt(int64(attInt))
			case stringType:
				field.SetString(attValue)
			case stringPtrType:
				field.Set(reflect.ValueOf(&attValue))
			case boolPtrType:
				b := attValue == "yes"
				field.Set(reflect.ValueOf(&b))
			case boolType:
				field.SetBool(attValue == "yes")
			case scaledPointsType:
				var wd bag.ScaledPoint
				if cols, err := strconv.Atoi(attValue); err == nil {
					switch fieldName {
					case "width":
						wd = xd.currentGrid.width(coord(cols))
					case "height":
						wd = xd.currentGrid.height(coord(cols))
					}
				} else {
					wd, err = bag.Sp(attValue)
					if err != nil {
						return err
					}
				}
				field.Set(reflect.ValueOf(wd).Convert(scaledPointsType))
			}

		}
	}
	return nil
}

func findAttribute(name string, element *goxml.Element, mustexist bool, allowXPath bool, dflt string, xp *xpath.Parser) (string, error) {
	var value string
	var found bool
	for _, attrib := range element.Attributes() {
		if attrib.Name == name {
			found = true
			value = attrib.Value
			break
		}
	}
	if !found {
		if mustexist {
			bag.Logger.Errorf("Layout line %d: attribute %s on element %s not found", element.Line, name, element.Name)
			return "", fmt.Errorf("line %d: attribute %s on element %s not found", element.Line, name, element.Name)
		}
		value = dflt
	}

	value = attributeValueRE.ReplaceAllStringFunc(value, func(a string) string {
		// strip curly braces
		seq, err := xp.Evaluate(a[1 : len(a)-1])
		if err != nil {
			bag.Logger.Errorf("Layout line %d: %s", element.Line, err)
			return ""
		}
		return seq.Stringvalue()
	})
	return value, nil
}

func (xd *xtsDocument) getAttributeBool(name string, element *goxml.Element, mustexist bool, allowXPath bool, dflt string) (bool, error) {
	attr, err := findAttribute(name, element, mustexist, allowXPath, dflt, xd.data)
	if attr == "yes" {
		return true, err
	}

	return false, err
}

func (xd *xtsDocument) getAttributeString(name string, element *goxml.Element, mustexist bool, allowXPath bool, dflt string) (string, error) {
	return findAttribute(name, element, mustexist, allowXPath, dflt, xd.data)
}

func (xd *xtsDocument) getAttributeInt(name string, element *goxml.Element, mustexist bool, allowXPath bool, dflt string) (int, error) {
	val, err := findAttribute(name, element, mustexist, allowXPath, dflt, xd.data)
	if err != nil {
		return 0, err
	}
	if val == "" {
		return 0, nil
	}
	return strconv.Atoi(val)
}

// getAttributeSize returns the provided width in scaled points.
func (xd *xtsDocument) getAttributeSize(name string, element *goxml.Element, mustexist bool, allowXPath bool, dflt string) (bag.ScaledPoint, error) {
	val, err := findAttribute(name, element, mustexist, allowXPath, dflt, xd.data)
	if err != nil {
		return 0, err
	}
	if val == "" {
		return 0, nil
	}
	return bag.Sp(val)
}

// getAttributeWidth returns the width which ich provided either by grid cells or a length value.
func (xd *xtsDocument) getAttributeWidth(name string, element *goxml.Element, mustexist bool, allowXPath bool, dflt string) (bag.ScaledPoint, error) {
	val, err := findAttribute(name, element, mustexist, allowXPath, dflt, xd.data)
	if err != nil {
		return 0, err
	}
	if val == "" {
		return 0, nil
	}
	if cols, err := strconv.Atoi(val); err == nil {
		return xd.currentGrid.width(coord(cols)), nil

	}
	return bag.Sp(val)
}

// getAttributeHeight returns the width which ich provided either by grid cells or a length value.
func (xd *xtsDocument) getAttributeHeight(name string, element *goxml.Element, mustexist bool, allowXPath bool, dflt string) (bag.ScaledPoint, error) {
	val, err := findAttribute(name, element, mustexist, allowXPath, dflt, xd.data)
	if err != nil {
		return 0, err
	}
	if val == "" {
		return 0, nil
	}
	if cols, err := strconv.Atoi(val); err == nil {
		return xd.currentGrid.height(coord(cols)), nil

	}
	return bag.Sp(val)
}

func getFourValues(str string) map[string]string {
	fields := strings.Fields(str)
	fourvalues := make(map[string]string)
	switch len(fields) {
	case 1:
		fourvalues["top"] = fields[0]
		fourvalues["bottom"] = fields[0]
		fourvalues["left"] = fields[0]
		fourvalues["right"] = fields[0]
	case 2:
		fourvalues["top"] = fields[0]
		fourvalues["bottom"] = fields[0]
		fourvalues["left"] = fields[1]
		fourvalues["right"] = fields[1]
	case 3:
		fourvalues["top"] = fields[0]
		fourvalues["left"] = fields[1]
		fourvalues["right"] = fields[1]
		fourvalues["bottom"] = fields[2]
	case 4:
		fourvalues["top"] = fields[0]
		fourvalues["right"] = fields[1]
		fourvalues["bottom"] = fields[2]
		fourvalues["left"] = fields[3]
	}

	return fourvalues
}

func getFourValuesSP(str string) (map[string]bag.ScaledPoint, error) {
	var err error
	fields := strings.Fields(str)
	fieldsSP := make([]bag.ScaledPoint, len(fields))
	for i, f := range fields {
		if fieldsSP[i], err = bag.Sp(f); err != nil {
			return nil, err
		}
	}

	fourvalues := make(map[string]bag.ScaledPoint)
	switch len(fields) {
	case 1:
		fourvalues["top"] = fieldsSP[0]
		fourvalues["bottom"] = fieldsSP[0]
		fourvalues["left"] = fieldsSP[0]
		fourvalues["right"] = fieldsSP[0]
	case 2:
		fourvalues["top"] = fieldsSP[0]
		fourvalues["bottom"] = fieldsSP[0]
		fourvalues["left"] = fieldsSP[1]
		fourvalues["right"] = fieldsSP[1]
	case 3:
		fourvalues["top"] = fieldsSP[0]
		fourvalues["left"] = fieldsSP[1]
		fourvalues["right"] = fieldsSP[1]
		fourvalues["bottom"] = fieldsSP[2]
	case 4:
		fourvalues["top"] = fieldsSP[0]
		fourvalues["right"] = fieldsSP[1]
		fourvalues["bottom"] = fieldsSP[2]
		fourvalues["left"] = fieldsSP[3]
	}
	return fourvalues, nil
}
