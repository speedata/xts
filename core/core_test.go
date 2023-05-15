package core

import "testing"

func TestVersion(t *testing.T) {
	data := []struct {
		layoutVersion string
		xtsVersion    string
	}{
		{"1.2.3", "1.2.4"},
		{"1.2.3", "2.3.1"},
		{"1.2.3", "1.3.1"},
	}

	for _, vi := range data {
		if err := checkVersion(vi.layoutVersion, vi.xtsVersion); err != nil {
			t.Errorf("checkVersion(%s,%s) = %s", vi.layoutVersion, vi.xtsVersion, err)
		}
	}
	data = []struct {
		layoutVersion string
		xtsVersion    string
	}{
		{"1.2.5", "1.2.4"},
		{"1.2.5", "0.2.1"},
	}

	for _, vi := range data {
		if err := checkVersion(vi.layoutVersion, vi.xtsVersion); err == nil {
			t.Errorf("checkVersion(%s,%s) = nil, want err", vi.layoutVersion, vi.xtsVersion)
		}
	}
}
