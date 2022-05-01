package core

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/document"
	"github.com/speedata/boxesandglue/backend/lang"
	"github.com/speedata/boxesandglue/backend/node"
	"github.com/speedata/boxesandglue/csshtml"
	"github.com/speedata/boxesandglue/frontend"
	"github.com/speedata/goxml"
	xpath "github.com/speedata/goxpath"
)

var (
	errAttribNotFound = errors.New("Attribute not found")
	attributeValueRE  *regexp.Regexp
	oneCM             = bag.MustSp("1cm")
	// Version is a semantic version
	Version string
)

func init() {
	attributeValueRE = regexp.MustCompile(`\{(.*?)\}`)
}

type xtsDocument struct {
	cfg               *XTSCofig
	document          *frontend.Document
	layoutcss         *csshtml.CSS
	data              *xpath.Parser
	defaultLanguage   *lang.Lang
	pages             []*page
	groups            map[string]*group
	fontsources       map[string]*frontend.FontSource
	fontsizes         map[string][2]bag.ScaledPoint
	defaultGridWidth  bag.ScaledPoint
	defaultGridHeight bag.ScaledPoint
	defaultGridGapX   bag.ScaledPoint
	defaultGridGapY   bag.ScaledPoint
	defaultGridNx     int
	defaultGridNy     int
	pagetypes         []*pagetype
	currentPage       *page
	currentGrid       *grid
	currentGroup      *group
	tracing           VTrace
}

func newXTSDocument() *xtsDocument {
	return &xtsDocument{
		defaultGridWidth:  oneCM,
		defaultGridHeight: oneCM,
		defaultGridGapX:   0,
		defaultGridGapY:   0,
		layoutcss:         csshtml.NewCssParser(),
		groups:            make(map[string]*group),
	}
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
		bag.Logger.Error(err)
	}
	bag.Logger.Infof("Page %s created wd: %d, ht: %d", p.pagetype.name, p.pagegrid.nx, p.pagegrid.ny)
	xd.pages = append(xd.pages, p)
	xd.currentPage = p
	inSetupPage = false
	if f != nil {
		f()
	}
}

// XTSCofig is the configuration file for PDF generation.
type XTSCofig struct {
	Layoutfile  io.ReadCloser
	Datafile    io.ReadCloser
	Outfile     io.WriteCloser
	OutFilename string
	FindFile    func(string) (string, error)
}

// RunXTS is the entry point
func RunXTS(cfg *XTSCofig) error {
	starttime := time.Now()
	var err error
	var layoutxml *goxml.XMLDocument
	bag.Logger.Infof("XTS start version %s", Version)

	d := newXTSDocument()
	d.cfg = cfg
	if d.document, err = frontend.New(cfg.OutFilename); err != nil {
		return err
	}

	if err = d.defaultfont(); err != nil {
		return err
	}

	if layoutxml, err = goxml.Parse(cfg.Layoutfile); err != nil {
		return err
	}
	cfg.Layoutfile.Close()

	if d.data, err = xpath.NewParser(cfg.Datafile); err != nil {
		return err
	}
	d.data.Ctx.Store = map[interface{}]interface{}{
		"xd": d,
	}
	cfg.Datafile.Close()

	d.registerCallbacks()
	var defaultPagetype *pagetype
	if defaultPagetype, err = d.newPagetype("default page", "true()"); err != nil {
		return err
	}
	defaultPagetype.marginLeft = oneCM
	defaultPagetype.marginRight = oneCM
	defaultPagetype.marginTop = oneCM
	defaultPagetype.marginBottom = oneCM

	layoutRoot, err := layoutxml.Root()
	if err != nil {
		return err
	}

	dataNameSeq, err := d.data.Evaluate("local-name(/*)")
	if err != nil {
		return err
	}
	if len(dataNameSeq) != 1 {
		return fmt.Errorf("Could not find the root name for the data xml")
	}
	rootname := dataNameSeq[0].(string)
	_, err = dispatch(d, layoutRoot, d.data)
	if err != nil {
		return err
	}
	bag.Logger.Info("Start processing data")
	d.data.Ctx.Root()
	var startDispatcher *goxml.Element
	var ok bool
	if startDispatcher, ok = dataDispatcher[rootname][""]; !ok {
		bag.Logger.Errorf("Cannot find <Record> for root element %s", rootname)
		return fmt.Errorf("Cannot find <Record> for root element %s", rootname)
	}
	d.defaultLanguage, err = frontend.GetLanguage("en")
	if err != nil {
		return err
	}
	d.document.Doc.DefaultLanguage = d.defaultLanguage
	_, err = dispatch(d, startDispatcher, d.data)
	if err != nil {
		bag.Logger.Error(err)
		return err
	}
	if d.currentPage != nil {
		d.currentPage.bagPage.Shipout()
	}
	d.document.Doc.Finish()

	bag.Logger.Infof("Finished in %s", time.Now().Sub(starttime))
	return nil
}

