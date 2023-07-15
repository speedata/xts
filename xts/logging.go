package main

import (
	"context"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/xts/core"
	"golang.org/x/exp/slog"
)

var (
	lvl          = new(slog.LevelVar)
	enc          *xml.Encoder
	msgAttr      = xml.Attr{Name: xml.Name{Local: "msg"}}
	lvlAttr      = xml.Attr{Name: xml.Name{Local: "level"}}
	logElement   = xml.StartElement{Name: xml.Name{Local: "entry"}}
	protocolFile io.Writer
	errCount     = 0
	warnCount    = 0
)

type logHandler struct {
}

func (lh *logHandler) Enabled(_ context.Context, level slog.Level) bool {
	if level == slog.LevelError {
		errCount++
	} else if level == slog.LevelWarn {
		warnCount++
	}
	return level >= bag.LogLevel.Level()
}

func (lh *logHandler) Handle(_ context.Context, r slog.Record) error {
	lvlAttr.Value = r.Level.String()
	if r.Level == core.LevelNotice {
		lvlAttr.Value = "NOTICE"
	}
	msgAttr.Value = r.Message
	values := []string{}
	le := logElement.Copy()
	le.Attr = append(le.Attr, lvlAttr, msgAttr)
	r.Attrs(
		func(a slog.Attr) bool {
			var val string
			switch t := a.Value.Any().(type) {
			case slog.LogValuer:
				val = t.LogValue().String()
				values = append(values, fmt.Sprintf("%s=%s", a.Key, val))
			default:
				t = a.Value
				val = a.Value.String()
				values = append(values, fmt.Sprintf("%s=%s", a.Key, a.Value))
			}
			le.Attr = append(le.Attr, xml.Attr{Name: xml.Name{Local: a.Key}, Value: val})
			return true
		})
	enc.EncodeToken(xml.CharData([]byte("  ")))
	enc.EncodeToken(le)
	enc.EncodeToken(le.End())
	enc.EncodeToken(xml.CharData([]byte("\n")))
	if configuration.Verbose {
		lparen := ""
		rparen := ""
		if len(values) > 0 {
			lparen = "("
			rparen = ")"
		}
		fmt.Println(r.Message, lparen+strings.Join(values, ",")+rparen)
	}
	return nil
}

func (lh *logHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return lh
}

func (lh *logHandler) WithGroup(name string) slog.Handler {
	return lh
}

// protocol is the name of the XML protocol file
func setupLog(protocol string) error {
	var err error

	protocolFile, err = os.Create(protocol)
	if err != nil {
		return err
	}
	sl := slog.New(&logHandler{})
	slog.SetDefault(sl)

	enc = xml.NewEncoder(protocolFile)
	if err = enc.EncodeToken(xml.StartElement{
		Name: xml.Name{Local: "log"},
	}); err != nil {
		return err
	}
	enc.EncodeToken(xml.CharData([]byte("\n")))

	return nil
}

func teardownLog() error {
	if err := enc.EncodeToken(xml.EndElement{
		Name: xml.Name{Local: "log"},
	}); err != nil {
		return err
	}
	if err := enc.Flush(); err != nil {
		return err
	}
	return nil
}

// see the section on performance considerations in the slog package for a
// rationale why I chose this way to do the calculation.
type md5calc string

func (e md5calc) LogValue() slog.Value {
	data, err := os.ReadFile(string(e))
	if err != nil {
		return slog.StringValue(err.Error())
	}
	return slog.AnyValue(fmt.Sprintf("%x", md5.Sum(data)))
}
