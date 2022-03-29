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

// trueOrFalse returns if given value is likely to be a true or false flag
func trueOrFalse(value string) bool {
	value = strings.ToLower(value)
	if value == "true" || value == "t" || value == "1" {
		return true
	}
	return false
}

// process runs the package processing on the destination plugin
func process(bundle *compactor.Bundle) error {

	start := time.Now().UnixNano() / int64(time.Millisecond)
	err := compactor.Process(bundle)

	end := time.Now().UnixNano() / int64(time.Millisecond)
	processTime := end - start

	file := bundle.CleanPath(bundle.Item.Path)

	if err != nil {
		os.Printf(os.Fatal, "[ERROR] %s - %dms\n%v\n", file, processTime, err)
		return err
	}

	os.Printf(os.Success, "[PROCESSED] %s - %dms\n", file, processTime)

	return nil
}

// processRelated and process the bundle for linked bundles
func processRelated(bundle *compactor.Bundle) error {

	for _, theBundle := range compactor.GetBundles() {

		// Ignore if is the same bundle
		if theBundle.Item.Path == bundle.Item.Path {
			continue
		}

		// Check on related items to the bundle
		for _, related := range theBundle.Item.Related {
			if related.Item.Path == bundle.Item.Path && related.Type == "link" {
				process(theBundle)
			}
		}

	}

	return nil
}

// main runs the program
func main() {

	// Plugins
	// compactor.AddPlugin("less", less.Plugin())
	// compactor.AddPlugin("styl", stylus.Plugin())
	// compactor.AddPlugin("apng", apng.Plugin())
	// compactor.AddPlugin("avif", avif.Plugin())
	// compactor.AddPlugin("ico", ico.Plugin())
	// compactor.AddPlugin("js", babel.Plugin())
	// compactor.AddPlugin("js", react.Plugin())
	// compactor.AddPlugin("jsx", react.Plugin())
	// compactor.AddPlugin("js", vue.Plugin())
	// compactor.AddPlugin("vue", vue.Plugin())
	// compactor.AddPlugin("js", svelte.Plugin())
	// compactor.AddPlugin("svelte", svelte.Plugin())
	// compactor.AddPlugin("coffee", coffee.Plugin())
	// compactor.AddPlugin("elm", elm.Plugin())
	// compactor.AddPlugin("eot", eot.Plugin())
	// compactor.AddPlugin("ttf", ttf.Plugin())
	// compactor.AddPlugin("woff", woff.Plugin())
	// compactor.AddPlugin("gql", graphql.Plugin())
	// compactor.AddPlugin("graphql", graphql.Plugin())
	// compactor.AddPlugin("yaml", yaml.Plugin())
	// compactor.AddPlugin("toml", toml.Plugin())

	compactor.AddPlugin(sass.Plugin())
	compactor.AddPlugin(css.Plugin())
	compactor.AddPlugin(javascript.Plugin())
	compactor.AddPlugin(typescript.Plugin())
	compactor.AddPlugin(json.Plugin())
	compactor.AddPlugin(xml.Plugin())
	compactor.AddPlugin(html.Plugin())
	compactor.AddPlugin(svg.Plugin())
	compactor.AddPlugin(gif.Plugin())
	compactor.AddPlugin(jpeg.Plugin())
	compactor.AddPlugin(png.Plugin())
	compactor.AddPlugin(webp.Plugin())
	compactor.AddPlugin(generic.Plugin())

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
		"disable",
		"Format: [PLUGIN,...]\nDescription: Defines which plugin should be disabled. When a plugin is disabled, the next available plugin that matches the file extension will be used. Otherwise, it forces the use of the generic plugin (simple copy to destination)",
		func(value string) error {

			list := strings.Split(value, ",")

			for _, namespace := range list {
				compactor.RemovePlugin(namespace)
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

	if watch {
		os.Printf(os.Notice, "[INFO] Running in watch mode!\n")
	} else {
		os.Printf(os.Notice, "[INFO] Running in process and exit mode\n")
	}

	// Index source files
	err := compactor.Index(options.Source.Path)

	if err != nil {
		os.Printf(os.Fatal, "[ERROR] %v\n", err)
		return
	}

	// Run bundle processing
	for _, bundle := range compactor.GetBundles() {
		process(bundle)
	}

	if !watch {
		return
	}

	os.Watch(
		options.Source.Path,
		time.Millisecond*250,
		func(path string) error {

			compactor.Index(os.Dir(path))
			bundle := compactor.GetBundle(path)

			err := process(bundle)
			if err != nil {
				return err
			}

			return processRelated(bundle)
		},
		func(path string) error {

			bundle := compactor.GetBundle(path)
			compactor.Remove(path)

			err := process(bundle)
			if err != nil {
				return err
			}

			return processRelated(bundle)
		},
	)

}
