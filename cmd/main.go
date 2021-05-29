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

// Colors
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

func processBundle(bundle *compactor.Bundle) {

	logger, err := compactor.Process(bundle)

	if err != nil {
		print(Fatal, "[ERROR] %s\n", bundle.Destination)
		print(Warn, "%v\n", err)
		return
	}

	for _, f := range logger.Processed {
		print(Success, "[PROCESSED] %s\n", bundle.CleanPath(f))
	}
	for _, f := range logger.Skipped {
		print(Warn, "[SKIPPED] %s\n", bundle.CleanPath(f))
	}
	for _, f := range logger.Ignored {
		print(Warn, "[IGNORED] %s\n", bundle.CleanPath(f))
	}
	for _, f := range logger.Written {
		print(Success, "[WRITTEN] %s\n", bundle.CleanPath(f))
	}
	for _, f := range logger.Deleted {
		print(Warn, "[DELETED] %s\n", bundle.CleanPath(f))
	}

}

func processFile(file string) {

	bundle := compactor.RetrieveBundleFor(file)
	processBundle(bundle)

}

func deleteFile(file string) {

	bundle := compactor.RetrieveBundleFor(file)
	bundle.RemoveFile(file)

	processBundle(bundle)

}

func runWatcher(path string) {

	w := watcher.New()

	go func() {
		for {
			select {
			case event := <-w.Event:
				if !event.IsDir() {

					if event.Op&watcher.Create == watcher.Create {
						processFile(event.Path)
					} else if event.Op&watcher.Write == watcher.Write {
						processFile(event.Path)
					} else if event.Op&watcher.Chmod == watcher.Chmod {
						processFile(event.Path)
					} else if event.Op&watcher.Rename == watcher.Rename {
						deleteFile(event.OldPath)
						processFile(event.Path)
					} else if event.Op&watcher.Move == watcher.Move {
						deleteFile(event.OldPath)
						processFile(event.Path)
					} else if event.Op&watcher.Remove == watcher.Remove {
						deleteFile(event.Path)
					}

				}
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	err := w.AddRecursive(path)
	if err != nil {
		log.Fatalln(err)
	}

	err = w.Start(time.Millisecond * 100)
	if err != nil {
		log.Fatalln(err)
	}

}

func main() {

	watch := false
	version := false
	source, _ := filepath.Abs("src/")
	destination, _ := filepath.Abs("dist/")
	bundles := map[string][]string{}

	options := compactor.Bundle{
		Extension:   "*",
		Source:      compactor.Source{Path: source},
		Destination: compactor.Destination{Path: destination},
		Compress:    compactor.Compress{Enabled: true},
		SourceMap:   compactor.SourceMap{Enabled: true},
		Progressive: compactor.Progressive{Enabled: true},
	}

	// Parsers
	compactor.RegisterProcessor("*", generic.Processor)
	compactor.RegisterProcessor("sass", sass.Processor)
	compactor.RegisterProcessor("scss", sass.Processor)
	compactor.RegisterProcessor("css", css.Processor)
	compactor.RegisterProcessor("ts", typescript.Processor)
	compactor.RegisterProcessor("tsx", typescript.Processor)
	compactor.RegisterProcessor("js", javascript.Processor)
	compactor.RegisterProcessor("json", json.Processor)
	compactor.RegisterProcessor("xml", xml.Processor)
	compactor.RegisterProcessor("html", html.Processor)
	compactor.RegisterProcessor("htm", html.Processor)
	compactor.RegisterProcessor("svg", svg.Processor)
	compactor.RegisterProcessor("gif", gif.Processor)
	compactor.RegisterProcessor("jpeg", jpeg.Processor)
	compactor.RegisterProcessor("jpg", jpeg.Processor)
	compactor.RegisterProcessor("png", png.Processor)
	compactor.RegisterProcessor("webp", webp.Processor)

	// compactor.RegisterProcessor("less", less.Processor)
	// compactor.RegisterProcessor("styl", stylus.Processor)
	// compactor.RegisterProcessor("apng", apng.Processor)
	// compactor.RegisterProcessor("avif", avif.Processor)
	// compactor.RegisterProcessor("ico", ico.Processor)
	// compactor.RegisterProcessor("js", babel.Processor)
	// compactor.RegisterProcessor("js", react.Processor)
	// compactor.RegisterProcessor("jsx", react.Processor)
	// compactor.RegisterProcessor("js", vue.Processor)
	// compactor.RegisterProcessor("vue", vue.Processor)
	// compactor.RegisterProcessor("js", svelte.Processor)
	// compactor.RegisterProcessor("svelte", svelte.Processor)
	// compactor.RegisterProcessor("coffee", coffee.Processor)
	// compactor.RegisterProcessor("elm", elm.Processor)
	// compactor.RegisterProcessor("eot", eot.Processor)
	// compactor.RegisterProcessor("ttf", ttf.Processor)
	// compactor.RegisterProcessor("woff", woff.Processor)
	// compactor.RegisterProcessor("gql", graphql.Processor)
	// compactor.RegisterProcessor("graphql", graphql.Processor)
	// compactor.RegisterProcessor("yaml", yaml.Processor)
	// compactor.RegisterProcessor("toml", toml.Processor)

	// Command line flags
	flag.Func(
		"source",
		"Path of project source files [DEFAULT: /src]",
		func(path string) error {

			source, err := filepath.Abs(path)
			if err != nil {
				options.Source.Path = source
			}

			return err
		})

	flag.Func(
		"destination",
		"Path to the destination folder [DEFAULT: /dist]",
		func(path string) error {

			destination, err := filepath.Abs(path)
			if err != nil {
				options.Destination.Path = destination
			}

			return err
		})

	flag.Func(
		"include",
		"Only include matching files from the given pattern",
		func(value string) error {
			patterns := strings.Split(value, ",")
			options.Source.Include = append(options.Source.Include, patterns...)
			return nil
		})

	flag.Func(
		"exclude",
		"Exclude matching files from the given pattern",
		func(value string) error {
			patterns := strings.Split(value, ",")
			options.Source.Exclude = append(options.Source.Exclude, patterns...)
			return nil
		})

	flag.Func(
		"ignore",
		"Ignore matching files from the given pattern",
		func(value string) error {
			patterns := strings.Split(value, ",")
			options.Source.Ignore = append(options.Source.Ignore, patterns...)
			return nil
		})

	flag.Func(
		"compress",
		"Compress or minify code/images to reduce size [DEFAULT: true]",
		func(value string) error {

			split := strings.Split(value, ":")
			enabled := trueOrFalse(split[0])

			if len(split) > 1 {

				patterns := strings.Split(split[1], ",")

				if enabled {
					options.Compress.Include = append(
						options.Compress.Include,
						patterns...,
					)
				} else {
					options.Compress.Exclude = append(
						options.Compress.Exclude,
						patterns...,
					)
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

				patterns := strings.Split(split[1], ",")

				if enabled {
					options.SourceMap.Include = append(
						options.SourceMap.Include,
						patterns...,
					)
				} else {
					options.SourceMap.Exclude = append(
						options.SourceMap.Exclude,
						patterns...,
					)
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

				patterns := strings.Split(split[1], ",")

				if enabled {
					options.Progressive.Include = append(
						options.Progressive.Include,
						patterns...,
					)
				} else {
					options.Progressive.Exclude = append(
						options.Progressive.Exclude,
						patterns...,
					)
				}

			} else {
				options.Progressive.Enabled = enabled
			}

			return nil
		})

	flag.Func(
		"bundle",
		"Create bundled final version from one or multiple files. Map matching files from the given pattern to target destination file",
		func(value string) error {

			split := strings.Split(value, ":")
			target := split[0]
			files := strings.Split(split[1], ",")

			bundles[target] = files

			return nil
		})

	flag.Func(
		"disable",
		"Comma separated. Defines which processors should be disabled. When a processor is disabled, it uses the generic copy processor",
		func(value string) error {

			list := strings.Split(value, ",")
			for _, item := range list {
				compactor.RemoveProcessors(item)
			}

			return nil
		})

	flag.BoolVar(
		&watch,
		"watch",
		false,
		"Enable watcher for live compilation [DEFAULT: false]")

	flag.BoolVar(
		&version,
		"version",
		false,
		"Print program version")

	// Parse values
	flag.Parse()

	// Print information
	if version {
		print("", "Compactor version 0.0.3\n")
		return
	}

	print(Purple, ":::| COMPACTOR - 0.0.3 |:::\n")
	print(Info, "[INFO] Files source folder is %s\n", source)
	print(Info, "[INFO] Files destination folder is %s\n", destination)

	if !compactor.ExistDirectory(source) {
		print(Fatal, "[ERROR] Files source folder does not exists\n")
		return
	}

	// Set as default model
	compactor.Default = &options

	// Create custom defined bundles to process
	for target, files := range bundles {

		bundle := compactor.NewBundle()
		bundle.Extension = bundle.CleanExtension(files[0])
		bundle.Destination.File = bundle.CleanPath(target)

		for _, file := range files {
			bundle.AddFile(file)
		}

		compactor.RegisterBundle(bundle)

	}

	// Create default bundles from files
	files, err := compactor.ListFiles(options.Source.Path)

	if err != nil {
		print(Fatal, "[ERROR] Bundles could not be created: %v\n", err)
		return
	}

	for _, file := range files {
		_ = compactor.RetrieveBundleFor(file)
	}

	// Run compactor processing
	if watch {
		print(Info, "[INFO] Running in watch mode!\n")
	} else {
		print(Info, "[INFO] Running in process and exit mode\n")
	}

	for _, bundle := range compactor.RetrieveBundles() {
		processBundle(bundle)
	}

	if watch {
		runWatcher(options.Source.Path)
	}

}
