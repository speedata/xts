package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"html/template"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/gammazero/workerpool"
)

type compareStatus struct {
	Path     string
	Badpages []int
	Delta    float64
}

type byDelta []compareStatus

func (a byDelta) Len() int           { return len(a) }
func (a byDelta) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byDelta) Less(i, j int) bool { return a[i].Delta > a[j].Delta }

var (
	finished          chan bool
	exeSuffix         string
	cs                []compareStatus
	allPages          []compareStatus
	mutex             *sync.Mutex
	wp                *workerpool.WorkerPool
	referencefilename string
)

func init() {
	finished = make(chan bool)
	mutex = &sync.Mutex{}
}

func fileExists(filename string) bool {
	fi, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !fi.IsDir()
}

// doCompare starts comparing the files in the current directory and its
// subdirectory. This is the function to be called (first).
func doCompare(absdir string, withHTML bool, referencefn string) {
	switch runtime.GOOS {
	case "windows":
		exeSuffix = ".exe"
	default:
		exeSuffix = ""
	}
	wp = workerpool.New(runtime.NumCPU())
	referencefilename = referencefn
	statuschan := make(chan []compareStatus, 0)
	compareFunc := mkCompare(statuschan)
	filepath.Walk(absdir, compareFunc)
	go getCompareStatus(statuschan)
	wp.StopWait()

	finished <- true
	if withHTML {
		mkWebPage(!configuration.Verbose)
	}
}

func compareTwoPages(sourcefile, referencefile, dummyfile, path string) float64 {
	// More complicated than the trivial case because I need the different exit statuses.
	// See http://stackoverflow.com/a/10385867
	if !fileExists(filepath.Join(path, sourcefile)) || !fileExists(filepath.Join(path, referencefile)) {
		return 99.0
	}

	cmd := exec.Command("compare"+exeSuffix, "-metric", "mae", sourcefile, referencefile, dummyfile)
	cmd.Dir = path
	// err == 1 looks like an indicator that the comparison is OK but some diffs in the images
	// err == 2 seems to be a fatal error
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Println("Do you have imagemagick installed?")
		log.Fatalf("cmd.Start: %v", err)
	}

	r := bufio.NewReader(stderr)
	line, _ := r.ReadBytes('\n')

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on Mac and hopefully on Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				if status.ExitStatus() == 1 {
					// comparison ok with differences
					delta, nerr := strconv.ParseFloat(strings.Split(string(line), " ")[0], 32)
					if nerr != nil {
						log.Fatal(nerr)
					}
					return delta
				}
				log.Fatal(err)
			}
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}
	return 0.0
}

func newer(pdf, png string) bool {
	pngFi, err := os.Stat(png)
	if err != nil {
		return true
	}
	pdfFi, err := os.Stat(pdf)
	if err != nil {
		panic(fmt.Sprintf("Source %s does not exist!", pdf))
	}
	return pngFi.ModTime().Before(pdfFi.ModTime())
}

func calculateHash(filename string) []byte {
	fh, err := os.Open(filename)
	if err != nil {
		panic(err.Error())
	}
	defer fh.Close()

	h := sha256.New()
	_, err = io.Copy(h, fh)
	if err != nil {
		panic(err.Error())
	}
	return h.Sum(nil)
}

