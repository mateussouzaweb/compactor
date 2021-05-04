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

func trueOrFalse(value string) bool {
	value = strings.ToLower(value)
	if value == "true" || value == "t" || value == "1" {
		return true
	}
	return false
}

func main() {

	// Options
	source, _ := filepath.Abs("src/")
	destination, _ := filepath.Abs("dist/")

	options := compactor.Options{
		Source:      source,
		Destination: destination,
		Development: false,
		Watch:       false,
		Include:     []string{},
		Exclude:     []string{},
		Ignore:      []string{},
		Compress: compactor.Compress{
			Enabled: true,
		},
		SourceMap: compactor.SourceMap{
			Enabled: true,
		},
		Progressive: compactor.Progressive{
			Enabled: true,
		},
		Bundles: []compactor.Bundle{},
	}

	// Parsers
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
	compactor.Add("htm", html.Processor)
	compactor.Add("svg", svg.Processor)
	compactor.Add("gif", gif.Processor)
	compactor.Add("jpeg", jpeg.Processor)
	compactor.Add("jpg", jpeg.Processor)
	compactor.Add("png", png.Processor)
	compactor.Add("webp", webp.Processor)

	// Command line flags
	flag.Func(
		"source",
		"Path of project source files [DEFAULT: /src]",
		func(path string) error {
			options.Source, _ = filepath.Abs(path)
			return nil
		})

	flag.Func(
		"destination",
		"Path to the destination folder [DEFAULT: /dist]",
		func(path string) error {
			options.Destination, _ = filepath.Abs(path)
			return nil
		})

	flag.BoolVar(
		&options.Development,
		"development",
		options.Development,
		"Run on development mode (no compression) [DEFAULT: false]")

	flag.BoolVar(
		&options.Watch,
		"watch",
		options.Watch,
		"Enable watcher for live compilation [DEFAULT: false]")

	flag.Func(
		"include",
		"Only include matching files from the given pattern",
		func(value string) error {
			patterns := strings.Split(value, ",")
			options.Include = append(options.Include, patterns...)
			return nil
		})

	flag.Func(
		"exclude",
		"Exclude matching files from the given pattern",
		func(value string) error {
			patterns := strings.Split(value, ",")
			options.Exclude = append(options.Exclude, patterns...)
			return nil
		})

	flag.Func(
		"ignore",
		"Ignore matching files from the given pattern",
		func(value string) error {
			patterns := strings.Split(value, ",")
			options.Ignore = append(options.Ignore, patterns...)
			return nil
		})

	flag.Func(
		"compress",
		"Compress or minify code/images to reduce size [DEFAULT: true]",
		func(value string) error {

			split := strings.Split(value, ":")
			enabled := trueOrFalse(split[0])

			if len(split) > 1 {
				extensions := strings.Split(split[1], ",")
				if enabled {
					options.Compress.Include = append(options.Compress.Include, extensions...)
				} else {
					options.Compress.Exclude = append(options.Compress.Exclude, extensions...)
				}
			} else {
				options.Compress.Enabled = enabled
			}

			return nil
		})

	flag.Func(
		"source-map",
		"Include source map on compilation [DEFAULT: true]",
		func(value string) error {

			split := strings.Split(value, ":")
			enabled := trueOrFalse(split[0])

			if len(split) > 1 {
				extensions := strings.Split(split[1], ",")
				if enabled {
					options.SourceMap.Include = append(options.SourceMap.Include, extensions...)
				} else {
					options.SourceMap.Exclude = append(options.SourceMap.Exclude, extensions...)
				}
			} else {
				options.SourceMap.Enabled = enabled
			}

			return nil
		})

	flag.Func(
		"progressive",
		"Generate new images formats from origin as progressive enhancement [DEFAULT: true]",
		func(value string) error {

			split := strings.Split(value, ":")
			enabled := trueOrFalse(split[0])

			if len(split) > 1 {
				extensions := strings.Split(split[1], ",")
				if enabled {
					options.Progressive.Include = append(options.Progressive.Include, extensions...)
				} else {
					options.Progressive.Exclude = append(options.Progressive.Exclude, extensions...)
				}
			} else {
				options.Progressive.Enabled = enabled
			}

			return nil
		})

	flag.Func(
		"bundle",
		"Create bundled final version from multiple files. Map matching files from the given pattern to destination",
		func(value string) error {

			split := strings.Split(value, ":")
			destination := split[0]
			files := strings.Split(split[1], ",")

			options.Bundles = append(options.Bundles, compactor.Bundle{
				Destination: destination,
				Files:       files,
			})

			return nil
		})

	flag.Func(
		"disable",
		"Comma separated. Defines which processors should be disabled. When a processor is disabled, it uses the generic copy processor",
		func(value string) error {

			list := strings.Split(value, ",")
			for _, item := range list {
				compactor.Remove(compactor.Extension(item))
			}

			return nil
		})

	// Parse values
	version := flag.Bool("version", false, "Print program version")
	flag.Parse()

	// Print information
	if *version {
		print("", "Compactor version 0.0.2\n")
		return
	}

	print(Purple, ":::| COMPACTOR - 0.0.2 |:::\n")
	print(Info, "[INFO] Files source folder is %s\n", options.Source)
	print(Info, "[INFO] Files destination folder is %s\n", options.Destination)

	if !compactor.ExistDirectory(options.Source) {
		print(Fatal, "[ERROR] Files source folder does not exists\n")
		return
	}

	if options.Watch {
		print(Info, "[INFO] Running in watch mode!\n")
		runDefault(&options)
		runWatcher(&options)
	} else {
		print(Info, "[INFO] Running in process and exit mode\n")
		runDefault(&options)
	}

}

func runWatcher(options *compactor.Options) {

	w := watcher.New()

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
					} else if event.Op&watcher.Move == watcher.Move {
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
	} else if context.Skipped {
		print(Warn, "[SKIPPED] %s\n", context.Path)
	} else if context.Processed {
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
