package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/pkg/css"
	"github.com/mateussouzaweb/compactor/pkg/generic"
	"github.com/mateussouzaweb/compactor/pkg/gif"
	"github.com/mateussouzaweb/compactor/pkg/html"
	"github.com/mateussouzaweb/compactor/pkg/javascript"
	"github.com/mateussouzaweb/compactor/pkg/jpeg"
	"github.com/mateussouzaweb/compactor/pkg/json"
	"github.com/mateussouzaweb/compactor/pkg/png"
	"github.com/mateussouzaweb/compactor/pkg/sass"
	"github.com/mateussouzaweb/compactor/pkg/svg"
	"github.com/mateussouzaweb/compactor/pkg/typescript"
	"github.com/mateussouzaweb/compactor/pkg/webp"
	"github.com/mateussouzaweb/compactor/pkg/xml"
	"github.com/radovskyb/watcher"
)

var (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Purple  = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	White   = "\033[97m"
	Info    = Cyan
	Warn    = Yellow
	Fatal   = Red
	Success = Green
)

func print(color string, format string, args ...interface{}) {
	fmt.Printf(color+format+Reset, args...)
}

func main() {

	// Command line flags
	version := flag.Bool("version", false, "Print program version")
	source := flag.String("source", "src/", "Path of project source files")
	destination := flag.String("destination", "dist/", "Path to the destination folder")

	watch := flag.Bool("watch", false, "Enable live watch compilation")
	minify := flag.Bool("minify", true, "Minify code on compilation")
	sourceMap := flag.Bool("source-map", true, "Include source map on compilation")
	compress := flag.Bool("compress", true, "Compress images to reduce size")
	progressive := flag.Bool("progressive", true, "Generate progressive new images formats")

	var include []string
	flag.Func("include", "Include matching files on the pattern", func(value string) error {
		include = append(include, value)
		return nil
	})

	var exclude []string
	flag.Func("exclude", "Exclude matching files on glob pattern", func(value string) error {
		exclude = append(exclude, value)
		return nil
	})

	var maps []string
	flag.Func("maps", "Maps matching files on glob pattern to destination", func(value string) error {
		maps = append(maps, value)
		return nil
	})

	// Parse values
	flag.Parse()

	if *version {
		print("", "Compactor version 0.0.1\n")
		return
	}

	// Add parsers
	compactor.Add("*", generic.Processor)
	compactor.Add("sass", sass.Processor)
	compactor.Add("scss", sass.Processor)
	compactor.Add("css", css.Processor)
	compactor.Add("ts", typescript.Processor)
	compactor.Add("tsx", typescript.Processor)
	compactor.Add("js", javascript.Processor)
	compactor.Add("json", json.Processor)
	compactor.Add("xml", xml.Processor)
	compactor.Add("html", html.Processor)
	compactor.Add("svg", svg.Processor)
	compactor.Add("gif", gif.Processor)
	compactor.Add("jpeg", jpeg.Processor)
	compactor.Add("jpg", jpeg.Processor)
	compactor.Add("png", png.Processor)
	compactor.Add("webp", webp.Processor)

	// Find real path
	rootSource, _ := filepath.Abs(*source)
	rootDestination, _ := filepath.Abs(*destination)

	// Print information
	print(Purple, ":::| COMPACTOR |:::\n")
	print(Info, "[INFO] Files source folder is %s\n", rootSource)
	print(Info, "[INFO] Files destination folder is %s\n", rootDestination)

	if !compactor.ExistDirectory(rootSource) {
		print(Fatal, "[ERROR] Files source folder does not exists\n")
		return
	}

	// Options
	options := compactor.Options{
		Source:      rootSource,
		Destination: rootDestination,
		Watch:       *watch,
		Minify:      *minify,
		SourceMap:   *sourceMap,
		Compress:    *compress,
		Progressive: *progressive,
		Include:     include,
		Exclude:     exclude,
		Maps:        maps,
	}

	if options.Watch {
		print(Info, "[INFO] Running in watch mode!\n")
		runWatcher(&options)
	} else {
		print(Info, "[INFO] Running in process and exit mode\n")
		runDefault(&options)
		print(Success, "[SUCCESS] Done\n")
	}

}

func runWatcher(options *compactor.Options) {

	w := watcher.New()
	w.SetMaxEvents(1)

	go func() {
		for {
			select {
			case event := <-w.Event:

				// print(Warn, "[EVENT] %v\n", event)

				if !event.IsDir() {
					if event.Op&watcher.Create == watcher.Create {
						processFile(event.Path, options)
					} else if event.Op&watcher.Write == watcher.Write {
						processFile(event.Path, options)
					} else if event.Op&watcher.Chmod == watcher.Chmod {
						processFile(event.Path, options)
					} else if event.Op&watcher.Rename == watcher.Rename {
						deleteFile(event.OldPath, options)
						processFile(event.Path, options)
					} else if event.Op&watcher.Remove == watcher.Remove {
						deleteFile(event.Path, options)
					}
				}

			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	err := w.AddRecursive(options.Source)
	if err != nil {
		log.Fatalln(err)
	}

	err = w.Start(time.Millisecond * 100)
	if err != nil {
		log.Fatalln(err)
	}

}

func runDefault(options *compactor.Options) {

	files, err := compactor.ListFiles(options.Source)

	if err != nil {
		log.Fatalln(err)
	}

	for _, filename := range files {
		processFile(filename, options)
	}

}

func processFile(filename string, options *compactor.Options) {

	context, err := compactor.Process(
		filename,
		options,
	)

	if err != nil {
		print(Fatal, "[ERROR] %s\n", context.Path)
		print(Warn, "%v\n", err)
	} else {
		print(Success, "[PROCESSED] %s\n", context.Path)
	}

}

func deleteFile(filename string, options *compactor.Options) {

	clean := strings.Replace(filename, options.Source, "", 1)
	destination := strings.Replace(filename, options.Source, options.Destination, 1)

	err := compactor.DeleteFile(destination)

	if err != nil {
		print(Fatal, "[ERROR] %s\n", clean)
		print(Warn, "%v\n", err)
	} else {
		print(Warn, "[DELETED] %s\n", clean)
	}

}
