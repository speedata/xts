package core

import (
	"bytes"
	"fmt"
	"math"
	"strings"

	pdf "github.com/boxesandglue/baseline-pdf"
)

// parsePlaceholderURL extracts width and height in PDF points from a
// placeholder:// URL. Format: placeholder://WxH (e.g. placeholder://200x150).
func parsePlaceholderURL(href string) (float64, float64, error) {
	spec := strings.TrimPrefix(href, "placeholder://")
	var w, h float64
	n, err := fmt.Sscanf(spec, "%fx%f", &w, &h)
	if err != nil || n != 2 {
		return 0, 0, fmt.Errorf("invalid placeholder URL %q, expected placeholder://WxH", href)
	}
	if w <= 0 || h <= 0 {
		return 0, 0, fmt.Errorf("placeholder dimensions must be positive, got %g x %g", w, h)
	}
	return w, h, nil
}

// loadPlaceholderImage generates a placeholder PDF in memory and loads it as an
// Imagefile via the document's PDF writer.
func loadPlaceholderImage(xd *xtsDocument, href string) (*pdf.Imagefile, error) {
	w, h, err := parsePlaceholderURL(href)
	if err != nil {
		return nil, err
	}
	pdfData, err := generatePlaceholderPDF(w, h)
	if err != nil {
		return nil, fmt.Errorf("generating placeholder PDF: %w", err)
	}
	reader := bytes.NewReader(pdfData)
	return xd.document.Doc.LoadImageFromReader(reader, "/MediaBox", 1)
}

// generatePlaceholderPDF creates a minimal single-page PDF with placeholder
// graphics: a light gray background, border, diagonal cross, a colored circle,
// and dimension text using the built-in Helvetica font.
func generatePlaceholderPDF(w, h float64) ([]byte, error) {
	var buf bytes.Buffer
	pw := pdf.NewPDFWriter(&buf)
	pw.DefaultPageWidth = w
	pw.DefaultPageHeight = h

	cs := pw.NewObject()
	writeContentStream(cs.Data, w, h)

	pg := pw.AddPage(cs, 0)
	pg.Width = w
	pg.Height = h
	pg.Dict = pdf.Dict{
		"Resources": pdf.Dict{
			"Font": pdf.Dict{
				pdf.Name("F1"): "<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>",
			},
		},
	}

	if err := pw.Finish(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// writeContentStream writes PDF drawing operators for the placeholder visual.
func writeContentStream(buf *bytes.Buffer, w, h float64) {
	cx, cy := w/2, h/2
	r := math.Min(w, h) * 0.15

	// Light gray background
	fmt.Fprintf(buf, "0.92 g\n")
	fmt.Fprintf(buf, "0 0 %s %s re f\n", ff(w), ff(h))

	// Border
	fmt.Fprintf(buf, "0.7 G 1 w\n")
	fmt.Fprintf(buf, "0.5 0.5 %s %s re S\n", ff(w-1), ff(h-1))

	// Diagonal cross
	fmt.Fprintf(buf, "0.8 G 0.5 w\n")
	fmt.Fprintf(buf, "0 0 m %s %s l S\n", ff(w), ff(h))
	fmt.Fprintf(buf, "%s 0 m 0 %s l S\n", ff(w), ff(h))

	// Colored circle at center (teal)
	if r > 2 {
		fmt.Fprintf(buf, "0.376 0.608 0.757 rg\n")
		writeBezierCircle(buf, cx, cy, r)
		fmt.Fprintf(buf, "f\n")
	}

	// Dimension text
	fontSize := math.Min(w, h) * 0.1
	fontSize = math.Max(6, math.Min(fontSize, 36))
	label := fmt.Sprintf("%g x %g", w, h)
	// Approximate text width: Helvetica averages about 0.52 * fontSize per
	// digit/symbol character. Count only ASCII runes since the label is simple.
	textWidth := float64(len(label)) * fontSize * 0.52
	tx := cx - textWidth/2
	ty := cy - fontSize/3

	fmt.Fprintf(buf, "0.3 g\n")
	fmt.Fprintf(buf, "BT\n")
	fmt.Fprintf(buf, "/F1 %s Tf\n", ff(fontSize))
	fmt.Fprintf(buf, "%s %s Td\n", ff(tx), ff(ty))
	fmt.Fprintf(buf, "(%s) Tj\n", pdfEscape(label))
	fmt.Fprintf(buf, "ET\n")
}

// writeBezierCircle approximates a circle with four cubic Bézier segments.
func writeBezierCircle(buf *bytes.Buffer, cx, cy, r float64) {
	const k = 0.5522847498 // kappa for circle approximation
	kr := k * r

	// Start at right (cx+r, cy)
	fmt.Fprintf(buf, "%s %s m\n", ff(cx+r), ff(cy))
	// Top quadrant
	fmt.Fprintf(buf, "%s %s %s %s %s %s c\n",
		ff(cx+r), ff(cy+kr), ff(cx+kr), ff(cy+r), ff(cx), ff(cy+r))
	// Left quadrant
	fmt.Fprintf(buf, "%s %s %s %s %s %s c\n",
		ff(cx-kr), ff(cy+r), ff(cx-r), ff(cy+kr), ff(cx-r), ff(cy))
	// Bottom quadrant
	fmt.Fprintf(buf, "%s %s %s %s %s %s c\n",
		ff(cx-r), ff(cy-kr), ff(cx-kr), ff(cy-r), ff(cx), ff(cy-r))
	// Right quadrant (close)
	fmt.Fprintf(buf, "%s %s %s %s %s %s c\n",
		ff(cx+kr), ff(cy-r), ff(cx+r), ff(cy-kr), ff(cx+r), ff(cy))
}

// ff formats a float for PDF content streams.
func ff(f float64) string {
	return pdf.FloatToPoint(f)
}

// pdfEscape escapes special characters for a PDF string literal.
func pdfEscape(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `(`, `\(`)
	s = strings.ReplaceAll(s, `)`, `\)`)
	return s
}
