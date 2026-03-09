package luaxlsx

import (
	lua "github.com/speedata/go-lua"
	"github.com/speedata/goxlsx"
)

const luaSpreadsheetTypeName = "spreadsheet"
const luaWorksheetTypeName = "worksheet"

func lerr(l *lua.State, errormessage string) int {
	l.SetTop(0)
	l.PushBoolean(false)
	l.PushString(errormessage)
	return 2
}

// ----------------------- spreadsheet

func indexSpreadSheet(l *lua.State) int {
	sh := checkSpreadsheet(l)
	n, _ := l.ToInteger(-1)
	ws, err := sh.GetWorksheet(n - 1)
	if err != nil {
		lua.Errorf(l, "%s", err.Error())
	}

	if lua.NewMetaTable(l, luaWorksheetTypeName) {
		l.PushGoFunction(callWorksheet)
		l.SetField(-2, "__call")
		l.PushGoFunction(indexWorksheet)
		l.SetField(-2, "__index")
	}

	l.PushUserData(ws)
	lua.SetMetaTableNamed(l, luaWorksheetTypeName)

	return 1
}

func lenSpreadSheet(l *lua.State) int {
	sh := checkSpreadsheet(l)
	l.PushInteger(sh.NumWorksheets())
	return 1
}

func checkSpreadsheet(l *lua.State) *goxlsx.Spreadsheet {
	ud := lua.CheckUserData(l, 1, luaSpreadsheetTypeName)
	if v, ok := ud.(*goxlsx.Spreadsheet); ok {
		return v
	}
	lua.ArgumentError(l, 1, "spreadsheet expected")
	return nil
}

// ----------------------- worksheet

func checkWorksheet(l *lua.State) *goxlsx.Worksheet {
	ud := lua.CheckUserData(l, 1, luaWorksheetTypeName)
	if v, ok := ud.(*goxlsx.Worksheet); ok {
		return v
	}
	lua.ArgumentError(l, 1, "worksheet expected")
	return nil
}

func indexWorksheet(l *lua.State) int {
	ws := checkWorksheet(l)
	arg, _ := l.ToString(2)
	switch arg {
	case "minrow":
		l.PushInteger(ws.MinRow)
		return 1
	case "maxrow":
		l.PushInteger(ws.MaxRow)
		return 1
	case "mincol":
		l.PushInteger(ws.MinColumn)
		return 1
	case "maxcol":
		l.PushInteger(ws.MaxColumn)
		return 1
	case "name":
		l.PushString(ws.Name)
		return 1
	}
	return 0
}

func callWorksheet(l *lua.State) int {
	ws := checkWorksheet(l)
	y, _ := l.ToInteger(-1)
	x, _ := l.ToInteger(-2)
	str := ws.Cell(x, y)
	l.PushString(str)
	return 1
}

func stringToDate(l *lua.State) int {
	n := lua.CheckString(l, 1)
	t := goxlsx.DateFromString(n)
	l.NewTable()
	l.PushInteger(t.Day())
	l.SetField(-2, "day")
	l.PushInteger(int(t.Month()))
	l.SetField(-2, "month")
	l.PushInteger(t.Year())
	l.SetField(-2, "year")
	l.PushInteger(t.Hour())
	l.SetField(-2, "hour")
	l.PushInteger(t.Minute())
	l.SetField(-2, "minute")
	l.PushInteger(t.Second())
	l.SetField(-2, "second")
	return 1
}

func openfile(l *lua.State) int {
	if l.Top() < 1 {
		return lerr(l, "The first argument of open must be the filename of the Excel file.")
	}
	filename := lua.CheckString(l, 1)

	sh, err := goxlsx.OpenFile(filename)
	if err != nil {
		return lerr(l, err.Error())
	}

	if lua.NewMetaTable(l, luaSpreadsheetTypeName) {
		l.PushGoFunction(indexSpreadSheet)
		l.SetField(-2, "__index")
		l.PushGoFunction(lenSpreadSheet)
		l.SetField(-2, "__len")
	}

	l.PushUserData(sh)
	lua.SetMetaTableNamed(l, luaSpreadsheetTypeName)

	return 1
}

// Open sets up the XLSX Lua module.
func Open(l *lua.State) int {
	lua.NewLibrary(l, []lua.RegistryFunction{
		{"open", openfile},
		{"string_to_date", stringToDate},
	})
	return 1
}
