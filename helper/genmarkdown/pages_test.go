package genmarkdown

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnchorize(t *testing.T) {
	for heading, want := range map[string]string{
		"Grid allocation":                      "grid-allocation",
		"Tables across pages":                  "tables-across-pages",
		"`<Flow>` and the grid":                "flow-and-the-grid",
		"What's next?":                         "whats-next",
		"Size and [line height](/x/y)":         "size-and-line-height",
		"Font size & line height":              "font-size--line-height",
		"break-before: column":                 "break-before-column",
		"The `bottom` attribute, step by step": "the-bottom-attribute-step-by-step",
	} {
		if got := anchorize(heading); got != want {
			t.Errorf("anchorize(%q) = %q, want %q", heading, got, want)
		}
	}
}

func TestPageTitle(t *testing.T) {
	dir := t.TempDir()
	write := func(name, text string) {
		t.Helper()
		fn := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(fn), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fn, []byte(text), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("manual/a/grid.md", "---\nweight: 20\nlinktitle: The Grid\n---\n\n# The Grid itself\n\n## Grid allocation\n\n```xml\n# not a heading\n```\n\n## Own id {#own}\n")
	write("manual/a/_index.md", "---\ntype: docs\n---\n\n# Core Concepts\n")

	for _, tc := range []struct {
		href, title string
		fails       bool
	}{
		{"/manual/a/grid", "The Grid", false},
		{"/manual/a/grid/", "The Grid", false},
		{"/manual/a/grid#grid-allocation", "The Grid: Grid allocation", false},
		{"/manual/a/grid#own", "The Grid: Own id", false},
		{"/manual/a/grid#not-a-heading", "The Grid", true},
		{"/manual/a/grid#own-id", "The Grid", true},
		{"/manual/a", "Core Concepts", false},
		{"/manual/a/missing", "", true},
		{"manual/a/grid", "", true},
		{"/manual/a/grid.md", "", true},
	} {
		title, err := pageTitle(dir, tc.href)
		if (err != nil) != tc.fails {
			t.Errorf("pageTitle(%q): err = %v, want failure %t", tc.href, err, tc.fails)
		}
		if title != tc.title {
			t.Errorf("pageTitle(%q) = %q, want %q", tc.href, title, tc.title)
		}
	}
}
