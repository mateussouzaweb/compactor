package main

import (
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/mateussouzaweb/compactor/src/cli"
	"github.com/mateussouzaweb/compactor/src/plugins/css"
	"github.com/mateussouzaweb/compactor/src/plugins/generic"
	"github.com/mateussouzaweb/compactor/src/plugins/gif"
	"github.com/mateussouzaweb/compactor/src/plugins/html"
	"github.com/mateussouzaweb/compactor/src/plugins/javascript"
	"github.com/mateussouzaweb/compactor/src/plugins/jpeg"
	"github.com/mateussouzaweb/compactor/src/plugins/json"
	"github.com/mateussouzaweb/compactor/src/plugins/png"
	"github.com/mateussouzaweb/compactor/src/plugins/sass"
	"github.com/mateussouzaweb/compactor/src/plugins/svg"
	"github.com/mateussouzaweb/compactor/src/plugins/typescript"
	"github.com/mateussouzaweb/compactor/src/plugins/webp"
	"github.com/mateussouzaweb/compactor/src/plugins/xml"
	"github.com/mateussouzaweb/compactor/src/processor"
	"github.com/mateussouzaweb/compactor/src/server"
	"github.com/mateussouzaweb/compactor/src/system"
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
func process(options *processor.Options, file *processor.File) error {

	start := time.Now().UnixNano() / int64(time.Millisecond)
	err := processor.Process(options, file)

	end := time.Now().UnixNano() / int64(time.Millisecond)
	processTime := end - start

	if err != nil {
		cli.Printf(cli.Fatal, "[ERROR] %s - %dms\n%v\n", file.Location, processTime, err)
		return err
	}

	cli.Printf(cli.Success, "[PROCESSED] %s - %dms\n", file.Location, processTime)

	return nil
}

// shutdown runs cleanup process before exiting the program
func shutdown(options *processor.Options) error {

	err := processor.Shutdown(options)
	if err != nil {
		cli.Printf(cli.Fatal, "[ERROR] %v\n", err)
		os.Exit(1)
		return err
	}

	return nil
}

// main runs the program
func main() {

	// Plugins
	// processor.AddPlugin("less", less.Plugin())
	// processor.AddPlugin("styl", stylus.Plugin())
	// processor.AddPlugin("apng", apng.Plugin())
	// processor.AddPlugin("avif", avif.Plugin())
	// processor.AddPlugin("ico", ico.Plugin())
	// processor.AddPlugin("js", babel.Plugin())
	// processor.AddPlugin("js", react.Plugin())
	// processor.AddPlugin("jsx", react.Plugin())
	// processor.AddPlugin("js", vue.Plugin())
	// processor.AddPlugin("vue", vue.Plugin())
	// processor.AddPlugin("js", svelte.Plugin())
	// processor.AddPlugin("svelte", svelte.Plugin())
	// processor.AddPlugin("coffee", coffee.Plugin())
	// processor.AddPlugin("elm", elm.Plugin())
	// processor.AddPlugin("eot", eot.Plugin())
	// processor.AddPlugin("ttf", ttf.Plugin())
	// processor.AddPlugin("woff", woff.Plugin())
	// processor.AddPlugin("gql", graphql.Plugin())
	// processor.AddPlugin("graphql", graphql.Plugin())
	// processor.AddPlugin("yaml", yaml.Plugin())
	// processor.AddPlugin("toml", toml.Plugin())

	processor.AddPlugin(sass.Plugin())
	processor.AddPlugin(css.Plugin())
	processor.AddPlugin(javascript.Plugin())
	processor.AddPlugin(typescript.Plugin())
	processor.AddPlugin(json.Plugin())
	processor.AddPlugin(xml.Plugin())
	processor.AddPlugin(html.Plugin())
	processor.AddPlugin(svg.Plugin())
	processor.AddPlugin(gif.Plugin())
	processor.AddPlugin(jpeg.Plugin())
	processor.AddPlugin(png.Plugin())
	processor.AddPlugin(webp.Plugin())
	processor.AddPlugin(generic.Plugin())

	// Options
	version := false
	debug := false
	watch := false
	serverMode := false
	serverPort := "5000"

	source, _ := filepath.Abs("src/")
	destination, _ := filepath.Abs("dist/")

	options := &processor.Options{
		Source: processor.Source{
			Path: source,
		},
		Destination: processor.Destination{
			Path:   destination,
			Hashed: true,
		},
		Compress: processor.Compress{
			Enabled: true,
		},
		SourceMap: processor.SourceMap{
			Enabled: true,
		},
		Progressive: processor.Progressive{
			Enabled: true,
		},
	}

	// Version flag
	flag.BoolVar(
		&version,
		"version",
		false,
		"Description: Print program version")

	// Debug flag
	flag.BoolVar(
		&debug,
		"debug",
		false,
		"Description: Print debug information")

	// Develop flag
	flag.Func(
		"develop",
		"Default: false\nFormat: [BOOLEAN]\nDescription: Enable or disable development mode. When enabled, disables hash, compression and progressive enhancements and enable watch and server.",
		func(value string) error {

			if trueOrFalse(value) {
				watch = true
				serverMode = true
				options.Destination.Hashed = false
				options.Compress.Enabled = false
				options.Progressive.Enabled = false
			}

			return nil
		},
	)

	// Watch flag
	flag.BoolVar(
		&watch,
		"watch",
		false,
		"Description: Enables file watching to live compile on code change")

	// Server flag
	flag.Func(
		"server",
		"Default: false\nFormats: [BOOLEAN] or :[PORT]\nDescription: Enable or disable local server on given port - if port is not specified, defaults to :5000. NOTE: Not existing paths will be automatically translated to index.html for a SPA like feature.",
		func(value string) error {

			if strings.Contains(value, ":") {
				serverMode = true
				serverPort = strings.Replace(value, ":", "", 1)
			} else {
				serverMode = trueOrFalse(value)
			}

			return nil
		},
	)

	// Compilation flags
	flag.Func(
		"source",
		"Default: /src\nFormat: [PATH]\nDescription: Set the path of project source files",
		func(path string) error {

			source, err := filepath.Abs(path)
			if err == nil {
				options.Source.Path = source
			}

			return err
		})

	flag.Func(
		"include",
		"Format: [PATTERN,...]\nDescription: Only include matching files from the given pattern",
		func(value string) error {
			patterns := strings.Split(value, ",")
			options.Source.Include = append(options.Source.Include, patterns...)
			return nil
		})

	flag.Func(
		"exclude",
		"Format: [PATTERN,...]\nDescription: Exclude matching files from the given pattern",
		func(value string) error {
			patterns := strings.Split(value, ",")
			options.Source.Exclude = append(options.Source.Exclude, patterns...)
			return nil
		})

	flag.Func(
		"destination",
		"Default: /dist\nFormat: [PATH]\nDescription: Set the path to the destination folder",
		func(path string) error {

			destination, err := filepath.Abs(path)
			if err == nil {
				options.Destination.Path = destination
			}

			return err
		})

	flag.Func(
		"hashed",
		"Default: true\nDescription: Defines if destination file should have the hash key in its name to avoid server caching on files that can be constantly updated.",
		func(value string) error {

			enabled := trueOrFalse(value)
			options.Destination.Hashed = enabled

			return nil
		})

	flag.Func(
		"compress",
		"Default: true\nFormats: [BOOLEAN] or [PATTERN,...]:[BOOLEAN]\nDescription: Defines if should compress or minify code/images to reduce size",
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
		"Default: true\nFormats: [BOOLEAN] or [PATTERN,...]:[BOOLEAN]\nDescription: Defines if should include source map reference on file compilation",
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
		"Default: true\nFormats: [BOOLEAN] or [PATTERN,...]:[BOOLEAN]\nDescription: Defines if should generate new images formats from original image format as progressive enhancement",
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

	// Plugin flag
	flag.Func(
		"disable",
		"Format: [PLUGIN,...]\nDescription: Defines which plugin should be disabled. When a plugin is disabled, the next available plugin that matches the file extension will be used. Otherwise, it forces the use of the generic plugin (simple copy to destination)",
		func(value string) error {

			list := strings.SplitSeq(value, ",")
			for namespace := range list {
				processor.RemovePlugin(namespace)
			}

			return nil
		})

	// Parse values
	flag.Parse()

	// Print information
	if version {
		cli.Printf("", "Compactor version 0.2.2\n")
		return
	}

	cli.Printf(cli.Purple, ":::| COMPACTOR - 0.2.2 |:::\n")
	cli.Printf(cli.Notice, "[INFO] Files source folder is %s\n", options.Source.Path)
	cli.Printf(cli.Notice, "[INFO] Files destination folder is %s\n", options.Destination.Path)

	if !system.Exist(options.Source.Path) {
		cli.Printf(cli.Fatal, "[ERROR] Files source folder does not exists\n")
		os.Exit(1)
		return
	}

	// Start a signal watcher to capture program interrupt
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-exit
		shutdown(options)
		cli.Printf(cli.Notice, "[INFO] Goodbye :)")
		os.Exit(0)
	}()

	// Index source files
	err := processor.IndexFiles(options, options.Source.Path)
	if err != nil {
		cli.Printf(cli.Fatal, "[ERROR] %v\n", err)
		os.Exit(1)
		return
	}

	// Detect packages
	packages := processor.FindPackages(options)

	// Debug info
	if debug {

		cli.Printf(cli.Purple, "[DEBUG] --- RUNTIME SETTINGS ---\n")
		cli.Printf(cli.Notice, "[DEBUG] Source ==> %+v\n", options.Source)
		cli.Printf(cli.Notice, "[DEBUG] Destination ==> %+v\n", options.Destination)
		cli.Printf(cli.Notice, "[DEBUG] Compress ==> %+v\n", options.Compress)
		cli.Printf(cli.Notice, "[DEBUG] SourceMap ==> %+v\n", options.SourceMap)
		cli.Printf(cli.Notice, "[DEBUG] Progressive ==> %+v\n", options.Progressive)
		cli.Printf(cli.Notice, "[DEBUG] Watch ==> %+v\n", watch)
		cli.Printf(cli.Notice, "[DEBUG] Server ==> %+v\n", serverMode)
		cli.Printf(cli.Notice, "[DEBUG] Server Port ==> %+v\n", serverPort)

		cli.Printf(cli.Purple, "[DEBUG] --- INDEXED FILES ---\n")
		for _, file := range processor.GetFiles() {
			cli.Printf(cli.Notice, "[DEBUG] %s\n", options.CleanPath(file.Path))
		}

		cli.Printf(cli.Purple, "[DEBUG] --- FINAL PACKAGES ---\n")
		for _, file := range packages {
			cli.Printf(cli.Notice, "[DEBUG] %s", options.CleanPath(file.Path))
			cli.Printf(cli.Purple, " ==> %s\n", options.CleanPath(file.Destination))
		}

	}

	// Server process
	if serverMode {

		go func() {
			cli.Printf(cli.Notice, "[INFO] Starting server at \033[1m%s\033[0m\n", "http://localhost:"+serverPort)
			server.Start(
				options.Destination.Path,
				serverPort,
				func(uri string) error {
					cli.Printf(cli.Notice, "[GET] %s\n", uri)
					return nil
				},
			)
		}()

	}

	// Watch mode
	if watch {

		go func() {
			cli.Printf(cli.Notice, "[INFO] Starting file watch process\n")
			system.Watch(
				options.Source.Path,
				func(path string, action string) error {

					err = processor.IndexFiles(options, system.Dir(path))
					if err != nil {
						cli.Printf(cli.Fatal, "[ERROR] %v\n", err)
						os.Exit(1)
						return err
					}

					packages = processor.FindPackages(options)
					file := processor.FindPackage(options, path)

					if file.Extension == "" {
						return nil
					}

					err := process(options, file)
					if err != nil {
						return err
					}

					// Try to process related packages
					for _, thePackage := range packages {

						// Ignore if is the same package
						if thePackage.Path == file.Path {
							continue
						}

						// Check on related items of the package
						for _, related := range thePackage.Related {
							if !related.Dependency && related.File.Path == file.Path {
								err := process(options, thePackage)
								if err != nil {
									return err
								}
							}
						}

					}

					return nil
				},
			)
		}()

	}

	// Compilation
	cli.Printf(cli.Notice, "[INFO] Running compilation on each package\n")

	for _, item := range packages {
		process(options, item)
	}

	// Keep process alive
	if watch || serverMode {
		<-exit
	}

	// Shutdown
	shutdown(options)

}
