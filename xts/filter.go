package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

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

func validateRelaxNG(l *lua.State) int {
	xmlfile := lua.CheckString(l, 1)
	rngfile := lua.CheckString(l, 2)

	cmd := exec.Command("java", "-jar", filepath.Join(configuration.libdir, "jing.jar"), rngfile, xmlfile)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return lerr(err.Error())
	}
	var b bytes.Buffer

	err = cmd.Start()
	if err != nil {
		return lerr(err.Error())
	}

	go io.Copy(&b, stdoutPipe)
	err = cmd.Wait()
	if err != nil {
		return lerr(b.String())
	}

	l.PushBoolean(true)
	return 1
}

func runSaxon(l *lua.State) int {
	numberArguments := l.Top()
	var command []string
	command = []string{"-jar", filepath.Join(configuration.libdir, "saxon-he-12.9.jar")}
	if numberArguments == 1 {
		if l.IsTable(-1) {
			m := map[string]string{
				"initialtemplate": "-it:%s",
				"source":          "-s:%s",
				"stylesheet":      "-xsl:%s",
				"out":             "-o:%s",
			}
			for k, val := range m {
				l.Field(-1, k)
				if l.IsString(-1) {
					str, _ := l.ToString(-1)
					command = append(command, fmt.Sprintf(val, str))
				}
				l.Pop(1)
			}
			// parameters at the end
			l.Field(-1, "params")
			if l.IsString(-1) {
				str, _ := l.ToString(-1)
				command = append(command, str)
			} else if l.IsTable(-1) {
				l.PushNil()
				for l.Next(-2) {
					key, _ := l.ToString(-2)
					value, _ := l.ToString(-1)
					command = append(command, fmt.Sprintf("%s=%s", key, value))
					l.Pop(1)
				}
			}
			l.Pop(1)
		} else {
			return lerr("The single argument must be a table (run_saxon)")
		}
	} else if numberArguments < 3 {
		return lerr("command requires 3 or 4 arguments")
	} else {
		xsl := lua.CheckString(l, 1)
		src := lua.CheckString(l, 2)
		out := lua.CheckString(l, 3)

		command = append(command, fmt.Sprintf("-xsl:%s", xsl), fmt.Sprintf("-s:%s", src), fmt.Sprintf("-o:%s", out))

		if numberArguments > 3 {
			command = append(command, lua.CheckString(l, 4))
		}
	}
	if configuration.Verbose {
		fmt.Println(command)
	}

	cmd := exec.Command("java", command...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	cmd.Start()

	if err := cmd.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			l.PushBoolean(false)
		}
	} else {
		l.PushBoolean(true)
	}
	return 1
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
	{Name: "validate_relaxng", Function: validateRelaxNG},
	{Name: "run_saxon", Function: runSaxon},
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
