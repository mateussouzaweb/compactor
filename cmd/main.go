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

	err := compactor.Process(bundle)

	if err != nil {
		os.Printf(os.Fatal, "[ERROR] %s\n", bundle.Destination)
		os.Printf(os.Warn, "%v\n", err)
		return err
	}

	for _, f := range bundle.Logs.Processed {
		os.Printf(os.Success, "[PROCESSED] %s\n", bundle.CleanPath(f))
	}
	for _, f := range bundle.Logs.Skipped {
		os.Printf(os.Warn, "[SKIPPED] %s\n", bundle.CleanPath(f))
	}
	for _, f := range bundle.Logs.Ignored {
		os.Printf(os.Warn, "[IGNORED] %s\n", bundle.CleanPath(f))
	}
	for _, f := range bundle.Logs.Written {
		os.Printf(os.Success, "[WRITTEN] %s\n", bundle.CleanPath(f))
	}
	for _, f := range bundle.Logs.Deleted {
		os.Printf(os.Warn, "[DELETED] %s\n", bundle.CleanPath(f))
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
	watch := false
	version := false
	source, _ := filepath.Abs("src/")
	destination, _ := filepath.Abs("dist/")
	maps := map[string][]string{}

	options := compactor.Bundle{
		Source:      compactor.Source{Path: source},
		Destination: compactor.Destination{Path: destination, Hashed: true},
		Compress:    compactor.Compress{Enabled: true},
		SourceMap:   compactor.SourceMap{Enabled: true},
		Progressive: compactor.Progressive{Enabled: true},
	}

	// Command line flags
	var err error

	flag.Func(
		"source",
		"Path of project source files [DEFAULT: /src]",
		func(path string) error {

			source, err = filepath.Abs(path)
			if err == nil {
				options.Source.Path = source
			}

			return err
		})

	flag.Func(
		"destination",
		"Path to the destination folder [DEFAULT: /dist]",
		func(path string) error {

			destination, err = filepath.Abs(path)
			if err == nil {
				options.Destination.Path = destination
			}

			return err
		})

	flag.Func(
		"hashed",
		"Set if destination file should be hashed to avoid caching on outputs that support it [DEFAULT: true]",
		func(value string) error {

			enabled := trueOrFalse(value)
			options.Destination.Hashed = enabled

			return nil
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
		"Create bundled version from one or multiple files by mapping matching files from given pattern to target destination file or folder",
		func(value string) error {

			split := strings.Split(value, ":")
			files := strings.Split(split[0], ",")
			target := split[1]

			maps[target] = files

			return nil
		})

	flag.Func(
		"disable",
		"Comma separated. Defines which plugin should be disabled. When a plugin is disabled, it uses the generic plugin instead (just copy to destination)",
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
		os.Printf("", "Compactor version 0.0.4\n")
		return
	}

	os.Printf(os.Purple, ":::| COMPACTOR - 0.0.4 |:::\n")
	os.Printf(os.Notice, "[INFO] Files source folder is %s\n", source)
	os.Printf(os.Notice, "[INFO] Files destination folder is %s\n", destination)

	if !os.Exist(source) {
		os.Printf(os.Fatal, "[ERROR] Files source folder does not exists\n")
		return
	}

	// Index source files
	compactor.Index(options.Source.Path)

	// Set as default model
	compactor.DefaultBundle(&options)

	// Create custom defined maps to process
	for target, files := range maps {

		for index, file := range files {
			files[index] = options.CleanPath(file)
		}

		// TODO: target should be mapped on index to allow checksum tracking
		// Maybe this need to goes to every file on dest
		target = options.CleanPath(target)
		compactor.Map(files, target)

	}

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
