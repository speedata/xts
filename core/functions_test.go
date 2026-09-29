package core

import (
	"strings"
	"testing"

	xpath "github.com/speedata/goxpath"
)

// TestCurrentRowBeforeFirstPage checks that sd:current-row() can be asked
// before anything has created a page.
func TestCurrentRowBeforeFirstPage(t *testing.T) {
	parser, err := xpath.NewParser(strings.NewReader("<root/>"))
	if err != nil {
		t.Fatal(err)
	}
	xd := &xtsDocument{data: parser}
	parser.Ctx.Store = map[any]any{"xd": xd}
	seq, err := fnCurrentRow(parser.Ctx, nil)
	if err != nil {
		t.Fatalf("fnCurrentRow: %v", err)
	}
	if len(seq) != 1 || seq[0] != 1 {
		t.Errorf("sd:current-row() = %v, want 1", seq)
	}
	if xd.currentPage != nil {
		t.Error("sd:current-row() created a page")
	}
}

// TestGridSizeBeforeFirstPage checks that sd:number-of-columns() and
// sd:number-of-rows() report an error rather than panic before the first page.
func TestGridSizeBeforeFirstPage(t *testing.T) {
	parser, err := xpath.NewParser(strings.NewReader("<root/>"))
	if err != nil {
		t.Fatal(err)
	}
	xd := &xtsDocument{data: parser}
	parser.Ctx.Store = map[any]any{"xd": xd}
	if _, err := fnNumberOfColumns(parser.Ctx, nil); err == nil {
		t.Error("sd:number-of-columns() before the first page: no error")
	}
	if _, err := fnNumberOfRows(parser.Ctx, nil); err == nil {
		t.Error("sd:number-of-rows() before the first page: no error")
	}
}
