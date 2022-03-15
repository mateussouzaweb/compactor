package main

import (
	"flag"
	"path/filepath"
	"strings"
	"time"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
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
)

func trueOrFalse(value string) bool {
	value = strings.ToLower(value)
	if value == "true" || value == "t" || value == "1" {
		return true
	}
	return false
}

func processBundle(bundle *compactor.Bundle) error {

	start := time.Now().UnixNano() / int64(time.Millisecond)
	err := compactor.Process(bundle)

	end := time.Now().UnixNano() / int64(time.Millisecond)
	processTime := end - start

	if err != nil {
		os.Printf(os.Fatal, "[ERROR] %s - %dms \n", bundle.Destination, processTime)
		os.Printf(os.Fatal, "%v\n", err)
		return err
	}

	for _, f := range bundle.Logs.Processed {
		os.Printf(os.Success, "[PROCESSED] %s - %dms \n", bundle.CleanPath(f), processTime)
	}
	for _, f := range bundle.Logs.Skipped {
		os.Printf(os.Warn, "[SKIPPED] %s - %dms \n", bundle.CleanPath(f), processTime)
	}
	for _, f := range bundle.Logs.Ignored {
		os.Printf(os.Warn, "[IGNORED] %s - %dms \n", bundle.CleanPath(f), processTime)
	}
	for _, f := range bundle.Logs.Written {
		os.Printf(os.Success, "[WRITTEN] %s - %dms \n", bundle.CleanPath(f), processTime)
	}
	for _, f := range bundle.Logs.Deleted {
		os.Printf(os.Warn, "[DELETED] %s - %dms \n", bundle.CleanPath(f), processTime)
	}

	return err
}

