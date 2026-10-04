package genmarkdown

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

var (
	headingLine = regexp.MustCompile(`^#{1,6}\s+(.*?)\s*$`)
	headingID   = regexp.MustCompile(`\s*\{#([^}]+)\}$`)
	mdLink      = regexp.MustCompile(`\[([^\]]*)\]\([^)]*\)`)
)

// pageTitle returns the title of the manual page that href (such as
// /manual/core-concepts/grid#grid-allocation) points to, read from the
// linktitle of its front matter or its first heading, followed by the heading
// of the anchor: "The Grid: Grid allocation". contentdir is the manual's
// content directory. It is an error when the page does not exist or has no
// heading for the anchor.
func pageTitle(contentdir, href string) (string, error) {
	path, anchor, _ := strings.Cut(href, "#")
	if !strings.HasPrefix(path, "/") || strings.HasSuffix(path, ".md") {
		return "", fmt.Errorf("page %q: write the path from the manual's root, such as /manual/core-concepts/grid", href)
	}
	rel := filepath.FromSlash(strings.Trim(path, "/"))
	var data []byte
	var err error
	for _, fn := range []string{rel + ".md", filepath.Join(rel, "_index.md")} {
		if data, err = os.ReadFile(filepath.Join(contentdir, fn)); err == nil {
			break
		}
	}
	if err != nil {
		return "", fmt.Errorf("page %q: no such page in the manual", href)
	}
	title, headings := scanPage(data)
	if anchor == "" {
		return title, nil
	}
	heading, ok := headings[anchor]
	if !ok {
		return title, fmt.Errorf("page %q: no heading with the anchor %q", href, anchor)
	}
	return title + ": " + heading, nil
}

// scanPage returns the linktitle of the page, or its first heading when
// there is none, and the text of its headings by their anchors.
func scanPage(data []byte) (string, map[string]string) {
	var title, linktitle string
	ids := map[string]string{}
	sc := bufio.NewScanner(bytes.NewReader(data))
	frontmatter, fence := false, false
	for n := 0; sc.Scan(); n++ {
		line := sc.Text()
		switch {
		case n == 0 && line == "---":
			frontmatter = true
		case frontmatter:
			if line == "---" {
				frontmatter = false
			} else if v, ok := strings.CutPrefix(line, "linktitle:"); ok {
				linktitle = strings.Trim(strings.TrimSpace(v), `"'`)
			}
		case strings.HasPrefix(line, "```"):
			fence = !fence
		case fence:
		default:
			m := headingLine.FindStringSubmatch(line)
			if m == nil {
				continue
			}
			text := m[1]
			if id := headingID.FindStringSubmatch(text); id != nil {
				text = headingID.ReplaceAllString(text, "")
				ids[id[1]] = text
			} else {
				ids[anchorize(text)] = text
			}
			if title == "" {
				title = text
			}
		}
	}
	if linktitle != "" {
		return linktitle, ids
	}
	return title, ids
}

// anchorize returns the id Hugo gives a heading with its default github
// style: letters, digits and _ in lower case, a space or - as -, the rest
// dropped. Markdown links count with their text.
func anchorize(heading string) string {
	heading = mdLink.ReplaceAllString(heading, "$1")
	var b strings.Builder
	for _, r := range heading {
		switch {
		case r == '-' || unicode.IsSpace(r):
			b.WriteRune('-')
		case r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}
