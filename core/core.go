package core

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/slog"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/document"
	"github.com/speedata/boxesandglue/backend/node"
	"github.com/speedata/boxesandglue/csshtml"
	"github.com/speedata/boxesandglue/frontend"
	"github.com/speedata/boxesandglue/frontend/pdfdraw"
	"github.com/speedata/goxml"
	xpath "github.com/speedata/goxpath"
)

var (
	// XPath escape sequence for attributes
	attributeValueRE   = regexp.MustCompile(`\{(.*?)\}`)
	oneCM              = bag.MustSp("1cm")
	destinationNumbers = make(chan int)
	// Version is a semantic version
	Version string
)

const (
	// SDNAMESPACE is the speedata XTS layout rules namespace
	SDNAMESPACE string = "urn:speedata.de/2021/xts/en"
	// LevelNotice is used for messages from Message
	LevelNotice = slog.Level(2)
)

func init() {
	go genIntegerSequence(destinationNumbers)
}

func genIntegerSequence(ids chan int) {
	i := int(0)
	for {
		ids <- i
		i++
	}
}

type xtsDocument struct {
	cfg               *XTSConfig
	document          *frontend.Document
	layoutcss         *csshtml.CSS
	data              *xpath.Parser
	pages             []*page
	groups            map[string]*group
	defaultGridWidth  bag.ScaledPoint
	defaultGridHeight bag.ScaledPoint
	defaultGridGapX   bag.ScaledPoint
	defaultGridGapY   bag.ScaledPoint
	defaultGridNx     int
	defaultGridNy     int
	masterpages       []*pagetype
	marker            mapmarker
	aux               *auxfile // contents of the previous run, if available
	currentPage       *page
	currentGrid       *grid
	currentGroup      *group
	currentPagenumber int
	tracing           VTrace
	layoutNS          map[string]string
	// for “global” variables
	store map[any]any
}

func newXTSDocument() *xtsDocument {
	xd := &xtsDocument{
		defaultGridWidth:  oneCM,
		defaultGridHeight: oneCM,
		defaultGridGapX:   0,
		defaultGridGapY:   0,
		layoutcss:         csshtml.NewCSSParserWithDefaults(),
		groups:            make(map[string]*group),
		store:             make(map[any]any),
		marker:            make(mapmarker),
	}
	return xd
}

// Check if requestedVersion can be used in productVersion.
func checkVersion(requestedVersion, productVersion string) error {
	if requestedVersion == "" {
		// no version information in the layout file, ok!
		return nil
	}

	xtsVersionSplit := strings.Split(productVersion, ".")
	if len(xtsVersionSplit) != 3 {
		return fmt.Errorf("XTS version %q looks incorrect", productVersion)
	}
	if strings.Contains(xtsVersionSplit[2], "-") {
		// this is probably a SHA1 based development version and should always work.
		return nil
	}
	var xtsVersionArray [3]int
	var err error
	for i, v := range xtsVersionSplit {
		if xtsVersionArray[i], err = strconv.Atoi(v); err != nil {
			return err
		}
	}

	layoutVersionSplit := strings.Split(requestedVersion, ".")
	if len(layoutVersionSplit) > 0 {
		if i, err := strconv.Atoi(layoutVersionSplit[0]); err == nil {
			if i > xtsVersionArray[0] {
				goto versionMismatch
			} else if i < xtsVersionArray[0] {
				return nil
			}
		} else {
			return err
		}
	}
	if len(requestedVersion) > 1 {
		if i, err := strconv.Atoi(layoutVersionSplit[1]); err == nil {
			if i > xtsVersionArray[1] {
				goto versionMismatch
			} else if i < xtsVersionArray[1] {
				return nil
			}
		} else {
			return err
		}
	}
	if len(requestedVersion) > 2 {
		if i, err := strconv.Atoi(layoutVersionSplit[2]); err == nil {
			if i > xtsVersionArray[2] {
				goto versionMismatch
			}
		} else {
			return err
		}
	}
	return nil
versionMismatch:
	return fmt.Errorf("requested layout version %q and xts version %q don't match", requestedVersion, productVersion)
}

