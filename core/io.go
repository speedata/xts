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

// statCache holds the absolute paths of the files FindFile found by looking
// at the file system during the current run. A layout that places the same
// image on hundreds of pages would otherwise stat the file (and resolve the
// working directory) for every placement. RunXTS clears it, so a file that
// disappears between two runs in watch mode is noticed.
var statCache = make(map[string]string)

// resetStatCache forgets the files found on the file system.
func resetStatCache() {
	statCache = make(map[string]string)
}

// AddDir recursively adds a directory to the file list
func AddDir(dirname string) error {
	slog.Debug("Add directory to recursive file list", "dir", dirname)
	return filepath.WalkDir(dirname, dirWalker)
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
	if fn, ok := statCache[filename]; ok {
		slog.Debug("File lookup", "src", filename, "found", fn)
		return fn, nil
	}
	if _, err := os.Stat(filename); err == nil {
		var fn string
		fn, err = filepath.Abs(filename)
		if err != nil {
			return "", err
		}
		slog.Debug("File lookup", "src", filename, "found", fn)
		statCache[filename] = fn
		return fn, nil
	}
	slog.Debug("File lookup (not found)", "src", filename)
	return "", fmt.Errorf("%w: %s", os.ErrNotExist, filename)
}

func fileexists(fn string) bool {
	if _, ok := filelist[fn]; ok {
		return true
	}
	_, err := os.Stat(fn)
	return err == nil
}
