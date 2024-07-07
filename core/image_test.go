package core

import (
	"testing"

	"github.com/boxesandglue/boxesandglue/backend/bag"
	bagimage "github.com/boxesandglue/boxesandglue/backend/image"
)

func TestImageSize(t *testing.T) {
	img := &bagimage.Image{}
	fiftypt := 50 * bag.Factor
	seventypt := 70 * bag.Factor
	hundretpt := 100 * bag.Factor
	twohundretpt := 200 * bag.Factor
	img.Width = hundretpt
	img.Height = seventypt

	testdata := []struct {
		wdWant  bag.ScaledPoint
		htWant  bag.ScaledPoint
		wd      *bag.ScaledPoint
		ht      *bag.ScaledPoint
		minwd   *bag.ScaledPoint
		maxwd   *bag.ScaledPoint
		minht   *bag.ScaledPoint
		maxht   *bag.ScaledPoint
		stretch bool
	}{
		{fiftypt, fiftypt, &fiftypt, &fiftypt, nil, nil, nil, nil, false},
		{fiftypt, fiftypt, &fiftypt, &fiftypt, nil, nil, nil, nil, true},
		{twohundretpt, twohundretpt, &twohundretpt, &twohundretpt, nil, nil, nil, nil, true},
		{twohundretpt, 2 * seventypt, nil, nil, nil, &twohundretpt, nil, &twohundretpt, true},
		{fiftypt, seventypt / 2, nil, nil, nil, &fiftypt, nil, &fiftypt, true},
		{hundretpt, seventypt, nil, nil, &hundretpt, nil, nil, nil, true},
	}
	for _, tc := range testdata {

		wd, ht := calculateImageSize(img, tc.wd, tc.ht, tc.minwd, tc.maxwd, tc.minht, tc.maxht, tc.stretch)
		if wd != tc.wdWant || ht != tc.htWant {
			t.Errorf("calculateImageSize: got (%s,%s), want (%s,%s)", wd, ht, tc.wdWant, tc.htWant)
		}

	}
}
