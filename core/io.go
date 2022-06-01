package core

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/speedata/boxesandglue/backend/bag"
)

var filelist = make(map[string]string)

// AddDir recursively adds a directory to the file list
func AddDir(dirname string) error {
	bag.Logger.Debugf("Add directory %q to recursive file list", dirname)
	filepath.WalkDir(dirname, dirWalker)
	return nil
}

// InitDirs starts indexing the files.
func InitDirs(basedir string) error {
	var err error
	for _, dir := range []string{"fonts", "img"} {
		dir = filepath.Join(basedir, dir)
		if err = AddDir(dir); err != nil {
			return err
		}
	}
	return nil
}

func dirWalker(path string, d fs.DirEntry, err error) error {
	if d != nil && d.Type().IsRegular() {
		filelist[filepath.Base(path)] = path
	}
	return nil
}

// FindFile returns the full path to the file name.
func FindFile(filename string) (string, error) {
	if fn, ok := filelist[filename]; ok {
		bag.Logger.Debugf("File lookup %q -> %q", filename, fn)
		return fn, nil
	}
	if _, err := os.Stat(filename); err == nil {
		var fn string
		fn, err = filepath.Abs(filename)
		if err != nil {
			return "", err
		}
		bag.Logger.Debugf("File lookup %q -> %q", filename, fn)
		return fn, nil
	}

	bag.Logger.Debugf("File lookup %q not found", filename)
	return "", fmt.Errorf("%w: %s", os.ErrNotExist, filename)
}

func isFontFile(filename string) bool {
	l := strings.ToLower(filename)
	return strings.HasSuffix(l, ".ttf") || strings.HasSuffix(l, ".otf")
}

// FindFontFiles returns a list of all font files (otf,ttf)
func FindFontFiles() []string {
	var ret []string
	for _, fn := range filelist {
		if isFontFile(fn) {
			ret = append(ret, fn)
		}
	}
	return ret
}

func fileexists(fn string) bool {
	_, ok := filelist[fn]
	return ok
}
