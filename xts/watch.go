package main

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/speedata/xts/core"
)

// watchDirs returns the directories to watch in watch mode: the current
// directory, the directories of the layout and the data file and all extra
// directories. Files that cannot be found yet are skipped, creating them
// later in a watched directory triggers a run.
func watchDirs() []string {
	seen := make(map[string]bool)
	dirs := []string{}
	add := func(dir string) {
		abs, err := filepath.Abs(dir)
		if err != nil {
			return
		}
		if !seen[abs] {
			seen[abs] = true
			dirs = append(dirs, abs)
		}
	}
	add(".")
	if p, err := core.FindFile(configuration.Layout); err == nil {
		add(filepath.Dir(p))
	}
	if !configuration.Dummy {
		if p, err := core.FindFile(configuration.Data); err == nil {
			add(filepath.Dir(p))
		}
	}
	for _, d := range configuration.ExtraDir {
		add(d)
	}
	return dirs
}

// isGeneratedFile reports whether the file is written by xts itself during a
// publishing run (PDF, protocol, aux files, XML dump). Those must not trigger
// a new run, otherwise the watch loop would run forever.
func isGeneratedFile(filename string, dumpOutputFileName string) bool {
	base := filepath.Base(filename)
	jobname := configuration.Jobname
	if base == jobname+".pdf" {
		return true
	}
	if strings.HasPrefix(base, jobname+"-") && strings.HasSuffix(base, ".xml") {
		return true
	}
	if strings.HasSuffix(base, ".protocol") {
		return true
	}
	if dumpOutputFileName != "" && base == filepath.Base(dumpOutputFileName) {
		return true
	}
	return false
}

// isEditorArtifact reports whether the file is a temporary or backup file
// that editors write while editing (vim swap files, backup~ files, hidden
// files).
func isEditorArtifact(filename string) bool {
	base := filepath.Base(filename)
	if strings.HasPrefix(base, ".") {
		return true
	}
	switch {
	case strings.HasSuffix(base, "~"),
		strings.HasSuffix(base, ".swp"),
		strings.HasSuffix(base, ".swx"),
		strings.HasSuffix(base, ".tmp"),
		base == "4913": // vim checks with this file if it can write to the directory
		return true
	}
	return false
}

// triggersRun reports whether the event might start a new publishing run.
func triggersRun(ev fsnotify.Event, dumpOutputFileName string) bool {
	if !ev.Has(fsnotify.Write) && !ev.Has(fsnotify.Create) && !ev.Has(fsnotify.Rename) && !ev.Has(fsnotify.Remove) {
		return false
	}
	if isGeneratedFile(ev.Name, dumpOutputFileName) || isEditorArtifact(ev.Name) {
		return false
	}
	return true
}

// contentChanged reports whether the file content differs from the content at
// the last publishing run that this file triggered. Without this check a Lua
// filter that rewrites its data file on every run would trigger runs forever.
// Unreadable (for example removed) files always count as changed.
func contentChanged(lastSeen map[string]string, filename string) bool {
	data, err := os.ReadFile(filename)
	if err != nil {
		delete(lastSeen, filename)
		return true
	}
	sum := fmt.Sprintf("%x", md5.Sum(data))
	if lastSeen[filename] == sum {
		return false
	}
	lastSeen[filename] = sum
	return true
}

// discardEvents reads and drops all watcher events for the given duration.
// Editors emit several events for a single file save, this collects the whole
// burst before the run starts.
func discardEvents(watcher *fsnotify.Watcher, d time.Duration) {
	timer := time.NewTimer(d)
	defer timer.Stop()
	for {
		select {
		case <-watcher.Events:
		case <-watcher.Errors:
		case <-timer.C:
			return
		}
	}
}

// watchAndRun runs the publishing process once and re-runs it whenever an
// input file changes. Errors from a publishing run do not stop the loop, the
// watcher waits for the next change instead.
func watchAndRun(dumpOutputFileName string, configFileRead []string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	dirs := watchDirs()
	for _, dir := range dirs {
		if err = watcher.Add(dir); err != nil {
			return fmt.Errorf("cannot watch directory %s: %w", dir, err)
		}
	}

	runOnce := func() {
		if err := runPublisher(dumpOutputFileName, configFileRead); err != nil {
			if terr, ok := err.(core.TypesettingError); !ok || !terr.Logged {
				fmt.Println("Error:", err)
			}
		}
		fmt.Println("Waiting for changes (press ctrl-c to quit)")
	}

	lastSeen := make(map[string]string)
	runOnce()

	for {
		select {
		case ev, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if !triggersRun(ev, dumpOutputFileName) {
				continue
			}
			// let the editor finish writing before the file gets read
			discardEvents(watcher, 200*time.Millisecond)
			if !contentChanged(lastSeen, ev.Name) {
				continue
			}
			fmt.Printf("Change detected: %s\n", ev.Name)
			runOnce()
		case werr, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Println("Watch error:", werr)
		}
	}
}
