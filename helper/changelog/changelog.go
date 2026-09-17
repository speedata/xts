// Package changelog reads doc/changelog.xml and renders it as Markdown for
// the manual and the GitHub release notes.
package changelog

import (
	"encoding/xml"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/speedata/xts/helper/config"
)

// Repo is the GitHub repository the commit and issue links point to.
const Repo = "https://github.com/speedata/xts"

// Text is the description of an entry: a one sentence summary and the
// details as inner XML (may contain <tt> markup).
type Text struct {
	Summary string `xml:"summary,attr"`
	Text    string `xml:",innerxml"`
}

// Entry is one user visible change.
type Entry struct {
	Version string `xml:"version,attr"`
	Date    string `xml:"date,attr"`
	SHA1    string `xml:"sha1,attr"`
	En      Text   `xml:"en"`
}

// Chapter groups the entries of a minor version.
type Chapter struct {
	Version string  `xml:"version,attr"`
	Entries []Entry `xml:"entry"`
}

// Changelog is the parsed changelog.xml.
type Changelog struct {
	Chapter []Chapter `xml:"chapter"`
}

// Release is the list of entries of one version, in file order.
type Release struct {
	Version string
	Date    time.Time
	Entries []Entry
}

// Parse reads and parses a changelog.xml file.
func Parse(filename string) (*Changelog, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cl := &Changelog{}
	if err = xml.Unmarshal(data, cl); err != nil {
		return nil, err
	}
	return cl, nil
}

// Read parses basedir/doc/changelog.xml.
func Read(cfg *config.Config) (*Changelog, error) {
	return Parse(filepath.Join(cfg.Basedir(), "doc", "changelog.xml"))
}

// Releases groups the entries by version in file order. The date of a
// release is the latest entry date.
func (cl *Changelog) Releases() []Release {
	var rr []Release
	idx := map[string]int{}
	for _, chap := range cl.Chapter {
		for _, e := range chap.Entries {
			i, ok := idx[e.Version]
			if !ok {
				i = len(rr)
				idx[e.Version] = i
				rr = append(rr, Release{Version: e.Version})
			}
			d, _ := time.Parse("2006-01-02", e.Date)
			if d.After(rr[i].Date) {
				rr[i].Date = d
			}
			rr[i].Entries = append(rr[i].Entries, e)
		}
	}
	return rr
}

// Release returns the entries of the given version (with or without a
// leading v) or nil.
func (cl *Changelog) Release(version string) *Release {
	version = strings.TrimPrefix(version, "v")
	for _, r := range cl.Releases() {
		if r.Version == version {
			r := r
			return &r
		}
	}
	return nil
}

var sha1Re = regexp.MustCompile(`^[0-9a-f]{7,40}$`)

// Validate checks the structural rules of the file: every entry has a
// version that belongs to its chapter, a parseable date, a summary and
// well formed commit ids, and the versions are in descending order.
func (cl *Changelog) Validate() error {
	var prev []int
	for _, chap := range cl.Chapter {
		for _, e := range chap.Entries {
			where := fmt.Sprintf("entry %s (%s)", e.Version, e.Date)
			v, err := parseVersion(e.Version)
			if err != nil {
				return fmt.Errorf("%s: %w", where, err)
			}
			if !strings.HasPrefix(e.Version, chap.Version+".") {
				return fmt.Errorf("%s: not in chapter %s", where, chap.Version)
			}
			if prev != nil && compareVersion(v, prev) > 0 {
				return fmt.Errorf("%s: versions must be in descending order", where)
			}
			prev = v
			if _, err := time.Parse("2006-01-02", e.Date); err != nil {
				return fmt.Errorf("%s: bad date: %w", where, err)
			}
			if strings.TrimSpace(e.En.Summary) == "" {
				return fmt.Errorf("%s: missing summary", where)
			}
			if strings.TrimSpace(e.En.Text) == "" {
				return fmt.Errorf("%s: missing text", where)
			}
			for _, sha := range splitSHA1(e.SHA1) {
				if !sha1Re.MatchString(sha) {
					return fmt.Errorf("%s: bad sha1 %q", where, sha)
				}
			}
		}
	}
	return nil
}

