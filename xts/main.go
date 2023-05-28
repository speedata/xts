package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/pelletier/go-toml/v2"
	"github.com/speedata/boxesandglue/backend/bag"
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
		Systemfonts:  false,
		Verbose:      false,
		VariablesMap: make(map[string]any),
	}
	configfilename string = "xts.cfg"
)

// config holds global configuration that is not document dependant.
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
	Systemfonts  bool           `mapstructure:"systemfonts"`
	Verbose      bool           `mapstructure:"verbose"`
	Trace        []string       `mapstructure:"trace"`
	Variables    []string       `mapstructure:"-" toml:"-"`
	VariablesMap map[string]any `mapstructure:"-" toml:"variables"`
}

// Create a new logger instance which logs info to stdout and perhaps more to
// the protocol file.
func newZapLogger() (*zap.SugaredLogger, error) {
	protocolFile := configuration.Jobname + ".protocol"
	w, err := os.Create(protocolFile)
	if err != nil {
		return nil, err
	}
	errorPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	var protocolPriority, consolePriority zap.LevelEnablerFunc
	if configuration.Verbose {
		protocolPriority = debugPriority
	} else {
		protocolPriority = infoPriority
	}

	switch configuration.LogLevel {
	case "debug":
		consolePriority = debugPriority
	case "info":
		consolePriority = infoPriority
	case "warn":
		consolePriority = warnPriority
	case "error":
		consolePriority = errorPriority
	default:
		return nil, fmt.Errorf("could not parse the log level %q", configuration.LogLevel)
	}
	colorEncoder := zapcore.EncoderConfig{
		EncodeLevel: zapcore.LowercaseColorLevelEncoder,
		LevelKey:    "level",
		MessageKey:  "message",
	}
	simpleEncoder := zapcore.EncoderConfig{
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		LevelKey:    "level",
		MessageKey:  "message",
	}

	var consoleDebugging zapcore.WriteSyncer
	if configuration.Quiet {
		consoleDebugging = zapcore.AddSync(io.Discard)
	} else {
		consoleDebugging = zapcore.Lock(os.Stdout)
	}
	consoleEncoder := zapcore.NewConsoleEncoder(colorEncoder)

	fileEncoder := zapcore.NewConsoleEncoder(simpleEncoder)
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleDebugging, consolePriority),
		zapcore.NewCore(fileEncoder, w, protocolPriority),
	)

	logger := zap.New(core)
	return logger.Sugar(), nil
}

func listFonts() error {
	core.InitDirs(configuration.basedir)
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
		fmt.Printf("<LoadFontfile name=%q filename=%q />\n", fp.PostscriptName(), filepath.Base(fontfile))
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

	err = ioutil.WriteFile("data.xml", []byte(dataTxt), 0644)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("layout.xml", []byte(layoutTxt), 0644)
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
	op.On("--quiet", "Run XTS in quiet mode", cmdline)
	op.On("--systemfonts", "Use system fonts", cmdline)
	op.On("--trace NAMES", "Set the trace to one or more of grid, allocation", cmdline)
	op.On("--verbose", "Put more debugging information into the protocol file", cmdline)
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
			if v == jobname+"-aux.xml" {
				fmt.Printf("Remove %s\n", v)
				if err = os.Remove(v); err != nil {
					return err
				}
			}
		}
	case "doc":
		return openURL("https://doc.speedata.de/xts/")
	case "list-fonts":
		if err = listFonts(); err != nil {
			bag.Logger.Error(err)
			return err
		}
	case "new":
		if err = scaffold(op.Extra[1:]...); err != nil {
			return err
		}
		os.Exit(0)
	case "run":
		if bag.Logger, err = newZapLogger(); err != nil {
			return err
		}
		for _, cfg := range configFileRead {
			bag.Logger.Infof("Use configuration file %s", cfg)
		}
		core.InitDirs(configuration.basedir)
		core.AddDir(currentDir)
		if luafile := configuration.Filter; luafile != "" {
			if err = runLuaScript(luafile); err != nil {
				return err
			}
		}
		var layoutpath, datapath string
		var lr, dr io.ReadCloser
		if layoutpath, err = core.FindFile(configuration.Layout); err != nil {
			return err
		}
		if lr, err = os.Open(layoutpath); err != nil {
			return err
		}

		if configuration.Dummy {
			dr = io.NopCloser(strings.NewReader(`<data />`))
		} else {
			if datapath, err = core.FindFile(configuration.Data); err != nil {
				return err
			}
			if dr, err = os.Open(datapath); err != nil {
				return err
			}
		}

		if configuration.Verbose {
			data, err := os.ReadFile(layoutpath)
			if err != nil {
				return err
			}
			bag.Logger.Debugf("md5 checksum %s: %x", configuration.Layout, md5.Sum(data))
			data, err = os.ReadFile(datapath)
			if err != nil {
				return err
			}
			bag.Logger.Debugf("md5 checksum %s: %x", configuration.Data, md5.Sum(data))
		}

		xc := &core.XTSConfig{
			Layoutfile:  lr,
			Datafile:    dr,
			OutFilename: configuration.Jobname + ".pdf",
			FindFile:    core.FindFile,
			Variables:   configuration.VariablesMap,
			Mode:        configuration.Mode,
			Tracing:     configuration.Trace,
		}

		if fn := dumpOutputFileName; fn != "" {
			w, err := os.Create(fn)
			if err != nil {
				return err
			}
			xc.DumpFile = w

		}

		if err = core.RunXTS(xc); err != nil {
			return err
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