func runComparison(path string, statuschan chan []compareStatus) {
	cs := compareStatus{}
	allPages := compareStatus{}
	cs.Path = path
	allPages.Path = path
	allPages.Badpages = append(allPages.Badpages, 0)
	var err error
	if configuration.Verbose {
		fmt.Println(path)
	}
	cmd := exec.Command("xts"+exeSuffix, "--suppressinfo")
	cmd.Dir = path
	err = cmd.Run()
	if err != nil {
		log.Println(path)
		log.Fatal("Error running command 'xts': ", err)
	}

	p := calculateHash(filepath.Join(path, "xts.pdf"))
	r := calculateHash(filepath.Join(path, fmt.Sprintf("%s.pdf", referencefilename)))
	if bytes.Equal(p, r) {
		if configuration.Verbose {
			fmt.Printf("Files in %q have the same checksum\n", path)
		}
		cs.Delta = 0
		allPages.Delta = 0
		statuschan <- []compareStatus{cs, allPages}
		return
	}
	if configuration.Verbose {
		fmt.Printf("Run convert for %q\n", path)
	}
	sourceFiles, err := filepath.Glob(filepath.Join(path, "source-*.png"))
	if err != nil {
		log.Fatal(err)
	}

	// Let's remove the old source files, otherwise
	// the number of pages (below) might
	// be incorrect which results in a fatal
	// error
	for _, name := range sourceFiles {
		err = os.Remove(name)
		if err != nil {
			log.Println(err)
		}
	}

	cmd = exec.Command("convert"+exeSuffix, "-density", "150", "-trim", "xts.pdf", "source-%02d.png")
	cmd.Dir = path
	cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// convert the reference pdf to png for later comparisons
	// we only do that when the pdf is newer than the png files
	// (that is: the pdf has been updated)
	if newer(filepath.Join(path, fmt.Sprintf("%s.pdf", referencefilename)), filepath.Join(path, "reference-00.png")) {
		cmd := exec.Command("convert"+exeSuffix, "-density", "150", "-trim", fmt.Sprintf("%s.pdf", referencefilename), referencefilename+"-%02d.png")
		cmd.Dir = path
		err = cmd.Run()
		if err != nil {
			log.Fatal("Error converting reference. Do you have ghostscript installed?", err)
		}
	}

	sourceFiles, err = filepath.Glob(filepath.Join(path, "source-*.png"))
	if err != nil {
		log.Fatal("No source files found. ", err)
	}

	for i := 0; i < len(sourceFiles); i++ {
		sourceFile := fmt.Sprintf("source-%02d.png", i)
		referenceFile := fmt.Sprintf("%s-%02d.png", referencefilename, i)
		dummyFile := fmt.Sprintf("pagediff-%02d.png", i)
		if delta := compareTwoPages(sourceFile, referenceFile, dummyFile, path); delta > 0 {
			cs.Delta = math.Max(cs.Delta, delta)
			allPages.Delta = cs.Delta
			if delta > 0.3 {
				if i > 0 {
					allPages.Badpages = append(allPages.Badpages, i)
				}
				cs.Badpages = append(cs.Badpages, i)
			}
		}
	}

	statuschan <- []compareStatus{cs, allPages}
}

func mkWebPage(onlyErrorPages bool) error {
	if onlyErrorPages && len(cs) == 0 {
		return nil
	}
	var pages []compareStatus
	if onlyErrorPages {
		pages = cs
	} else {
		pages = allPages
	}

	sort.Sort(byDelta(pages))

	tmpl := `<!DOCTYPE html>
<html>
<head>
	<title>speedata compare result</title>
	<style type="text/css">
		img { height: 150px ; border: 1px solid black; max-width: 75%; width: auto; height: auto; }
		tr.img td	{ border-bottom: 1px solid black; }
		tr  {vertical-align: top;}
	</style>
</head>
<body>
	<table>
	{{ range .CompareStatus -}}
	{{ $path := .Path}}
	<tr>
		<td colspan="1">{{ .Path }} ({{ .Delta | printf "%.3f" }})</td>
	</tr>
	<tr>
		<td>
		{{range .Badpages}}{{.}}: <a href="{{ $path}}/{{. | printf $.ImageToShow }}"><img src="{{ $path}}/{{. | printf $.ImageToShow }}" ></a>{{end}}
		</td>
	</tr>
	{{- end }}
</table>
</body>
</html>
`

	var buf bytes.Buffer

	t := template.Must(template.New("html").Parse(tmpl))
	data := struct {
		CompareStatus []compareStatus
		ImageToShow   string
	}{
		CompareStatus: pages,
	}
	if onlyErrorPages {
		data.ImageToShow = "pagediff-%.2d.png"
	} else {
		data.ImageToShow = "source-%.2d.png"
	}
	err := t.Execute(&buf, data)
	if err != nil {
		return err
	}
	outfile := "compare-report.html"
	f, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = buf.WriteTo(f)
	if err != nil {
		return err
	}
	fmt.Println("Output written to", outfile)
	return nil
}

func getCompareStatus(statuschan chan []compareStatus) {
	for {
		select {
		case st := <-statuschan:
			allPages = append(allPages, st[1])
			if len(st[0].Badpages) > 0 {
				mutex.Lock()
				cs = append(cs, st[0])
				mutex.Unlock()
				fmt.Println("---------------------------")
				fmt.Println("Finished with comparison in")
				fmt.Println(st[0].Path)
				fmt.Println("Comparison failed. Bad pages are:", st[0].Badpages)
				fmt.Println("Max delta is", fmt.Sprintf("%.2f", st[0].Delta))
			}
		case <-finished:
			// now that we have read from the channel, we are all done
		}
	}
}

// Return a filepath.WalkFunc that looks into a directory, runs convert to generate the PNG files from the PDF and
// compares the two resulting files. The function puts the result into the channel compareStatus.
func mkCompare(statuschan chan []compareStatus) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if info == nil || !info.IsDir() {
			return nil
		}
		pdffn := fmt.Sprintf("%s.pdf", referencefilename)
		if _, err := os.Stat(filepath.Join(path, pdffn)); err == nil {
			wp.Submit(func() { runComparison(path, statuschan) })
		} else if _, err := os.Stat(filepath.Join(path, "layout.xml")); err == nil {
			fmt.Println("Warning: directory", path, "has layout.xml but not", pdffn)
		}
		return nil
	}
}
