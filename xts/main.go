package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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
		Verbose:     false,
		Systemfonts: false,
		Dummy:       false,
		Jobname:     "publisher",
		Data:        "data.xml",
		Layout:      "layout.xml",
	}
	configfilename string = "publisher.cfg"
)

// config holds global configuration that is not document dependant.
type config struct {
	Verbose     bool
	Systemfonts bool
	Dummy       bool
	Basedir     string
	Jobname     string
	Data        string
	Layout      string
}

// Create a new logger instance which logs info to stdout and perhaps more to
// the protocol file.
func newZapLogger() (*zap.SugaredLogger, error) {
	protocolFile := configuration.Jobname + ".protocol"
	w, err := os.Create(protocolFile)
	if err != nil {
		return nil, err
	}

	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	var protocolPriority zap.LevelEnablerFunc
	if configuration.Verbose {
		protocolPriority = debugPriority
	} else {
		protocolPriority = infoPriority
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

	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleEncoder := zapcore.NewConsoleEncoder(colorEncoder)

	fileEncoder := zapcore.NewConsoleEncoder(simpleEncoder)
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleDebugging, infoPriority),
		zapcore.NewCore(fileEncoder, w, protocolPriority),
	)

	logger := zap.New(core)
	return logger.Sugar(), nil
}

func listFonts() error {
	core.InitDirs(configuration.Basedir)
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

func dothings() error {
	pathToXTS, err := os.Executable()
	if err != nil {
		return err
	}
	configuration.Basedir = filepath.Join(filepath.Dir(pathToXTS), "..")
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	var dumpOutputFileName string

	op := optionparser.NewOptionParser()
	op.On("-c NAME", "--config", "Read the config file with the given NAME. Default: 'publisher.cfg'", &configfilename)
	op.On("--jobname NAME", "The name of the resulting PDF file (without extension), default is 'publisher'", &configuration.Jobname)
	op.On("--systemfonts", "Use system fonts", &configuration.Systemfonts)
	op.On("--verbose", "Put more debugging information into the protocol file", &configuration.Verbose)
	op.On("--dummy", "Don't read a data file, use '<data />' as input", &configuration.Dummy)
	op.On("--dumpoutput FILENAME", "Complete XML dump of generated PDF file", &dumpOutputFileName)
	op.Command("list-fonts", "List installed fonts")
	op.Command("clean", "Remove auxiliary and protocol files")
	op.Command("new", "Create simple layout and data file to start. Provide optional directory.")
	op.Command("run", "Load layout and data files and create PDF (default)")
	op.Command("version", "Print version information")
	err = op.Parse()
	if err != nil {
		op.Help()
		return err
	}
	if data, err := os.ReadFile(configfilename); err == nil {
		if err = toml.Unmarshal(data, configuration); err != nil {
			fmt.Println(err.(*toml.DecodeError).String())
			return err
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

		core.InitDirs(configuration.Basedir)
		core.AddDir(currentDir)
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
		}

		if fn := dumpOutputFileName; fn != "" {
			fmt.Println("****dumping ", fn)
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
		bag.Logger.Error(err)
		os.Exit(1)
	}
}
