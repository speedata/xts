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

// watchDirs returns the top level directories to watch in watch mode: the
// current directory, the directories of the layout and the data file and all
// extra directories. Files that cannot be found yet are skipped, creating them
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

// watchRecursive adds the directory and all its subdirectories to the
// watcher. Hidden directories (.git, .vscode, ...) are skipped. Style sheets,
// fonts and images usually live in subdirectories of the layout directory, so
// changes there must trigger a run too.
func watchRecursive(watcher *fsnotify.Watcher, root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			// unreadable directory: skip it, watch the rest
			if d != nil && d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !d.IsDir() {
			return nil
		}
		if path != root && strings.HasPrefix(d.Name(), ".") {
			return filepath.SkipDir
		}
		if werr := watcher.Add(path); werr != nil {
			return fmt.Errorf("cannot watch directory %s: %w", path, werr)
		}
		return nil
	})
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

// contentChanged reports whether the file content differs from the content
// seen the last time this file was checked. The check is only applied to
// files that were written while a publishing run was in progress: a Lua
// filter that rewrites its data file on every run would otherwise trigger
// runs forever. Unreadable (for example removed) files always count as
// changed.
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
//
// A file change after a run always starts a new run, even if the file content
// is identical: saving the layout again is the natural way to force a re-run
// after editing something xts does not watch. Files written while a run is in
// progress are usually written by xts itself (a Lua filter rewriting the data
// file, for example). They only start another run when their content differs
// from the previous run, so the loop converges instead of running forever.
func watchAndRun(dumpOutputFileName string, configFileRead []string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	for _, dir := range watchDirs() {
		if err = watchRecursive(watcher, dir); err != nil {
			return err
		}
	}

	done := make(chan struct{})
	running := false
	start := func() {
		running = true
		go func() {
			if err := runPublisher(dumpOutputFileName, configFileRead); err != nil {
				if terr, ok := err.(core.TypesettingError); !ok || !terr.Logged {
					fmt.Println("Error:", err)
				}
			}
			done <- struct{}{}
		}()
	}

	lastSeen := make(map[string]string)
	// files written while a run was in progress
	dirty := make(map[string]bool)
	start()

	for {
		select {
		case <-done:
			running = false
			changed := []string{}
			for name := range dirty {
				if contentChanged(lastSeen, name) {
					changed = append(changed, name)
				}
			}
			dirty = make(map[string]bool)
			if len(changed) > 0 {
				fmt.Printf("Change detected during run: %s\n", strings.Join(changed, ", "))
				start()
				continue
			}
			fmt.Println("Waiting for changes (press ctrl-c to quit)")
		case ev, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if ev.Has(fsnotify.Create) {
				if fi, serr := os.Stat(ev.Name); serr == nil && fi.IsDir() {
					// a new subdirectory: watch it, but a directory itself
					// is no input
					if werr := watchRecursive(watcher, ev.Name); werr != nil {
						fmt.Println("Watch error:", werr)
					}
					continue
				}
			}
			if !triggersRun(ev, dumpOutputFileName) {
				continue
			}
			if running {
				dirty[ev.Name] = true
				continue
			}
			// let the editor finish writing before the file gets read
			discardEvents(watcher, 200*time.Millisecond)
			fmt.Printf("Change detected: %s\n", ev.Name)
			start()
		case werr, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Println("Watch error:", werr)
		}
	}
}
