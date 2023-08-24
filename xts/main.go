package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/speedata/optionparser"
	"github.com/speedata/textlayout/fonts/truetype"
	"github.com/speedata/xts/core"
)

var (
	configuration = &config{
		Data:         "data.xml",
		Dummy:        false,
		Jobname:      "xts",
		Layout:       "layout.xml",
		LogLevel:     "info",
		Runs:         1,
		SuppressInfo: false,
		Systemfonts:  false,
		Verbose:      false,
		VariablesMap: make(map[string]any),
	}
	configfilename string = "xts.cfg"
)

// config holds global configuration that is not document dependant. The
// mapstructure are for the Lua filter to map between these settings and the Lua
// values.
type config struct {
	basedir      string
	libdir       string
	Data         string         `mapstructure:"data"`
	Dummy        bool           `mapstructure:"dummy"`
	Filter       string         `mapstructure:"filter"`
	Jobname      string         `mapstructure:"jobname"`
	Layout       string         `mapstructure:"layout"`
	LogLevel     string         `mapstructure:"loglevel"`
	Mode         []string       `mapstructure:"mode"`
	Quiet        bool           `mapstructure:"quiet"`
	Runs         int            `mapstructure:"runs"`
	Systemfonts  bool           `mapstructure:"systemfonts"`
	SuppressInfo bool           `mapstructure:"suppressinfo"`
	Verbose      bool           `mapstructure:"verbose"`
	Trace        []string       `mapstructure:"trace"`
	Variables    []string       `mapstructure:"-" toml:"-"`
	VariablesMap map[string]any `mapstructure:"-" toml:"variables"`
}

func listFonts() error {
	ff := core.FindFontFiles()
	ret := make([]string, len(ff))
	for i, fontfile := range ff {
		ret[i] = filepath.Base(fontfile)
	}
	sort.Slice(ff, func(i, j int) bool {
		return filepath.Base(strings.ToLower(ff[i])) < filepath.Base(strings.ToLower(ff[j]))
	})
	for _, fontfile := range ff {
		f, err := os.Open(fontfile)
		if err != nil {
			return err
		}

		fp, err := truetype.Parse(f)
		if err != nil {
			return nil
		}
		l := strings.ToLower(fontfile)
		weight := "normal"
		style := "normal"
		switch {
		case strings.Contains(l, "regular"):
			weight = "normal"
		case strings.Contains(l, "bolditalic"):
			style = "italic"
			weight = "bold"
		case strings.Contains(l, "italic"):
			style = "italic"
		case strings.Contains(l, "bold"):
			weight = "bold"
		}
		fmt.Printf("@font-face { font-family: %q; src: url(%q);", fp.Names.SelectEntry(truetype.NameFontFamily), filepath.Base(fontfile))
		if weight != "normal" {
			fmt.Printf(" font-weight: %s; ", weight)
		}
		if style != "normal" {
			fmt.Printf(" font-style: %s;", style)
		}
		fmt.Println("}")
		if err = f.Close(); err != nil {
			return err
		}
	}
	return nil
}
func scaffold(extra ...string) error {
	var err error
	fmt.Print("Creating layout.xml and data.xml in ")
	if len(extra) > 0 {
		dir := extra[0]
		fmt.Println("a new directory", dir)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
		err = os.Chdir(dir)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("current directory")
	}

	// Let's not overwrite existing files
	_, err = os.Stat("data.xml")
	if err == nil {
		return fmt.Errorf("data.xml already exists")
	}
	_, err = os.Stat("layout.xml")
	if err == nil {
		return fmt.Errorf("layout.xml already exists")
	}

	dataTxt := `<data>Hello, world!</data>
`
	layoutTxt := `<Layout xmlns="urn:speedata.de/2021/xts/en"
    xmlns:sd="urn:speedata.de/2021/xtsfunctions/en">
    <Record element="data">
        <PlaceObject>
            <Textblock>
                <Paragraph>
                    <Value select="."/>
                </Paragraph>
            </Textblock>
        </PlaceObject>
    </Record>
</Layout>
`

	err = os.WriteFile("data.xml", []byte(dataTxt), 0644)
	if err != nil {
		return err
	}

	err = os.WriteFile("layout.xml", []byte(layoutTxt), 0644)
	if err != nil {
		return err
	}

	return nil
}

func openURL(url string) error {
	cmd := []string{}
	switch runtime.GOOS {
	case "darwin":
		cmd = append(cmd, "open", "-u")
	case "linux":
		cmd = append(cmd, "xdg-open")
	case "windows":
		cmd = append(cmd, "start")
	default:
		fmt.Printf("Open your browser at %s\n", url)
		return nil
	}
	cmd = append(cmd, url)
	ecmd := exec.Command(cmd[0], cmd[1:]...)
	return ecmd.Run()
}