var inSetupPage bool

func (xd *xtsDocument) setupPage() {
	if xd.currentGroup != nil {
		return
	}
	if xd.currentPage != nil {
		return
	}
	if inSetupPage {
		return
	}
	inSetupPage = true
	p, f, err := newPage(xd)
	if err != nil {
		slog.Error(err.Error())
	}
	slog.Info(fmt.Sprintf("Page %d type %s created wd: %d, ht: %d grid cells", p.pagenumber, p.pagetype.name, p.pagegrid.nx, p.pagegrid.ny))
	xd.pages = append(xd.pages, p)
	xd.currentPage = p
	inSetupPage = false
	if f != nil {
		f()
	}
}

// XTSConfig is the configuration file for PDF generation.
type XTSConfig struct {
	Datafile     io.Reader
	DumpFile     io.Writer
	FindFile     func(string) (string, error)
	Jobname      string
	Layoutfile   io.Reader
	Mode         []string
	Outfile      io.WriteCloser
	OutFilename  string
	SuppressInfo bool
	Tracing      []string
	Variables    map[string]any
}

// RunXTS is the entry point
func RunXTS(cfg *XTSConfig) error {
	var err error
	var layoutxml *goxml.XMLDocument
	slog.Info(fmt.Sprintf("XTS start version %s", Version))

	d := newXTSDocument()
	d.cfg = cfg
	if d.document, err = frontend.New(cfg.OutFilename); err != nil {
		return err
	}
	bag.SetLogger(slog.Default())
	if cfg.SuppressInfo {
		d.document.SetSuppressInfo(true)
		slog.Info("Creating reproducible build")
	}
	curWD, err := os.Getwd()
	if err != nil {
		return err
	}
	d.layoutcss.FrontendDocument = d.document
	d.layoutcss.PushDir(curWD)

	for _, tr := range cfg.Tracing {
		switch tr {
		case "grid":
			d.SetVTrace(VTraceGrid)
		case "gridallocation":
			d.SetVTrace(VTraceAllocation)
		}
	}
	d.document.Doc.CompressLevel = 9
	slog.Info("Setup defaults ...")
	if err = d.defaultfont(); err != nil {
		return err
	}

	var defaultPagetype *pagetype
	if defaultPagetype, err = d.newPagetype("default page", "true()"); err != nil {
		return err
	}
	defaultPagetype.marginLeft = oneCM
	defaultPagetype.marginRight = oneCM
	defaultPagetype.marginTop = oneCM
	defaultPagetype.marginBottom = oneCM

	slog.Info("Setup defaults ... done")
	if layoutxml, err = goxml.Parse(cfg.Layoutfile); err != nil {
		return err
	}

	if d.data, err = xpath.NewParser(cfg.Datafile); err != nil {
		return err
	}
	if d.aux, err = d.readAuxFile(); err != nil {
		return err
	}

	// connect the main document to the xpath parser, so we can access the
	// document related information in the layout functions
	d.data.Ctx.Store = map[any]any{
		"xd": d,
	}
	for k, v := range cfg.Variables {
		d.data.SetVariable(k, xpath.Sequence{v})
	}
	d.registerCallbacks()

	layoutRoot, err := layoutxml.Root()
	if err != nil {
		return err
	}
	// Check if the layout file is in the correct name space
	if ns := layoutRoot.Namespaces[layoutRoot.Prefix]; ns != SDNAMESPACE {
		if ns == "" {
			ns = "none"
		}
		return newTypesettingErrorFromStringf("the layout file must be in the name space %s, found %s", SDNAMESPACE, ns)
	}
	for _, attr := range layoutRoot.Attributes() {
		if attr.Name == "version" {
			if err = checkVersion(attr.Value, Version); err != nil {
				return err
			}
			break
		}
	}

	d.document.Doc.DefaultLanguage, err = frontend.GetLanguage("en")
	if err != nil {
		return err
	}

	dataNameSeq, err := d.data.Evaluate("local-name(/*)")
	if err != nil {
		return err
	}
	if len(dataNameSeq) != 1 {
		return newTypesettingErrorFromString("Could not find the root name for the data xml")
	}
	slog.Info("Start processing data")

	rootname := dataNameSeq[0].(string)
	_, err = dispatch(d, layoutRoot, d.data)
	if err != nil {
		return newTypesettingError(err)
	}

	d.data.Ctx.Root()
	var startDispatcher *goxml.Element
	var ok bool
	if startDispatcher, ok = dataDispatcher[rootname][""]; !ok {
		return newTypesettingErrorFromString(fmt.Sprintf("Cannot find <Record> for root element %s", rootname))
	}
	_, err = dispatch(d, startDispatcher, d.data)
	if err != nil {
		return newTypesettingError(err)
	}
	if d.currentPage != nil {
		d.currentPage.bagPage.Shipout()
	}
	if err = d.document.Finish(); err != nil {
		return err
	}
	if cfg.DumpFile != nil {
		d.document.Doc.OutputXMLDump(cfg.DumpFile)
		if closer, ok := cfg.DumpFile.(io.WriteCloser); ok {
			err = closer.Close()
			if err != nil {
				return err
			}
		}
	}
	if err = d.writeAuxXML(); err != nil {
		return err
	}
	return nil
}

