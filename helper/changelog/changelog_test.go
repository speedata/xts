package changelog

import (
	"strings"
	"testing"
)

// The checked-in changelog must always pass the structural rules, so a
// broken entry fails go test before it reaches a release.
func TestCheckedInChangelog(t *testing.T) {
	cl, err := Parse("../../doc/changelog.xml")
	if err != nil {
		t.Fatal(err)
	}
	if err := cl.Validate(); err != nil {
		t.Fatal(err)
	}
	if len(cl.Releases()) == 0 {
		t.Fatal("no releases")
	}
}

func TestCheck(t *testing.T) {
	cl := &Changelog{Chapter: []Chapter{{Version: "0.0", Entries: []Entry{
		{Version: "0.0.30", Date: "2026-09-17", SHA1: "cec2eff", En: Text{Summary: "s", Text: "t"}},
		{Version: "0.0.29", Date: "2026-09-14", En: Text{Summary: "s", Text: "t"}},
	}}}}
	if err := cl.Check("v0.0.30"); err != nil {
		t.Errorf("v0.0.30: %v", err)
	}
	if err := cl.Check("0.0.29"); err == nil || !strings.Contains(err.Error(), "newer") {
		t.Errorf("0.0.29 should report newer entries, got %v", err)
	}
	if err := cl.Check("0.0.31"); err == nil || !strings.Contains(err.Error(), "no changelog entries") {
		t.Errorf("0.0.31 should report missing entries, got %v", err)
	}
}

func TestValidateOrder(t *testing.T) {
	cl := &Changelog{Chapter: []Chapter{{Version: "0.0", Entries: []Entry{
		{Version: "0.0.29", Date: "2026-09-14", En: Text{Summary: "s", Text: "t"}},
		{Version: "0.0.30", Date: "2026-09-17", En: Text{Summary: "s", Text: "t"}},
	}}}}
	if err := cl.Validate(); err == nil {
		t.Error("ascending versions must be rejected")
	}
}

func TestMarkdown(t *testing.T) {
	r := &Release{Entries: []Entry{{
		SHA1: "abc1234,def5678",
		En:   Text{Summary: "New <tt>&lt;Foo&gt;</tt> command (#12).", Text: "Uses <tt>a_b</tt> and\n   <tt>*</tt>, see *docs*."},
	}}}
	got := r.Markdown(true)
	for _, want := range []string{
		"- **New `<Foo>` command ([#12](" + Repo + "/issues/12)).**",
		"[↗](" + Repo + "/commit/abc1234) [↗](" + Repo + "/commit/def5678)<br>",
		"Uses `a_b` and `*`, see \\*docs\\*.",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in\n%s", want, got)
		}
	}
}