// Check validates the file and verifies that it is ready for a release of
// the given version: the version has entries and no entry names a newer
// version.
func (cl *Changelog) Check(version string) error {
	if err := cl.Validate(); err != nil {
		return err
	}
	version = strings.TrimPrefix(version, "v")
	want, err := parseVersion(version)
	if err != nil {
		return err
	}
	if cl.Release(version) == nil {
		return fmt.Errorf("no changelog entries for version %s", version)
	}
	for _, r := range cl.Releases() {
		v, _ := parseVersion(r.Version)
		if compareVersion(v, want) > 0 {
			return fmt.Errorf("changelog has entries for %s, newer than the release %s", r.Version, version)
		}
	}
	return nil
}

func parseVersion(s string) ([]int, error) {
	parts := strings.Split(strings.TrimPrefix(s, "v"), ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("version %q must have the form major.minor.patch", s)
	}
	v := make([]int, 3)
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("version %q: %w", s, err)
		}
		v[i] = n
	}
	return v, nil
}

func compareVersion(a, b []int) int {
	for i := range a {
		if a[i] != b[i] {
			return a[i] - b[i]
		}
	}
	return 0
}

func splitSHA1(s string) []string {
	var out []string
	for _, sha := range strings.Split(s, ",") {
		if sha = strings.TrimSpace(sha); sha != "" {
			out = append(out, sha)
		}
	}
	return out
}

var (
	ghIssueRe = regexp.MustCompile(`#(\d+)`)
	ttBlockRe = regexp.MustCompile(`<tt>.*?</tt>`)
	ttRepl    = strings.NewReplacer(`<tt>`, "`", `</tt>`, "`")
	mdEscRepl = strings.NewReplacer(`\`, `\\`, `*`, `\*`, `_`, `\_`, `[`, `\[`, `]`, `\]`, `<`, `\<`, `>`, `\>`)
	spaceRe   = regexp.MustCompile(`\s+`)
)

// markdown turns the inner XML of a summary or text into Markdown: entity
// references are decoded, Markdown characters outside <tt> are escaped,
// <tt> becomes a code span and issue numbers become links.
func markdown(s string) string {
	s = spaceRe.ReplaceAllString(strings.TrimSpace(s), " ")
	s = html.UnescapeString(s)
	var b strings.Builder
	last := 0
	for _, loc := range ttBlockRe.FindAllStringIndex(s, -1) {
		b.WriteString(mdEscRepl.Replace(s[last:loc[0]]))
		b.WriteString(s[loc[0]:loc[1]])
		last = loc[1]
	}
	b.WriteString(mdEscRepl.Replace(s[last:]))
	s = ghIssueRe.ReplaceAllString(b.String(), `[#$1](`+Repo+`/issues/$1)`)
	return ttRepl.Replace(s)
}

// Markdown renders the entries of a release as a bullet list. With
// commitLinks, every entry gets a link to its commit(s).
func (r *Release) Markdown(commitLinks bool) string {
	var b strings.Builder
	for _, e := range r.Entries {
		fmt.Fprintf(&b, "- **%s**", markdown(e.En.Summary))
		if commitLinks {
			for _, sha := range splitSHA1(e.SHA1) {
				fmt.Fprintf(&b, " [↗](%s/commit/%s)", Repo, sha)
			}
		}
		fmt.Fprintf(&b, "<br>\n  %s\n", markdown(e.En.Text))
	}
	return b.String()
}

// ReleaseNotes renders the notes for the GitHub release of a version.
func (cl *Changelog) ReleaseNotes(version string) (string, error) {
	r := cl.Release(version)
	if r == nil {
		return "", fmt.Errorf("no changelog entries for version %s", version)
	}
	return r.Markdown(false), nil
}

// WriteManualPage writes the changelog section of the Hugo manual.
func (cl *Changelog) WriteManualPage(cfg *config.Config) error {
	dir := filepath.Join(cfg.Basedir(), "doc", "manual", "content", "changelog")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	var b strings.Builder
	b.WriteString("---\ntitle: Changelog\nlinktitle: Changelog\ntype: docs\nweight: 60\n---\n\n")
	b.WriteString("<!-- Generated by rake doc from doc/changelog.xml, do not edit. -->\n\n")
	b.WriteString("All user visible changes of XTS, newest first. The version at the top may not be released yet.\n\n")
	for _, r := range cl.Releases() {
		fmt.Fprintf(&b, "## %s (%s)\n\n%s\n", r.Version, r.Date.Format("2006-01-02"), r.Markdown(true))
	}
	return os.WriteFile(filepath.Join(dir, "_index.md"), []byte(b.String()), 0o644)
}
