package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/optionparser"
	"github.com/speedata/textlayout/fonts/truetype"
	"github.com/speedata/xts/core"
)

var (
	version       string
	configuration = &config{
		verbose:     false,
		systemfonts: false,
	}
)

// config holds global configuration that is not document dependant.
type config struct {
	verbose     bool
	systemfonts bool
	basedir     string
}

func newZapLogger(verbose bool) (*zap.SugaredLogger, error) {
	cfg := zap.Config{
		Encoding:    "console",
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			EncodeLevel: zapcore.LowercaseColorLevelEncoder,
			LevelKey:    "level",
			MessageKey:  "message",
		},
	}
	if verbose {
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	} else {
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
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

func dothings() error {
	pathToXTS, err := os.Executable()
	if err != nil {
		return err
	}
	configuration.basedir = filepath.Join(filepath.Dir(pathToXTS), "..")
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	op := optionparser.NewOptionParser()
	op.On("--systemfonts", "Use system fonts", &configuration.systemfonts)
	op.On("--verbose", "Print more debugging information", &configuration.verbose)
	op.Command("list-fonts", "List installed fonts")
	op.Command("run", "Load layout and data files and create PDF")
	op.Command("version", "Print version information")
	err = op.Parse()
	if err != nil {
		op.Help()
		return err
	}
	if bag.Logger, err = newZapLogger(configuration.verbose); err != nil {
		return err
	}
	if configuration.systemfonts {
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
	case "list-fonts":
		if err = listFonts(); err != nil {
			bag.Logger.Error(err)
			return err
		}
	case "run":
		core.InitDirs(configuration.basedir)
		core.AddDir(currentDir)
		var layoutpath, datapath string
		var lr, dr io.ReadCloser

		if layoutpath, err = core.FindFile("layout.xml"); err != nil {
			return err
		}
		if lr, err = os.Open(layoutpath); err != nil {
			return err
		}
		if datapath, err = core.FindFile("data.xml"); err != nil {
			return err
		}
		if dr, err = os.Open(datapath); err != nil {
			return err
		}

		xc := &core.XTSConfig{
			Layoutfile:  lr,
			Datafile:    dr,
			OutFilename: "publisher.pdf",
			FindFile:    core.FindFile,
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