func dothings() error {
	pathToXTS, err := os.Executable()
	if err != nil {
		return err
	}
	configuration.basedir = filepath.Join(filepath.Dir(pathToXTS), "..")
	configuration.libdir = filepath.Join(configuration.basedir, "lib")
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	var dumpOutputFileName string
	cmdline := make(map[string]string)

	op := optionparser.NewOptionParser()
	op.On("-c NAME", "--config", "Read the config file with the given NAME. Default: 'xts.cfg'", &configfilename)
	op.On("--data NAME", "Name of the data file. Defaults to 'data.xml'", cmdline)
	op.On("--dummy", "Don't read a data file, use '<data />' as input", cmdline)
	op.On("--dumpoutput FILENAME", "Complete XML dump of generated PDF file", &dumpOutputFileName)
	op.On("--filter NAME", "Run Lua process before the publishing run", cmdline)
	op.On("--jobname NAME", "The name of the resulting PDF file (without extension), default is 'xts'", cmdline)
	op.On("--layout NAME", "Name of the layout file. Defaults to 'layout.xml'", cmdline)
	op.On("--loglevel LVL", "Set the log level for the console to one of debug, info, warn, error", cmdline)
	op.On("--mode NAME", "Set mode. Multiple modes given in a comma separated list.", cmdline)
	op.On("--quiet", "Run XTS in quiet mode (no output on STDOUT)", cmdline)
	op.On("--runs N", "Run XTS N times", cmdline)
	op.On("--suppressinfo", "Create a reproducible document", cmdline)
	op.On("--systemfonts", "Use system fonts", cmdline)
	op.On("--trace NAMES", "Set the trace to one or more of grid, allocation", cmdline)
	op.On("--verbose", "Show log output in the terminal window (STDOUT)", cmdline)
	op.On("-v", "--var=VALUE", "Set a variable for the publishing run", cmdline)
	op.Command("list-fonts", "List installed fonts")
	op.Command("clean", "Remove auxiliary and protocol files")
	op.Command("doc", "Open the documentation (web page)")
	op.Command("new", "Create simple layout and data file to start. Provide optional directory.")
	op.Command("run", "Load layout and data files and create PDF (default)")
	op.Command("version", "Print version information")
	err = op.Parse()
	if err != nil {
		op.Help()
		fmt.Println()
		return err
	}
	var configFileRead []string
	if data, err := os.ReadFile(configfilename); err == nil {
		if err = toml.Unmarshal(data, configuration); err != nil {
			switch t := err.(type) {
			case *toml.DecodeError:
				fmt.Println(t.String())
			default:
				return err
			}
			return err
		}
		configFileRead = append(configFileRead, configfilename)
	}
	for k, v := range cmdline {
		switch k {
		case "data":
			configuration.Data = v
		case "dummy":
			configuration.Dummy = (v == "true")
		case "filter":
			configuration.Filter = v
		case "jobname":
			configuration.Jobname = v
		case "layout":
			configuration.Layout = v
		case "loglevel":
			configuration.LogLevel = v
		case "mode":
			configuration.Mode = strings.Split(v, ",")
		case "quiet":
			configuration.Quiet = (v == "true")
		case "suppressinfo":
			configuration.SuppressInfo = (v == "true")
		case "systemfonts":
			configuration.Systemfonts = (v == "true")
		case "trace":
			configuration.Trace = strings.Split(v, ",")
		case "verbose":
			configuration.Verbose = (v == "true")
		case "var":
			for _, assignment := range strings.Split(v, ",") {
				v1 := strings.Split(assignment, "=")
				if len(v1) == 2 {
					configuration.VariablesMap[v1[0]] = v1[1]
				}
			}
		case "runs":
			if configuration.Runs, err = strconv.Atoi(v); err != nil {
				return err
			}
		default:
			return fmt.Errorf("could not handle configuration %s", k)
		}
	}

	for _, vars := range configuration.Variables {
		kv := strings.Split(vars, "=")
		if len(kv) == 2 {
			configuration.VariablesMap[kv[0]] = kv[1]
		}
	}

	switch configuration.LogLevel {
	case "trace":
		loglevel.Set(slog.Level(-8))
	case "debug":
		loglevel.Set(slog.LevelDebug)
	case "info":
		loglevel.Set(slog.LevelInfo)
	case "notice":
		loglevel.Set(core.LevelNotice)
	case "warn":
		loglevel.Set(slog.LevelWarn)
	case "error":
		loglevel.Set(slog.LevelError)
	}

	if configuration.Quiet {
		os.Stdout.Close()
	}
	if configuration.Systemfonts {
		fontfolders, err := core.FontFolder()
		if err != nil {
			return err
		}
		a := strings.Split(fontfolders, ":")
		for _, d := range a {
			core.AddDir(d)
		}
	}

	cmd := "run"
	if len(op.Extra) > 0 {
		cmd = op.Extra[0]
	}
	protocolFilename := configuration.Jobname + "-protocol.xml"
	if err != nil {
		return err
	}
	switch cmd {
	case "clean":
		jobname := configuration.Jobname
		files, err := filepath.Glob(jobname + "*")
		if err != nil {
			return err
		}
		for _, v := range files {
			switch filepath.Ext(v) {
			case ".protocol":
				fmt.Printf("Remove %s\n", v)
				if err = os.Remove(v); err != nil {
					return err
				}
			}
			if strings.HasPrefix(v, jobname+"-") && strings.HasSuffix(v, ".xml") {
				fmt.Printf("Remove %s\n", v)
				if err = os.Remove(v); err != nil {
					return err
				}
			}
		}
	case "doc":
		return openURL("https://doc.speedata.de/xts/")
	case "list-fonts":
		setupLog(protocolFilename)
		defer teardownLog()

		if err = core.InitDirs(configuration.basedir); err != nil {
			return err
		}
		if err = core.AddDir(currentDir); err != nil {
			return err
		}

		if err = listFonts(); err != nil {
			slog.Error(err.Error())
			return err
		}
	case "new":
		if err = scaffold(op.Extra[1:]...); err != nil {
			return err
		}
		os.Exit(0)
	case "run":
		starttime := time.Now()
		setupLog(protocolFilename)
		defer teardownLog()

		for _, cfg := range configFileRead {
			slog.Info(fmt.Sprintf("Use configuration file %s", cfg))
		}
		core.InitDirs(configuration.basedir)
		core.AddDir(currentDir)
		if luafile := configuration.Filter; luafile != "" {
			if err = runLuaScript(luafile); err != nil {
				return err
			}
		}

		var layoutpath, datapath string
		var lr, dr io.ReadSeeker
		if layoutpath, err = core.FindFile(configuration.Layout); err != nil {
			return err
		}
		if lr, err = os.Open(layoutpath); err != nil {
			return err
		}

		if configuration.Dummy {
			dr = strings.NewReader(`<data />`)
		} else {
			if datapath, err = core.FindFile(configuration.Data); err != nil {
				slog.Error(err.Error())
				return err
			}
			if dr, err = os.Open(datapath); err != nil {
				slog.Error(err.Error())
				return err
			}
		}

		slog.Debug("checksum", "filename", configuration.Layout, "md5", md5calc(configuration.Layout))
		slog.Debug("checksum", "filename", configuration.Data, "md5", md5calc(configuration.Data))

		for i := 0; i < int(configuration.Runs); i++ {
			if cr := configuration.Runs; cr > 1 {
				slog.Info(fmt.Sprintf("Run %d of %d", i+1, cr))
			}
			lr.Seek(0, io.SeekStart)
			dr.Seek(0, io.SeekStart)
			xc := &core.XTSConfig{
				Datafile:     dr,
				FindFile:     core.FindFile,
				Layoutfile:   lr,
				Mode:         configuration.Mode,
				OutFilename:  configuration.Jobname + ".pdf",
				Jobname:      configuration.Jobname,
				SuppressInfo: configuration.SuppressInfo,
				Tracing:      configuration.Trace,
				Variables:    configuration.VariablesMap,
			}

			if fn := dumpOutputFileName; fn != "" {
				w, err := os.Create(fn)
				if err != nil {
					return err
				}
				xc.DumpFile = w

			}
			if err = core.RunXTS(xc); err != nil {
				goto finished
			}
		}
		if lrc, ok := lr.(io.ReadCloser); ok {
			if err = lrc.Close(); err != nil {
				return err
			}
		}
		if drc, ok := dr.(io.ReadCloser); ok {
			if err = drc.Close(); err != nil {
				return err
			}
		}
	finished:
		dur := time.Now().Sub(starttime)
		slog.Info(fmt.Sprintf("Finished in %s", dur))
		fmt.Printf("Finished with %d error(s) and %d warning(s) in %s.\nOutput written to %s\n  and protocol file to %s.\n", errCount, warnCount, dur, configuration.Jobname+".pdf", protocolFilename)
		if errCount > 0 {
			return core.TypesettingError{Logged: true}
		}
	case "version":
		fmt.Println("xts version", core.Version)
	}
	return nil
}

func main() {
	if err := dothings(); err != nil {
		if terr, ok := err.(core.TypesettingError); ok {
			if !terr.Logged {
				fmt.Println("Error:", err)
			}
		} else {
			fmt.Println("Error:", err)
		}
		os.Exit(1)
	}
}
