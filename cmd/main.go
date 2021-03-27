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
	watch := flag.Bool("watch", false, "Enable live watch compilation")
	_source := flag.String("source", "src/", "Path of project source files")
	_destination := flag.String("destination", "dist/", "Path to the destination folder")

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
	source, _ := filepath.Abs(*_source)
	destination, _ := filepath.Abs(*_destination)

	// Print information
	print(Purple, ":::| COMPACTOR |:::\n")
	print(Info, "[INFO] Files source folder is %s\n", source)
	print(Info, "[INFO] Files destination folder is %s\n", destination)

	if !compactor.ExistDirectory(source) {
		print(Fatal, "[ERROR] Files source folder does not exists\n")
		return
	}

	if *watch {
		print(Info, "[INFO] Running in watch mode!\n")
		runWatcher(source, destination)
	} else {
		print(Info, "[INFO] Running in process and exit mode\n")
		runDefault(source, destination)
		print(Success, "[SUCCESS] Done\n")
	}

}

func runWatcher(rootSource string, rootDestination string) {

	w := watcher.New()
	w.SetMaxEvents(1)

	go func() {
		for {
			select {
			case event := <-w.Event:

				// print(Warn, "[EVENT] %v\n", event)

				if !event.IsDir() {
					if event.Op&watcher.Create == watcher.Create {
						processFile(event.Path, rootSource, rootDestination)
					} else if event.Op&watcher.Write == watcher.Write {
						processFile(event.Path, rootSource, rootDestination)
					} else if event.Op&watcher.Chmod == watcher.Chmod {
						processFile(event.Path, rootSource, rootDestination)
					} else if event.Op&watcher.Rename == watcher.Rename {
						deleteFile(event.OldPath, rootSource, rootDestination)
						processFile(event.Path, rootSource, rootDestination)
					} else if event.Op&watcher.Remove == watcher.Remove {
						deleteFile(event.Path, rootSource, rootDestination)
					}
				}

			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	err := w.AddRecursive(rootSource)
	if err != nil {
		log.Fatalln(err)
	}

	err = w.Start(time.Millisecond * 100)
	if err != nil {
		log.Fatalln(err)
	}

}

func runDefault(rootSource string, rootDestination string) {

	files, err := compactor.ListFiles(rootSource)

	if err != nil {
		log.Fatalln(err)
	}

	for _, filename := range files {
		processFile(filename, rootSource, rootDestination)
	}

}

func processFile(filename string, source string, destination string) {

	context, err := compactor.Process(
		filename,
		source,
		destination,
		compactor.Options{},
	)

	if err != nil {
		print(Fatal, "[ERROR] %s\n", context.Path)
		print(Warn, "%v\n", err)
	} else {
		print(Success, "[PROCESSED] %s\n", context.Path)
	}

}

func deleteFile(filename string, rootSource string, rootDestination string) {

	clean := strings.Replace(filename, rootSource, "", 1)
	destination := strings.Replace(filename, rootSource, rootDestination, 1)

	err := compactor.DeleteFile(destination)

	if err != nil {
		print(Fatal, "[ERROR] %s\n", clean)
		print(Warn, "%v\n", err)
	} else {
		print(Warn, "[DELETED] %s\n", clean)
	}

}
