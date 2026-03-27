package luaxml

import (
	"encoding/xml"
	"io"
	"os"

	lua "github.com/speedata/go-lua"
)

func decodeXML(l *lua.State) int {
	filename := lua.CheckString(l, 1)
	f, err := os.Open(filename)
	if err != nil {
		return lerr(l, err.Error())
	}

	defer f.Close()

	dec := xml.NewDecoder(f)

	// We track parent tables via a stack of absolute indices.
	// The root table will be left on the Lua stack at the end.
	type stackEntry struct {
		absIdx int
	}
	var parentStack []stackEntry

done:
	for {
		tok, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				break done
			}
			return lerr(l, err.Error())
		}
		switch t := tok.(type) {
		case xml.StartElement:
			l.NewTable()
			curIdx := l.AbsIndex(-1)

			l.PushString("element")
			l.SetField(curIdx, "type")
			l.PushString(t.Name.Local)
			l.SetField(curIdx, "name")

			if len(t.Attr) > 0 {
				l.NewTable()
				for _, attr := range t.Attr {
					key := attr.Name.Local
					if attr.Name.Space != "" {
						key = attr.Name.Space + ":" + attr.Name.Local
					}
					l.PushString(attr.Value)
					l.SetField(-2, key)
				}
				l.SetField(curIdx, "attribs")
			}

			if len(parentStack) > 0 {
				parentIdx := parentStack[len(parentStack)-1].absIdx
				// append to parent: parent[#parent+1] = cur
				n := l.RawLength(parentIdx)
				l.PushValue(curIdx)
				l.RawSetInt(parentIdx, n+1)
			}

			parentStack = append(parentStack, stackEntry{curIdx})

		case xml.CharData:
			if len(parentStack) > 0 {
				parentIdx := parentStack[len(parentStack)-1].absIdx
				n := l.RawLength(parentIdx)
				l.PushString(string(t.Copy()))
				l.RawSetInt(parentIdx, n+1)
			}

		case xml.EndElement:
			if len(parentStack) > 1 {
				// Pop current element from the Lua stack (parent still references it)
				// but keep it if it was appended
				parentStack = parentStack[:len(parentStack)-1]
				// The table is still on the Lua stack; remove it since
				// the parent already holds a reference via RawSetInt.
				l.Remove(-1)
			} else if len(parentStack) == 1 {
				// Root element — leave it on the stack
				parentStack = parentStack[:0]
			}
		}
	}
	// stack should have the root table on top
	l.PushBoolean(true)
	l.Insert(-2) // put true before the root table
	return 2
}