func main() {

	// Plugins
	compactor.Register(sass.Plugin())
	compactor.Register(css.Plugin())
	compactor.Register(typescript.Plugin())
	compactor.Register(javascript.Plugin())
	compactor.Register(json.Plugin())
	compactor.Register(xml.Plugin())
	compactor.Register(html.Plugin())
	compactor.Register(svg.Plugin())
	compactor.Register(gif.Plugin())
	compactor.Register(jpeg.Plugin())
	compactor.Register(png.Plugin())
	compactor.Register(webp.Plugin())
	compactor.Register(generic.Plugin())

	// compactor.Register("less", less.Plugin())
	// compactor.Register("styl", stylus.Plugin())
	// compactor.Register("apng", apng.Plugin())
	// compactor.Register("avif", avif.Plugin())
	// compactor.Register("ico", ico.Plugin())
	// compactor.Register("js", babel.Plugin())
	// compactor.Register("js", react.Plugin())
	// compactor.Register("jsx", react.Plugin())
	// compactor.Register("js", vue.Plugin())
	// compactor.Register("vue", vue.Plugin())
	// compactor.Register("js", svelte.Plugin())
	// compactor.Register("svelte", svelte.Plugin())
	// compactor.Register("coffee", coffee.Plugin())
	// compactor.Register("elm", elm.Plugin())
	// compactor.Register("eot", eot.Plugin())
	// compactor.Register("ttf", ttf.Plugin())
	// compactor.Register("woff", woff.Plugin())
	// compactor.Register("gql", graphql.Plugin())
	// compactor.Register("graphql", graphql.Plugin())
	// compactor.Register("yaml", yaml.Plugin())
	// compactor.Register("toml", toml.Plugin())

	// Options
	version := false
	watch := false
	source, _ := filepath.Abs("src/")
	destination, _ := filepath.Abs("dist/")

	options := compactor.Default
	options.Source.Path = source
	options.Destination.Path = destination

	// Command line flags
	flag.BoolVar(
		&version,
		"version",
		false,
		"Description: Print program version")

	flag.BoolVar(
		&watch,
		"watch",
		false,
		"Default: false\nDescription: Enable watcher for live compilation")

	flag.Func(
		"source",
		"Default: /src\nDescription: Set the path of project source files",
		func(path string) error {

			source, err := filepath.Abs(path)
			if err == nil {
				options.Source.Path = source
			}

			return err
		})

	flag.Func(
		"destination",
		"Default: /dist\nDescription: Set the path to the destination folder",
		func(path string) error {

			destination, err := filepath.Abs(path)
			if err == nil {
				options.Destination.Path = destination
			}

			return err
		})

	flag.Func(
		"hashed",
		"Default: true\nDescription: Set if destination file should be hashed to avoid caching on outputs that support it\nImportant: If you are running in watch mode, we recommend to disable this option",
		func(value string) error {

			enabled := trueOrFalse(value)
			options.Destination.Hashed = enabled

			return nil
		})

	flag.Func(
		"include",
		"Description: Only include matching files from the given pattern",
		func(value string) error {
			patterns := strings.Split(value, ",")
			options.Source.Include = append(options.Source.Include, patterns...)
			return nil
		})

	flag.Func(
		"exclude",
		"Description: Exclude matching files from the given pattern",
		func(value string) error {
			patterns := strings.Split(value, ",")
			options.Source.Exclude = append(options.Source.Exclude, patterns...)
			return nil
		})

	flag.Func(
		"compress",
		"Default: true\nFormats: [BOOLEAN] or [PATTERN,...]:[BOOLEAN]\nDescription: Define if should compress or minify code/images to reduce size",
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
		"Default: true\nFormats: [BOOLEAN] or [PATTERN,...]:[BOOLEAN]\nDescription: Define if should include source map reference on compilation",
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
		"Default: true\nFormats: [BOOLEAN] or [PATTERN,...]:[BOOLEAN]\nDescription: Define if should generate new images formats from origin as progressive enhancement",
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
		"Formats: [FILES]:[DESTINATION] or [FILES]:[DEPENDENCIES]:[DESTINATION]\nDescription: Create bundled version from one or multiple files and dependencies by mapping matching files from given pattern to target destination file or folder",
		func(value string) error {

			var destination string
			var files []string
			var dependencies []string

			split := strings.Split(value, ":")
			files = strings.Split(split[0], ",")

			if len(split) > 2 {
				dependencies = strings.Split(split[1], ",")
				destination = split[2]
			} else {
				destination = split[1]
			}

			for index, file := range files {
				files[index] = options.CleanPath(file)
			}
			for index, file := range dependencies {
				dependencies[index] = options.CleanPath(file)
			}

			destination = options.CleanPath(destination)

			// TODO: Destination should be mapped on index to allow checksum tracking
			// Maybe this need to goes to every file on dest
			compactor.Map(files, dependencies, destination)

			return nil
		})

	flag.Func(
		"disable",
		"Format: [EXTENSION,...]\nDescription: Defines which plugin should be disabled. When a plugin is disabled, it uses the generic plugin instead (just copy to destination)",
		func(value string) error {

			list := strings.Split(value, ",")

			for _, item := range list {
				if !strings.Contains(item, ".") {
					compactor.Unregister("." + item)
				} else {
					compactor.Unregister(item)
				}
			}

			return nil
		})

	// Parse values
	flag.Parse()

	// Print information
	if version {
		os.Printf("", "Compactor version 0.0.11\n")
		return
	}

	os.Printf(os.Purple, ":::| COMPACTOR - 0.0.11 |:::\n")
	os.Printf(os.Notice, "[INFO] Files source folder is %s\n", options.Source.Path)
	os.Printf(os.Notice, "[INFO] Files destination folder is %s\n", options.Destination.Path)

	if !os.Exist(options.Source.Path) {
		os.Printf(os.Fatal, "[ERROR] Files source folder does not exists\n")
		return
	}

	// Index source files
	compactor.Index(options.Source.Path)

	// Run compactor processing
	if watch {
		os.Printf(os.Notice, "[INFO] Running in watch mode!\n")
	} else {
		os.Printf(os.Notice, "[INFO] Running in process and exit mode\n")
	}

	for _, bundle := range compactor.GetBundles() {
		processBundle(bundle)
	}

	if !watch {
		return
	}

	os.Watch(
		options.Source.Path,
		time.Millisecond*250,
		func(path string) error {

			existing := compactor.Get(path)

			if existing.Path == "" {
				compactor.Append(path)
			} else {
				compactor.Update(path)
			}

			bundle := compactor.GetBundleFor(path)
			processBundle(bundle)

			// TODO: .html files should be reprocessed when dependencies update

			return nil
		},
		func(path string) error {

			bundle := compactor.GetBundleFor(path)
			compactor.Remove(path)

			processBundle(bundle)

			return nil
		},
	)

}