// Add necessary callbacks to boxes and glue callback mechanism for tracing
// purpose.
func (xd *xtsDocument) registerCallbacks() {
	preShipout := func(pg *document.Page) {
		xtspage := pg.Userdata["xtspage"].(*page)
		curGrid := xtspage.pagegrid
		pageArea := curGrid.areas[pageAreaName]
		// Draw grid when requested
		if xd.IsTrace(VTraceAllocation) {
			vlist := node.NewVList()
			rule := node.NewRule()
			pdfinstructions := make([]string, 0, len(xd.currentGrid.allocatedBlocks))
			pdfinstructions = append(pdfinstructions, "q 1 1 0 rg")

			for k, v := range curGrid.allocatedBlocks {
				if v > 0 {
					x, y := k.XY()
					pdfinstructions = append(pdfinstructions, fmt.Sprintf("%s %s %s %s re f", curGrid.posX(x, pageArea), xtspage.pageHeight-curGrid.posY(y, pageArea), curGrid.gridWidth, -curGrid.gridHeight))
				}
			}
			pdfinstructions = append(pdfinstructions, " Q")
			rule.Pre = strings.Join(pdfinstructions, " ")

			vlist.List = node.Hpack(rule)
			pg.Background = append(pg.Background, document.Object{Vlist: vlist, X: 0, Y: 0})

		}
		if xd.IsTrace(VTraceGrid) {
			vlist := node.NewVList()
			rule := node.NewRule()
			x := curGrid.marginLeft
			y := curGrid.marginBottom
			wd := xtspage.pageWidth - curGrid.marginLeft - curGrid.marginRight
			ht := xtspage.pageHeight - curGrid.marginTop - curGrid.marginBottom
			// page
			halfpt := bag.ScaledPointFromFloat(0.5)
			lightgray := *xd.document.GetColor("lightgray")
			gray := *xd.document.GetColor("gray")
			darkgray := *xd.document.GetColor("darkgray")
			red := *xd.document.GetColor("red")
			darkblue := *xd.document.GetColor("darkblue")

			gridDebug := pdfdraw.NewStandalone().LineWidth(halfpt).ColorStroking(darkblue).Rect(x, y, wd, ht).Stroke()

			// pdfinstructions = append(pdfinstructions, rect)
			gridMaxX := xtspage.pageWidth - curGrid.marginRight
			gridDebug.LineWidth(halfpt).SetDash([]uint{2}, 0)

			gridX := x
			// vertical grid rules
			for i := 0; gridX <= gridMaxX; i++ {
				switch {
				case i%10 == 0:
					gridDebug.ColorStroking(darkgray)
				case i%5 == 0:
					gridDebug.ColorStroking(gray)
				default:
					gridDebug.ColorStroking(lightgray)
				}
				gridDebug.Moveto(gridX, 0).Lineto(gridX, xtspage.pageHeight).Stroke()
				if xd.currentGrid.gridGapX > 0 && gridX < gridMaxX && i > 0 {
					gridX += xd.currentGrid.gridGapX
					gridDebug.Moveto(gridX, 0).Lineto(gridX, xtspage.pageHeight).Stroke()
				}
				gridX += xd.currentGrid.gridWidth
			}

			// horizontal grid rules from top to bottom
			gridY := xtspage.pageHeight - curGrid.marginTop
			for i := 0; gridY >= y; i++ {
				switch {
				case i%10 == 0:
					gridDebug.ColorStroking(darkgray)
				case i%5 == 0:
					gridDebug.ColorStroking(gray)
				default:
					gridDebug.ColorStroking(lightgray)
				}
				gridDebug.Moveto(0, gridY).Lineto(xtspage.pageWidth, gridY).Stroke()
				if xd.currentGrid.gridGapY > 0 && gridY > y && i > 0 {
					gridY -= xd.currentGrid.gridGapY
					gridDebug.Moveto(0, gridY).Lineto(xtspage.pageWidth, gridY).Stroke()
				}
				gridY -= xd.currentGrid.gridHeight
			}
			gridDebug.Rect(0, 0, xtspage.pageWidth, xtspage.pageHeight).Stroke().SetDash([]uint{}, 0).LineWidth(bag.ScaledPointFromFloat(0.4)).ColorStroking(red)

			for _, area := range curGrid.areas {
				for _, rect := range area.frame {
					posX := xd.currentGrid.posX(1, area)
					posY := xtspage.pageHeight - xd.currentGrid.posY(1, area)
					wd := xd.currentGrid.width(rect.width)
					ht := xd.currentGrid.height(rect.height) * -1
					gridDebug.Rect(posX, posY, wd, ht).Stroke()
				}
			}

			rule.Pre = gridDebug.String()
			vlist.List = node.Hpack(rule)
			pg.Background = append(pg.Background, document.Object{Vlist: vlist, X: 0, Y: 0})

		}
		if xd.IsTrace(VTraceHyphenation) {
			for _, v := range pg.Objects {
				showDiscNodes(v.Vlist.List)
			}
		}
	}

	xd.document.Doc.RegisterCallback(document.CallbackPreShipout, preShipout)
}

func showDiscNodes(n node.Node) {
	for e := n; e != nil; e = e.Next() {
		switch t := e.(type) {
		case *node.HList:
			showDiscNodes(t.List)
		case *node.Disc:
			r := node.NewRule()
			r.Pre = "q 0.3 w 0 2 m 0 7 l S Q"
			node.InsertAfter(n, e, r)
		default:
			// ignore
		}
	}
}

// A TypesettingError contains the information if it has been logged, so it does
// not appear more than once in the output.
type TypesettingError struct {
	Logged bool
	Msg    string
}

func (te TypesettingError) Error() string {
	return te.Msg
}

func newTypesettingError(err error) error {
	if terr, ok := err.(TypesettingError); ok {
		if !terr.Logged {
			slog.Error(terr.Msg)
			terr.Logged = true
		}
		return terr
	}
	return TypesettingError{
		Msg: err.Error(),
	}
}

func newTypesettingErrorFromString(msg string) error {
	slog.Error(msg)
	return TypesettingError{
		Msg:    msg,
		Logged: true,
	}
}

func newTypesettingErrorFromStringf(format string, a ...any) error {
	return newTypesettingErrorFromString(fmt.Sprintf(format, a...))
}
