package core

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/speedata/boxesandglue/backend/bag"
	"github.com/speedata/boxesandglue/backend/lang"
	"github.com/speedata/boxesandglue/backend/node"
	"github.com/speedata/boxesandglue/csshtml"
	"github.com/speedata/boxesandglue/document"
	"github.com/speedata/goxml"
	"github.com/speedata/goxpath/xpath"
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
	doc               *document.Document
	layoutcss         *csshtml.CSS
	data              *xpath.Parser
	defaultLanguage   *lang.Lang
	pages             []*page
	fontsources       map[string]*document.FontSource
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
	tracing           VTrace
}

func newXTSDocument() *xtsDocument {
	return &xtsDocument{
		defaultGridWidth:  oneCM,
		defaultGridHeight: oneCM,
		defaultGridGapX:   0,
		defaultGridGapY:   0,
		layoutcss:         csshtml.NewCssParser(),
	}
}

func (xd *xtsDocument) setupPage() {
	if xd.currentPage != nil {
		return
	}
	p, err := newPage(xd)
	if err != nil {
		bag.Logger.Error(err)
	}
	bag.Logger.Infof("Create page %s", p.pagetype.name)
	xd.pages = append(xd.pages, p)
	xd.currentPage = p
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

	if layoutxml, err = goxml.Parse(cfg.Layoutfile); err != nil {
		return err
	}
	cfg.Layoutfile.Close()

	if d.data, err = xpath.NewParser(cfg.Datafile); err != nil {
		return err
	}
	cfg.Datafile.Close()

	d.doc = document.NewDocument(cfg.Outfile)
	d.doc.Filename = cfg.OutFilename
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
	d.defaultLanguage, err = d.doc.GetLanguage("en")
	if err != nil {
		return err
	}
	d.doc.SetDefaultLanguage(d.defaultLanguage)
	_, err = dispatch(d, startDispatcher, d.data)
	if err != nil {
		return err
	}
	d.doc.CurrentPage.Shipout()
	d.doc.Finish()
	cfg.Outfile.Close()
	bag.Logger.Infof("Finished in %s", time.Now().Sub(starttime))
	return nil
}

func (xd *xtsDocument) registerCallbacks() {
	preShipout := func(pg *document.Page) {
		xtspage := pg.Userdata["xtspage"].(*page)
		// Draw grid when requested
		if xd.IsTrace(VTraceAllocation) {
			vlist := node.NewVList()
			rule := node.NewRule()

			pdfinstructions := make([]string, 0, len(xd.currentGrid.allocatedBlocks))
			pdfinstructions = append(pdfinstructions, "q")
			curGrid := xtspage.pagegrid

			for k, v := range curGrid.allocatedBlocks {
				if v > 0 {
					x, y := k.XY()
					pdfinstructions = append(pdfinstructions, fmt.Sprintf(" 1 1 0 rg %s %s %s %s re f", curGrid.posX(x), xtspage.pageHeight-curGrid.posY(y), curGrid.gridWidth, -curGrid.gridHeight))
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
			x := xtspage.marginLeft
			y := xtspage.marginBottom
			wd := xtspage.pageWidth - xtspage.marginLeft - xtspage.marginRight
			ht := xtspage.pageHeight - xtspage.marginTop - xtspage.marginBottom
			var pdfinstructions []string
			// page
			pdfinstructions = append(pdfinstructions, fmt.Sprintf("%s %s %s %s re S", x, y, wd, ht))
			gridHeight := xtspage.pageHeight - xtspage.marginTop
			gridWidth := xtspage.pageWidth - xtspage.marginRight
			pdfinstructions = append(pdfinstructions, "0.4 w")

			gridX := x + xtspage.pagegrid.gridWidth
			// vertical grid rules
			for i := 1; gridX < gridWidth; i++ {
				if i%5 == 0 {
					pdfinstructions = append(pdfinstructions, "0.5 G")
				} else {
					pdfinstructions = append(pdfinstructions, "0.9 G")
				}
				pdfinstructions = append(pdfinstructions, fmt.Sprintf("%s %s m %s %s l S", gridX, y, gridX, gridHeight))
				gridX += xd.currentGrid.gridGapX
				if xd.currentGrid.gridGapX > 0 && gridX < gridWidth {
					pdfinstructions = append(pdfinstructions, fmt.Sprintf("%s %s m %s %s l S", gridX, y, gridX, gridHeight))
				}
				gridX += xd.currentGrid.gridWidth
			}

			// horizontal grid rules from top to bottom
			gridY := xtspage.pageHeight - xtspage.pagegrid.gridHeight - xtspage.marginTop
			for i := 1; gridY > y; i++ {
				if i%5 == 0 {
					pdfinstructions = append(pdfinstructions, "0.5 G")
				} else {
					pdfinstructions = append(pdfinstructions, "0.9 G")
				}
				pdfinstructions = append(pdfinstructions, fmt.Sprintf("%s %s m %s %s l S", x, gridY, gridWidth, gridY))
				gridY -= xd.currentGrid.gridGapY
				if xd.currentGrid.gridGapY > 0 && gridY > y {
					pdfinstructions = append(pdfinstructions, fmt.Sprintf("%s %s m %s %s l S", x, gridY, gridWidth, gridY))
				}
				gridY -= xd.currentGrid.gridHeight
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

	xd.doc.RegisterCallback(document.CallbackPreShipout, preShipout)
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
