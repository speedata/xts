package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/mitchellh/mapstructure"
	lua "github.com/speedata/go-lua"
	"github.com/speedata/xts/core"
	"github.com/speedata/xts/xts/luacsv"
	"github.com/speedata/xts/xts/luahttp"
	"github.com/speedata/xts/xts/luaxlsx"
	"github.com/speedata/xts/xts/luaxml"
)

var (
	options map[string]any
	l       *lua.State
)

func lerr(errormessage string) int {
	l.SetTop(0)
	l.PushBoolean(false)
	l.PushString(errormessage)
	return 2
}

func findFile(l *lua.State) int {
	numberArguments := l.Top()
	if numberArguments != 1 {
		return lerr("find_file requires 1 argument: the file to find")
	}
	fn := lua.CheckString(l, 1)
	abspath, err := core.FindFile(fn)
	if abspath == "" {
		if err != nil {
			l.PushNil()
			l.PushString(err.Error())
			return 2
		}
		l.PushNil()
		return 1
	}
	l.PushString(abspath)
	return 1
}

func execute(l *lua.State) int {
	lua.CheckType(l, 1, lua.TypeTable)
	var cmd string
	var arguments []string

	length := l.RawLength(1)
	for i := 1; i <= length; i++ {
		l.RawGetInt(1, i)
		s, _ := l.ToString(-1)
		l.Pop(1)
		if i == 1 {
			cmd = s
		} else {
			arguments = append(arguments, s)
		}
	}
	command := exec.Command(cmd, arguments...)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	err := command.Run()
	if err != nil {
		return lerr(err.Error())
	}
	l.PushBoolean(true)
	return 1
}

var exports = []lua.RegistryFunction{
	{Name: "find_file", Function: findFile},
	{Name: "execute", Function: execute},
}

func runtimeLoader(l *lua.State) int {
	lua.NewLibrary(l, exports)
	fillRuntimeModule(l)
	return 1
}

// set projectdir and variables table
func fillRuntimeModule(l *lua.State) {
	// variables sub-table with metatable
	l.NewTable()
	l.NewTable()
	l.PushGoFunction(indexVariables)
	l.SetField(-2, "__index")
	l.PushGoFunction(newIndexVariables)
	l.SetField(-2, "__newindex")
	l.SetMetaTable(-2)
	l.SetField(-2, "variables")

	// options sub-table with metatable
	l.NewTable()
	l.NewTable()
	l.PushGoFunction(indexOptions)
	l.SetField(-2, "__index")
	l.PushGoFunction(newIndexOptions)
	l.SetField(-2, "__newindex")
	l.SetMetaTable(-2)
	l.SetField(-2, "options")

	// log sub-table
	l.NewTable()
	lua.SetFunctions(l, []lua.RegistryFunction{
		{Name: "debug", Function: debugLog},
		{Name: "info", Function: infoLog},
		{Name: "warn", Function: warnLog},
		{Name: "error", Function: errorLog},
	}, 0)
	l.SetField(-2, "log")

	wd, _ := os.Getwd()
	l.PushString(wd)
	l.SetField(-2, "projectdir")
}

func newIndexOptions(l *lua.State) int {
	numberArguments := l.Top()
	if numberArguments < 3 {
		l.PushNil()
		return 1
	}
	// 1: tbl, 2: key, 3: value
	optionName := lua.CheckString(l, 2)
	switch l.TypeOf(3) {
	case lua.TypeNil:
		delete(options, optionName)
	case lua.TypeBoolean:
		options[optionName] = fmt.Sprint(l.ToBoolean(3))
	case lua.TypeTable:
		str := []string{}
		length := l.RawLength(3)
		for i := 1; i <= length; i++ {
			l.RawGetInt(3, i)
			s, _ := l.ToString(-1)
			l.Pop(1)
			str = append(str, s)
		}
		options[optionName] = str
	default:
		optionValue := lua.CheckString(l, 3)
		options[optionName] = optionValue
	}
	return 0
}

func indexOptions(l *lua.State) int {
	numberArguments := l.Top()
	if numberArguments < 2 {
		l.PushNil()
		return 1
	}
	// 1: tbl, 2: key
	optionName := lua.CheckString(l, 2)
	l.PushString(fmt.Sprintf("%s", options[optionName]))
	return 1
}

func newIndexVariables(l *lua.State) int {
	numberArguments := l.Top()
	if numberArguments < 3 {
		l.PushNil()
		return 1
	}
	// 1: tbl, 2: key, 3: value
	variableName := lua.CheckString(l, 2)
	value := lua.CheckString(l, 3)
	configuration.VariablesMap[variableName] = value
	return 0
}

func indexVariables(l *lua.State) int {
	numberArguments := l.Top()
	if numberArguments < 2 {
		l.PushNil()
		return 1
	}
	// 1: tbl, 2: key
	variableName := lua.CheckString(l, 2)
	l.PushString(fmt.Sprintf("%s", configuration.VariablesMap[variableName]))
	return 1
}

func debugLog(l *lua.State) int {
	slog.Debug(lua.CheckString(l, 1))
	return 0
}

func infoLog(l *lua.State) int {
	slog.Info(lua.CheckString(l, 1))
	return 0
}

func warnLog(l *lua.State) int {
	slog.Warn(lua.CheckString(l, 1))
	return 0
}

func errorLog(l *lua.State) int {
	slog.Error(lua.CheckString(l, 1))
	return 0
}

// When runtime.finalizer is set, call that function after
// the publishing run
func runFinalizerCallback() {
	l.Global("runtime")
	if l.IsNil(-1) {
		l.Pop(1)
		return
	}
	if !l.IsTable(-1) {
		l.Pop(1)
		return
	}
	l.Field(-1, "finalizer")
	if l.IsFunction(-1) {
		l.Call(0, 0)
	} else {
		l.Pop(1)
	}
	l.Pop(1) // pop runtime table
}

// preloadModule registers a module loader in package.preload[name]
func preloadModule(l *lua.State, name string, loader lua.Function) {
	l.Global("package")
	l.Field(-1, "preload")
	l.PushGoFunction(loader)
	l.SetField(-2, name)
	l.Pop(2) // pop preload and package
}

func runLuaScript(filename string) error {
	if l == nil {
		l = lua.NewState()
		lua.OpenLibraries(l)
	}

	var err error
	if err = mapstructure.Decode(configuration, &options); err != nil {
		return err
	}

	preloadModule(l, "runtime", runtimeLoader)
	preloadModule(l, "csv", luacsv.Open)
	preloadModule(l, "xml", luaxml.Open)
	preloadModule(l, "xlsx", luaxlsx.Open)
	preloadModule(l, "http", luahttp.Open)

	if err := lua.DoFile(l, filename); err != nil {
		return err
	}

	decConfig := mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           configuration,
	}
	dec, err := mapstructure.NewDecoder(&decConfig)
	if err != nil {
		return err
	}
	if err = dec.Decode(options); err != nil {
		return err
	}

	return nil
}
