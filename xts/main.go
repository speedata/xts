package main

import (
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/optionparser"
	"github.com/speedata/xts/core"
)

var (
	version string
)

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
		core.InitDirs()
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

		xc := &core.XTSCofig{
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
		os.Exit(1)
	}
}
