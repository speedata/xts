package core

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var filelist = make(map[string]string)

// AddDir recursively adds a directory to the file list
func AddDir(dirname string) error {
	slog.Debug("Add directory to recursive file list", "dir", dirname)
	return filepath.WalkDir(dirname, dirWalker)
}

// InitDirs starts indexing the files.
func InitDirs(basedir string) error {
	var err error
	for _, dir := range []string{"img"} {
		dir = filepath.Join(basedir, dir)
		if err = AddDir(dir); err != nil {
			return err
		}
	}
	return nil
}

func dirWalker(path string, d fs.DirEntry, err error) error {
	if d == nil {
		return fmt.Errorf("%w %q", os.ErrNotExist, path)
	}
	if d.Type().IsRegular() {
		filelist[filepath.Base(path)] = path
	}
	return nil
}

// urldownloader downloads the given URI to a file. No caching is performed.
func urldownloader(uri string) (string, error) {
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	destfile := parsedURL.Hostname() + parsedURL.Path
	if parsedURL.RawQuery != "" {
		destfile += "?" + parsedURL.RawQuery
	}

	hashedFilename := fmt.Sprintf("%x", md5.Sum([]byte(destfile)))
	tmpdir, err := os.MkdirTemp("", "xtsimages")
	if err != nil {
		return "", err
	}
	w, err := os.Create(filepath.Join(tmpdir, hashedFilename))
	if err != nil {
		return "", err
	}
	resp, err := http.Get(uri)
	if err != nil {
		return "", err
	}
	if _, err = io.Copy(w, resp.Body); err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return w.Name(), nil
}

// FindFile returns the full path to the file name.
func FindFile(filename string) (string, error) {
	if fn, ok := filelist[filename]; ok {
		slog.Debug("File lookup", "src", filename, "found", fn)
		return fn, nil
	}
	if strings.HasPrefix(filename, "https://") || strings.HasPrefix(filename, "http://") {
		fn, err := urldownloader(filename)
		if err != nil {
			return "", err
		}
		slog.Info("Write URL to file", "url", filename, "file", fn)
		return fn, nil
	}
	if _, err := os.Stat(filename); err == nil {
		var fn string
		fn, err = filepath.Abs(filename)
		if err != nil {
			return "", err
		}
		slog.Debug("File lookup", "src", filename, "found", fn)
		return fn, nil
	}
	slog.Debug("File lookup (not found)", "src", filename)
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
