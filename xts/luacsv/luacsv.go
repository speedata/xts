package luacsv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"

	lua "github.com/speedata/go-lua"
	"golang.org/x/text/encoding/charmap"
)

var reCarriageReturn = regexp.MustCompile(`\r`)

func lerr(l *lua.State, errormessage string) int {
	l.SetTop(0)
	l.PushBoolean(false)
	l.PushString(errormessage)
	return 2
}

func decode(l *lua.State) int {
	if l.Top() < 1 {
		return lerr(l, "The first argument of decode must be the filename of the CSV.")
	}
	filename := lua.CheckString(l, 1)

	columns := []int{}
	var charset, separator string

	if l.Top() > 1 {
		if l.IsTable(-1) {
			l.Field(-1, "charset")
			if l.IsString(-1) {
				charset, _ = l.ToString(-1)
			}
			l.Pop(1)

			l.Field(-1, "separator")
			if l.IsString(-1) {
				separator, _ = l.ToString(-1)
			}
			l.Pop(1)

			l.Field(-1, "columns")
			if l.IsTable(-1) {
				length := l.RawLength(-1)
				for i := 1; i <= length; i++ {
					l.RawGetInt(-1, i)
					if n, ok := l.ToNumber(-1); ok {
						columns = append(columns, int(n))
					}
					l.Pop(1)
				}
			}
			l.Pop(1)
		}
	}

	var err error
	var rd io.Reader

	rd, err = os.Open(filename)
	if err != nil {
		return lerr(l, err.Error())
	}

	switch charset {
	case "ISO-8859-1":
		rd = charmap.ISO8859_1.NewDecoder().Reader(rd)
	}

	data, err := io.ReadAll(rd)
	if err != nil {
		return lerr(l, err.Error())
	}

	data = reCarriageReturn.ReplaceAll(data, []byte{10})
	br := bytes.NewReader(data)
	reader := csv.NewReader(br)
	if separator != "" {
		reader.Comma = rune(separator[0])
	}

	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		return lerr(l, err.Error())
	}

	l.NewTable() // rows
	for i, row := range records {
		if i == 0 && len(columns) == 0 {
			for z := 1; z <= len(row); z++ {
				columns = append(columns, z)
			}
		}
		l.NewTable() // col
		for j, entry := range columns {
			if entry-1 < 0 || entry > len(row) {
				return lerr(l, fmt.Sprintf("Column %d out of range. Must be between 1 and %d (# of columns)", entry, len(row)))
			}
			l.PushString(row[entry-1])
			l.RawSetInt(-2, j+1)
		}
		l.RawSetInt(-2, i+1)
	}

	return 1
}

// Open starts this lua module
func Open(l *lua.State) int {
	lua.NewLibrary(l, []lua.RegistryFunction{
		{Name: "decode", Function: decode},
	})
	return 1
}