// Add necessary callbacks to boxes and glue callback mechanism
// for tracing purpose.
func (xd *xtsDocument) registerCallbacks() {
	preShipout := func(pg *document.Page) {
		xtspage := pg.Userdata["xtspage"].(*page)
		curGrid := xtspage.pagegrid
		// Draw grid when requested
		if xd.IsTrace(VTraceAllocation) {
			vlist := node.NewVList()
			rule := node.NewRule()

			pdfinstructions := make([]string, 0, len(xd.currentGrid.allocatedBlocks))
			pdfinstructions = append(pdfinstructions, "q")

			for k, v := range curGrid.allocatedBlocks {
				if v > 0 {
					x, y := k.XY()
					pdfinstructions = append(pdfinstructions, fmt.Sprintf(" 1 1 0 rg %s %s %s %s re f", curGrid.posX(x, pageAreaName), xtspage.pageHeight-curGrid.posY(y, pageAreaName), curGrid.gridWidth, -curGrid.gridHeight))
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
			var pdfinstructions []string
			// page
			pdfinstructions = append(pdfinstructions, fmt.Sprintf("%s %s %s %s re S", x, y, wd, ht))
			gridMaxX := xtspage.pageWidth - curGrid.marginRight
			gridMaxY := xtspage.pageHeight - curGrid.marginTop
			pdfinstructions = append(pdfinstructions, "0.4 w")

			gridX := x + curGrid.gridWidth
			// vertical grid rules
			for i := 1; gridX < gridMaxX; i++ {
				switch {
				case i%10 == 0:
					pdfinstructions = append(pdfinstructions, "0.1 G")
				case i%5 == 0:
					pdfinstructions = append(pdfinstructions, "0.7 G")
				default:
					pdfinstructions = append(pdfinstructions, "0.9 G")
				}
				pdfinstructions = append(pdfinstructions, fmt.Sprintf("%s %s m %s %s l S", gridX, y, gridX, gridMaxY))
				gridX += xd.currentGrid.gridGapX
				if xd.currentGrid.gridGapX > 0 && gridX < gridMaxX {
					pdfinstructions = append(pdfinstructions, fmt.Sprintf("%s %s m %s %s l S", gridX, y, gridX, gridMaxY))
				}
				gridX += xd.currentGrid.gridWidth
			}

			// horizontal grid rules from top to bottom
			gridY := xtspage.pageHeight - curGrid.gridHeight - curGrid.marginTop
			for i := 1; gridY > y; i++ {
				switch {
				case i%10 == 0:
					pdfinstructions = append(pdfinstructions, "0.1 G")
				case i%5 == 0:
					pdfinstructions = append(pdfinstructions, "0.7 G")
				default:
					pdfinstructions = append(pdfinstructions, "0.9 G")
				}
				pdfinstructions = append(pdfinstructions, fmt.Sprintf("%s %s m %s %s l S", x, gridY, gridMaxX, gridY))
				gridY -= xd.currentGrid.gridGapY
				if xd.currentGrid.gridGapY > 0 && gridY > y {
					pdfinstructions = append(pdfinstructions, fmt.Sprintf("%s %s m %s %s l S", x, gridY, gridMaxX, gridY))
				}
				gridY -= xd.currentGrid.gridHeight
			}
			pageframe := fmt.Sprintf("0 0 %s %s re S", xtspage.pageWidth, xtspage.pageHeight)
			pdfinstructions = append(pdfinstructions, pageframe)

			pdfinstructions = append(pdfinstructions, " 2 w 1 0 0 RG ")
			for _, area := range curGrid.areas {
				for _, rect := range area.frame {
					posX := xd.currentGrid.posX(rect.col, pageAreaName)
					posY := xtspage.pageHeight - xd.currentGrid.posY(rect.row, pageAreaName)
					wd := xd.currentGrid.width(rect.width)
					ht := xd.currentGrid.height(rect.height) * -1
					frame := fmt.Sprintf("%s %s %s %s re S", posX, posY, wd, ht)
					pdfinstructions = append(pdfinstructions, frame)
				}
			}

			rule.Pre = strings.Join(pdfinstructions, " ")
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
