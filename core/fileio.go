package core

import (
	"io/fs"
	"os"
	"path/filepath"
)

var filelist map[string]string

func init() {
	filelist = make(map[string]string)

	root, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filepath.WalkDir(filepath.Join(root, "fonts"), dirWalker)
	filepath.WalkDir(filepath.Join(root, "hyphenationpatterns"), dirWalker)
}

func dirWalker(path string, d fs.DirEntry, err error) error {
	if d.Type().IsRegular() {
		filelist[filepath.Base(path)] = path
	}
	return nil
}

func findFile(filename string) (string, error) {
	if fn, ok := filelist[filename]; ok {
		logger.Debugf("File lookup %q -> %q", filename, fn)
		return fn, nil
	}

	logger.Debugf("File lookup %q not found", filename)
	return "", os.ErrNotExist
}
