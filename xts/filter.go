package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/speedata/xts/core"
	"github.com/speedata/xts/xts/luacsv"
	"github.com/speedata/xts/xts/luaxlsx"
	"github.com/speedata/xts/xts/luaxml"

	"github.com/cjoudrey/gluahttp"
	lua "github.com/yuin/gopher-lua"
)

var (
	options map[string]any
	l       *lua.LState
)

func lerr(errormessage string) int {
	l.SetTop(0)
	l.Push(lua.LFalse)
	l.Push(lua.LString(errormessage))
	return 2
}

func validateRelaxNG(l *lua.LState) int {
	xmlfile := l.CheckString(1)
	rngfile := l.CheckString(2)

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

	l.Push(lua.LTrue)
	return 1
}

func runSaxon(l *lua.LState) int {
	numberArguments := l.GetTop()
	var command []string
	command = []string{"-jar", filepath.Join(configuration.libdir, "saxon-he-12.1.jar")}
	if numberArguments == 1 {
		// hopefully a table
		lv := l.Get(-1)
		if tbl, ok := lv.(*lua.LTable); ok {
			m := map[string]string{
				"initialtemplate": "-it:%s",
				"source":          "-s:%s",
				"stylesheet":      "-xsl:%s",
				"out":             "-o:%s",
			}
			for k, val := range m {
				if str := tbl.RawGetString(k); str.Type() == lua.LTString {
					command = append(command, fmt.Sprintf(val, str.String()))
				}
			}
			// parameters at the end
			if str := tbl.RawGetString("params"); str.Type() == lua.LTString {
				command = append(command, str.String())
			} else if tbl := tbl.RawGetString("params"); tbl.Type() == lua.LTTable {
				if paramtbl, ok := tbl.(*lua.LTable); ok {
					paramtbl.ForEach(func(key lua.LValue, value lua.LValue) {
						command = append(command, fmt.Sprintf("%s=%s", key.String(), value.String()))
					})
				}
			}

		} else {
			return lerr("The single argument must be a table (run_saxon)")
		}
	} else if numberArguments < 3 {
		return lerr("command requires 3 or 4 arguments")
	} else {
		xsl := l.CheckString(1)
		src := l.CheckString(2)
		out := l.CheckString(3)

		command = append(command, fmt.Sprintf("-xsl:%s", xsl), fmt.Sprintf("-s:%s", src), fmt.Sprintf("-o:%s", out))

		// fourth argument param is optional
		if numberArguments > 3 {
			command = append(command, l.CheckString(4))
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
			l.Push(lua.LFalse)
		}
	} else {
		l.Push(lua.LTrue)
	}
	return 1
}

func findFile(l *lua.LState) int {
	numberArguments := l.GetTop()
	if numberArguments != 1 {
		return lerr("find_file requires 1 argument: the file to find")
	}
	fn := l.CheckString(1)
	abspath, err := core.FindFile(fn)
	if abspath == "" {
		if err != nil {
			l.Push(lua.LNil)
			l.Push(lua.LString(err.Error()))
			return 2
		}
		l.Push(lua.LNil)
		return 1
	}
	l.Push(lua.LString(abspath))
	return 1
}
func execute(l *lua.LState) int {
	cmdline := l.CheckTable(1)
	var cmd string
	var arguments []string

	for i := 1; i <= cmdline.Len(); i++ {
		val := cmdline.RawGetInt(i)
		if i == 1 {
			cmd = val.String()
		} else {
			arguments = append(arguments, val.String())
		}
	}
	command := exec.Command(cmd, arguments...)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	err := command.Run()
	if err != nil {
		return lerr(err.Error())
	}
	l.Push(lua.LTrue)
	return 1
}

var exports = map[string]lua.LGFunction{
	"validate_relaxng": validateRelaxNG,
	"run_saxon":        runSaxon,
	"find_file":        findFile,
	"execute":          execute,
}

func runtimeLoader(l *lua.LState) int {
	mod := l.SetFuncs(l.NewTable(), exports)
	fillRuntimeModule(mod)
	l.Push(mod)
	return 1

}

// set projectdir and variables table
func fillRuntimeModule(mod lua.LValue) {
	lvars := l.NewTable()
	// for k, v := range variables {
	// 	lvars.RawSetString(k, lua.LString(v))
	// }
	l.SetField(mod, "variables", lvars)
	l.SetField(mod, "options", getOptionsTable((l)))
	wd, _ := os.Getwd()
	l.SetField(mod, "projectdir", lua.LString(wd))
}

func getOptionsTable(l *lua.LState) *lua.LTable {
	options := l.NewTable()
	mt := l.NewTable()
	l.SetField(mt, "__index", l.NewFunction(indexOptions))
	l.SetField(mt, "__newindex", l.NewFunction(newIndexOptions))
	l.SetMetatable(options, mt)
	return options
}

// Set string
func newIndexOptions(l *lua.LState) int {
	numberArguments := l.GetTop()
	if numberArguments < 3 {
		l.Push(lua.LNil)
		return 1
	}
	// 1: tbl
	// 2: key
	// 3: value
	optionName := l.CheckString(2)
	switch l.Get(3).Type() {
	case lua.LTNil:
		delete(options, optionName)
	case lua.LTBool:
		options[optionName] = fmt.Sprint(l.CheckBool(3))
	case lua.LTTable:
		ltbl := l.CheckTable(3)
		str := []string{}
		for i := 1; i <= ltbl.Len(); i++ {
			str = append(str, ltbl.RawGetInt(i).String())
		}
		options[optionName] = str
	default:
		optionValue := l.CheckString(3)
		options[optionName] = optionValue
	}
	return 0
}

func indexOptions(l *lua.LState) int {
	numberArguments := l.GetTop()
	if numberArguments < 2 {
		l.Push(lua.LNil)
		return 1
	}
	// 1: tbl
	// 2: key
	optionName := l.CheckString(2)
	l.Push(lua.LString(fmt.Sprintf("%s", options[optionName])))
	return 1
}

// When runtime.finalizer is set, call that function after
// the publishing run
func runFinalizerCallback() {
	val := l.GetGlobal("runtime")
	if val == nil {
		return
	}

	tbl, ok := val.(*lua.LTable)
	if !ok {
		return
	}
	fun := tbl.RawGetString("finalizer")
	if fn, ok := fun.(*lua.LFunction); ok {
		l.Push(fn)
		l.Call(0, 0)
	}
}

func runLuaScript(filename string) error {
	if l == nil {
		l = lua.NewState()
	}

	var err error
	if err = mapstructure.Decode(configuration, &options); err != nil {
		return err
	}

	l.PreloadModule("runtime", runtimeLoader)
	l.PreloadModule("csv", luacsv.Open)
	l.PreloadModule("xml", luaxml.Open)
	l.PreloadModule("xlsx", luaxlsx.Open)
	l.PreloadModule("http", gluahttp.NewHttpModule(&http.Client{}).Loader)

	if err := l.DoFile(filename); err != nil {
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
