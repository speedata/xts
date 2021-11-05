package core

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/document"
	"github.com/speedata/goxml"
	"github.com/speedata/goxpath/xpath"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger            *zap.SugaredLogger
	errAttribNotFound = errors.New("Attribute not found")
	attributeValueRE  *regexp.Regexp
	doc               *document.Document
)

func init() {
	attributeValueRE = regexp.MustCompile(`\{(.*?)\}`)
}

func newZapLogger() *zap.SugaredLogger {
	logger, _ := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding:    "console",
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			EncodeLevel: zapcore.LowercaseColorLevelEncoder,
			LevelKey:    "level",
			MessageKey:  "message",
		},
	}.Build()
	return logger.Sugar()
}

// Dothings is the entry point
func Dothings() error {
	starttime := time.Now()
	logger = newZapLogger()
	bag.Logger = logger

	logger.Info("XTS start")
	layoutReader, err := os.Open("layout.xml")
	if err != nil {
		return err
	}
	layoutxml, err := goxml.Parse(layoutReader)
	if err != nil {
		return err
	}
	layoutReader.Close()

	dataReader, err := os.Open("data.xml")
	if err != nil {
		return err
	}
	data, err := xpath.NewParser(dataReader)
	if err != nil {
		return err
	}
	dataReader.Close()

	w, err := os.Create("out.pdf")
	if err != nil {
		return err
	}
	doc = document.NewDocument(w)
	doc.Filename = "out.pdf"

	layoutRoot, err := layoutxml.Root()
	if err != nil {
		return err
	}

	dataNameSeq, err := data.Evaluate("local-name(/*)")
	if err != nil {
		return err
	}
	if len(dataNameSeq) != 1 {
		return fmt.Errorf("Could not find the root name for the data xml")
	}
	rootname := dataNameSeq[0].(string)
	_, err = dispatch(layoutRoot, data)
	if err != nil {
		return err
	}
	logger.Info("Start processing data")
	data.Ctx.Root()
	var startDispatcher *goxml.Element
	var ok bool
	if startDispatcher, ok = dataDispatcher[rootname][""]; !ok {
		logger.Errorf("Cannot find <Record> for root element %s", rootname)
		return fmt.Errorf("Cannot find <Record> for root element %s", rootname)
	}

	dispatch(startDispatcher, data)
	doc.CurrentPage.Shipout()
	doc.Finish()
	w.Close()
	logger.Infof("Finished in %s", time.Now().Sub(starttime))
	return nil
}

func findAttribute(name string, element *goxml.Element, mustexist bool, allowXPath bool, dflt string, xp *xpath.Parser) (string, error) {
	var value string
	var found bool
	for _, attrib := range element.Attributes() {
		if attrib.Name == name {
			found = true
			value = attrib.Value
			break
		}
	}
	if !found {
		if mustexist {
			logger.Errorf("Layout line %d: attribute %s on element %s not found", element.Line, name, element.Name)
			return "", fmt.Errorf("line %d: attribute %s on element %s not found", element.Line, name, element.Name)
		}
		value = dflt
	}

	value = attributeValueRE.ReplaceAllStringFunc(value, func(a string) string {
		// strip curly braces
		seq, err := xp.Evaluate(a[1 : len(a)-1])
		if err != nil {
			logger.Errorf("Layout line %d: %s", element.Line, err)
			return ""
		}
		return seq.Stringvalue()
	})
	return value, nil
}

func getAttributeString(name string, element *goxml.Element, mustexist bool, allowXPath bool, dflt string, xp *xpath.Parser) (string, error) {
	return findAttribute(name, element, mustexist, allowXPath, dflt, xp)
}

func getAttributeSize(name string, element *goxml.Element, mustexist bool, allowXPath bool, dflt string, xp *xpath.Parser) (bag.ScaledPoint, error) {
	val, err := findAttribute(name, element, mustexist, allowXPath, dflt, xp)
	if err != nil {
		return 0, err
	}
	return bag.Sp(val)
}
