package luaxml

import (
	"fmt"
	"os"

	lua "github.com/speedata/go-lua"
	"github.com/speedata/goxml"
	"github.com/speedata/goxpath"
	"github.com/speedata/goxslt"
)

func runXSLT(l *lua.State) int {
	numberArguments := l.Top()
	var stylesheetPath, sourcePath, outPath, initialTemplate string
	params := map[string]goxpath.Sequence{}

	if numberArguments == 1 {
		if !l.IsTable(-1) {
			return lerr(l, "The single argument must be a table (run_xslt)")
		}
		for _, field := range []struct {
			name string
			dest *string
		}{
			{"stylesheet", &stylesheetPath},
			{"source", &sourcePath},
			{"out", &outPath},
			{"initialtemplate", &initialTemplate},
		} {
			l.Field(-1, field.name)
			if l.IsString(-1) {
				*field.dest, _ = l.ToString(-1)
			}
			l.Pop(1)
		}
		l.Field(-1, "params")
		if l.IsString(-1) {
			str, _ := l.ToString(-1)
			params[""] = goxpath.Sequence{str}
		} else if l.IsTable(-1) {
			l.PushNil()
			for l.Next(-2) {
				key, _ := l.ToString(-2)
				value, _ := l.ToString(-1)
				params[key] = goxpath.Sequence{value}
				l.Pop(1)
			}
		}
		l.Pop(1)
	} else if numberArguments < 3 {
		return lerr(l, "command requires 3 or 4 arguments")
	} else {
		stylesheetPath = lua.CheckString(l, 1)
		sourcePath = lua.CheckString(l, 2)
		outPath = lua.CheckString(l, 3)
	}

	if stylesheetPath == "" {
		return lerr(l, "stylesheet is required (run_xslt)")
	}

	ss, err := goxslt.CompileFile(stylesheetPath)
	if err != nil {
		return lerr(l, fmt.Sprintf("XSLT compile error: %s", err))
	}

	opts := goxslt.TransformOptions{
		Parameters: params,
	}
	if initialTemplate != "" {
		opts.InitialTemplate = initialTemplate
	}

	var result *goxslt.TransformResult
	if sourcePath != "" {
		sourceFile, err := os.Open(sourcePath)
		if err != nil {
			return lerr(l, fmt.Sprintf("cannot open source: %s", err))
		}
		defer sourceFile.Close()
		sourceDoc, err := goxml.Parse(sourceFile)
		if err != nil {
			return lerr(l, fmt.Sprintf("XML parse error: %s", err))
		}
		result, err = goxslt.TransformWithOptions(ss, sourceDoc, opts)
	} else {
		result, err = goxslt.TransformWithOptions(ss, nil, opts)
	}
	if err != nil {
		return lerr(l, fmt.Sprintf("XSLT transform error: %s", err))
	}

	output := goxslt.SerializeWithOutput(result.Document, result.Output)

	if outPath != "" {
		if err := os.WriteFile(outPath, []byte(output), 0644); err != nil {
			return lerr(l, fmt.Sprintf("cannot write output: %s", err))
		}
	}

	for href, doc := range result.SecondaryDocuments {
		secOutput := goxslt.SerializeWithOutput(doc, result.Output)
		if err := os.WriteFile(href, []byte(secOutput), 0644); err != nil {
			return lerr(l, fmt.Sprintf("cannot write secondary document %s: %s", href, err))
		}
	}

	l.PushBoolean(true)
	return 1
}
