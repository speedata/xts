package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/optionparser"
	"github.com/speedata/xts/core"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	filelist = make(map[string]string)
	version  string
)

func initDirs() error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, dir := range []string{"fonts", "hyphenationpatterns", "img"} {
		dir = filepath.Join(root, dir)
		bag.Logger.Debugf("Add directory %q to recursive file list", dir)
		filepath.WalkDir(dir, dirWalker)
	}
	return nil
}

func dirWalker(path string, d fs.DirEntry, err error) error {
	if d.Type().IsRegular() {
		filelist[filepath.Base(path)] = path
	}
	return nil
}

// findFile returns the full path to the file name.
func findFile(filename string) (string, error) {
	if fn, ok := filelist[filename]; ok {
		bag.Logger.Debugf("File lookup %q -> %q", filename, fn)
		return fn, nil
	}
	if _, err := os.Stat(filename); err == nil {
		var fn string
		fn, err = filepath.Abs(filename)
		if err != nil {
			return "", err
		}
		bag.Logger.Debugf("File lookup %q -> %q", filename, fn)
		return fn, nil
	}

	bag.Logger.Debugf("File lookup %q not found", filename)
	return "", fmt.Errorf("%w: %s", os.ErrNotExist, filename)
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

func dothings() error {
	var verbose bool
	op := optionparser.NewOptionParser()
	op.On("--verbose", "Print more debugging information", &verbose)
	op.Command("run", "Load layout and data files and create PDF")
	op.Command("version", "Print version information")
	err := op.Parse()
	if err != nil {
		op.Help()
		return err
	}
	if bag.Logger, err = newZapLogger(verbose); err != nil {
		return err
	}

	cmd := "run"
	if len(op.Extra) > 0 {
		cmd = op.Extra[0]
	}

	switch cmd {
	case "run":
		initDirs()
		var layoutpath, datapath string
		var lr, dr io.ReadCloser
		var pw io.WriteCloser
		if layoutpath, err = findFile("layout.xml"); err != nil {
			return err
		}
		if lr, err = os.Open(layoutpath); err != nil {
			return err
		}
		if datapath, err = findFile("data.xml"); err != nil {
			return err
		}
		if dr, err = os.Open(datapath); err != nil {
			return err
		}
		if pw, err = os.Create("publisher.pdf"); err != nil {
			return err
		}

		xc := &core.XTSCofig{
			Layoutfile:  lr,
			Datafile:    dr,
			Outfile:     pw,
			OutFilename: "publisher.pdf",
			FindFile:    findFile,
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
		fmt.Println(err)
		os.Exit(1)
	}
}
