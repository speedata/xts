package luaxml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"

	lua "github.com/speedata/go-lua"
)

func lerr(l *lua.State, errormessage string) int {
	l.SetTop(0)
	l.PushBoolean(false)
	l.PushString(errormessage)
	return 2
}

// encodeComment reads a comment from the table at the top of the stack
func encodeComment(l *lua.State, enc *xml.Encoder) error {
	l.Field(-1, "_value")
	if !l.IsString(-1) {
		l.Pop(1)
		return fmt.Errorf("error reading comment")
	}
	comment, _ := l.ToString(-1)
	l.Pop(1)

	c := xml.Comment([]byte(comment))
	return enc.EncodeToken(c)
}

// encodeElement reads an element from the table at the top of the stack
func encodeElement(l *lua.State, enc *xml.Encoder) error {
	l.Field(-1, "_name")
	localName, _ := l.ToString(-1)
	l.Pop(1)

	start := xml.StartElement{
		Name: xml.Name{
			Local: localName,
		},
	}

	// iterate table for attributes (string keys not starting with '_')
	l.PushNil()
	for l.Next(-2) {
		if l.IsString(-2) {
			key, _ := l.ToString(-2)
			if key[0] != '_' {
				val, _ := l.ToString(-1)
				attr := xml.Attr{
					Value: val,
					Name: xml.Name{
						Local: key,
					},
				}
				start.Attr = append(start.Attr, attr)
			}
		}
		l.Pop(1)
	}

	if err := enc.EncodeToken(start); err != nil {
		return err
	}

	// iterate table for children (integer keys)
	l.PushNil()
	for l.Next(-2) {
		if l.IsNumber(-2) {
			switch l.TypeOf(-1) {
			case lua.TypeTable:
				if err := encodeItem(l, enc); err != nil {
					l.Pop(2) // key + value
					return err
				}
			case lua.TypeString:
				s, _ := l.ToString(-1)
				if err := enc.EncodeToken(xml.CharData([]byte(s))); err != nil {
					l.Pop(2)
					return err
				}
			default:
				l.Pop(2)
				return fmt.Errorf("unknown type: %s", lua.TypeNameOf(l, -1))
			}
		}
		l.Pop(1)
	}

	return enc.EncodeToken(start.End())
}

// encodeItem encodes the table at the top of the stack
func encodeItem(l *lua.State, enc *xml.Encoder) error {
	l.Field(-1, "_type")
	typ := "element"
	if l.IsString(-1) {
		typ, _ = l.ToString(-1)
	}
	l.Pop(1)

	switch typ {
	case "element":
		return encodeElement(l, enc)
	case "comment":
		return encodeComment(l, enc)
	}
	return nil
}

func encodeTable(l *lua.State) int {
	filename := "data.xml"
	if l.Top() > 1 {
		filename = lua.CheckString(l, 2)
	}
	var b bytes.Buffer
	enc := xml.NewEncoder(&b)
	lua.CheckType(l, 1, lua.TypeTable)
	l.PushValue(1) // push table copy to top for encodeItem
	if err := encodeItem(l, enc); err != nil {
		l.Pop(1)
		return lerr(l, err.Error())
	}
	l.Pop(1)
	enc.Flush()
	if err := os.WriteFile(filename, b.Bytes(), 0644); err != nil {
		return lerr(l, err.Error())
	}
	l.SetTop(0)
	l.PushBoolean(true)
	return 1
}

// Open starts this lua module
func Open(l *lua.State) int {
	lua.NewLibrary(l, []lua.RegistryFunction{
		{"encode_table", encodeTable},
		{"decode_xml", decodeXML},
	})
	return 1
}
