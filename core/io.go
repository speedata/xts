package core

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/speedata/boxesandglue/backend/bag"
)

var filelist = make(map[string]string)

// InitDirs starts indexing the files.
func InitDirs() error {
	root, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, dir := range []string{"fonts", "img"} {
		dir = filepath.Join(root, dir)
		bag.Logger.Debugf("Add directory %q to recursive file list", dir)
		filepath.WalkDir(dir, dirWalker)
	}
	return nil
}

func dirWalker(path string, d fs.DirEntry, err error) error {
	if d.Type().IsRegular() {
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

func fileexists(fn string) bool {
	_, ok := filelist[fn]
	return ok
}
