package core

import "testing"

func TestVersion(t *testing.T) {
	pass := []struct {
		layoutVersion string
		xtsVersion    string
	}{
		{"", "0.1.0"},
		{"1.2.3", "1.2.4"},
		{"1.2.3", "2.3.1"},
		{"1.2.3", "1.3.1"},
		{"0.1", "0.1.0"},
		{"0.1", "0.1.5"},
		{"0.1", "0.2.0"},
		{"0", "0.1.0"},
		{"0.0.30", "0.1.0"},
		// development builds accept every version
		{"0.1", "0.1.0-3-gabcdef1"},
		{"0.1", "dev"},
		{"0.1", "50bb014e"},
		{"0.2", "0.1.0-3-gabcdef1"},
	}

	for _, vi := range pass {
		if err := checkVersion(vi.layoutVersion, vi.xtsVersion); err != nil {
			t.Errorf("checkVersion(%q,%q) = %s", vi.layoutVersion, vi.xtsVersion, err)
		}
	}
	fail := []struct {
		layoutVersion string
		xtsVersion    string
	}{
		{"1.2.5", "1.2.4"},
		{"1.2.5", "0.2.1"},
		{"0.2", "0.1.0"},
		{"0.1.1", "0.1.0"},
		{"1", "0.1.0"},
		{"x.y", "0.1.0"},
	}

	for _, vi := range fail {
		if err := checkVersion(vi.layoutVersion, vi.xtsVersion); err == nil {
			t.Errorf("checkVersion(%q,%q) = nil, want err", vi.layoutVersion, vi.xtsVersion)
		}
	}
}
