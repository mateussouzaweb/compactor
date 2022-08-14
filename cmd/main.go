package main

import (
	"flag"
	_os "os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
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
func process(options *compactor.Options, file *compactor.File) error {

	start := time.Now().UnixNano() / int64(time.Millisecond)
	err := compactor.Process(options, file)

	end := time.Now().UnixNano() / int64(time.Millisecond)
	processTime := end - start

	if err != nil {
		os.Printf(os.Fatal, "[ERROR] %s - %dms\n%v\n", file.Location, processTime, err)
		return err
	}

	os.Printf(os.Success, "[PROCESSED] %s - %dms\n", file.Location, processTime)

	return nil
}

// shutdown runs cleanup process before exiting the program
func shutdown(options *compactor.Options) error {

	err := compactor.Shutdown(options)

	if err != nil {
		os.Printf(os.Fatal, "[ERROR] %v\n", err)
		_os.Exit(1)
		return err
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
	debug := false
	watch := false
	server := false
	serverPort := "5000"

	source, _ := filepath.Abs("src/")
	destination, _ := filepath.Abs("dist/")

	options := &compactor.Options{
		Source: compactor.Source{
			Path: source,
		},
		Destination: compactor.Destination{
			Path:   destination,
			Hashed: true,
		},
		Compress: compactor.Compress{
			Enabled: true,
		},
		SourceMap: compactor.SourceMap{
			Enabled: true,
		},
		Progressive: compactor.Progressive{
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
				server = true
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
				server = true
				serverPort = strings.Replace(value, ":", "", 1)
			} else {
				server = trueOrFalse(value)
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
		os.Printf("", "Compactor version 0.1.2\n")
		return
	}

	os.Printf(os.Purple, ":::| COMPACTOR - 0.1.2 |:::\n")
	os.Printf(os.Notice, "[INFO] Files source folder is %s\n", options.Source.Path)
	os.Printf(os.Notice, "[INFO] Files destination folder is %s\n", options.Destination.Path)

	if !os.Exist(options.Source.Path) {
		os.Printf(os.Fatal, "[ERROR] Files source folder does not exists\n")
		_os.Exit(1)
		return
	}

	// Start a signal watcher to capture program interrupt
	exit := make(chan _os.Signal, 1)
	signal.Notify(exit, _os.Interrupt, syscall.SIGTERM)

	go func() {
		<-exit
		shutdown(options)
		os.Printf(os.Notice, "[INFO] Goodbye :)")
		_os.Exit(0)
	}()

	// Index source files
	err := compactor.IndexFiles(options, options.Source.Path)

	if err != nil {
		os.Printf(os.Fatal, "[ERROR] %v\n", err)
		_os.Exit(1)
		return
	}

	// Detect packages
	packages := compactor.FindPackages(options)

	// Debug info
	if debug {

		os.Printf(os.Purple, "[DEBUG] --- RUNTIME SETTINGS ---\n")
		os.Printf(os.Notice, "[DEBUG] Source ==> %+v\n", options.Source)
		os.Printf(os.Notice, "[DEBUG] Destination ==> %+v\n", options.Destination)
		os.Printf(os.Notice, "[DEBUG] Compress ==> %+v\n", options.Compress)
		os.Printf(os.Notice, "[DEBUG] SourceMap ==> %+v\n", options.SourceMap)
		os.Printf(os.Notice, "[DEBUG] Progressive ==> %+v\n", options.Progressive)
		os.Printf(os.Notice, "[DEBUG] Watch ==> %+v\n", watch)
		os.Printf(os.Notice, "[DEBUG] Server ==> %+v\n", server)
		os.Printf(os.Notice, "[DEBUG] Server Port ==> %+v\n", serverPort)

		os.Printf(os.Purple, "[DEBUG] --- INDEXED FILES ---\n")
		for _, file := range compactor.GetFiles() {
			os.Printf(os.Notice, "[DEBUG] %s\n", options.CleanPath(file.Path))
		}

		os.Printf(os.Purple, "[DEBUG] --- FINAL PACKAGES ---\n")
		for _, file := range packages {
			os.Printf(os.Notice, "[DEBUG] %s", options.CleanPath(file.Path))
			os.Printf(os.Purple, " ==> %s\n", options.CleanPath(file.Destination))
		}

	}

	// Server process
	if server {

		go func() {
			os.Printf(os.Notice, "[INFO] Starting server at \033[1m%s\033[0m\n", "http://localhost:"+serverPort)
			os.Server(
				options.Destination.Path,
				serverPort,
				func(uri string) error {
					os.Printf(os.Notice, "[GET] %s\n", uri)
					return nil
				},
			)
		}()

	}

	// Watch mode
	if watch {

		go func() {
			os.Printf(os.Notice, "[INFO] Starting file watch process\n")
			os.Watch(
				options.Source.Path,
				func(path string, action string) error {

					err = compactor.IndexFiles(options, os.Dir(path))

					if err != nil {
						os.Printf(os.Fatal, "[ERROR] %v\n", err)
						_os.Exit(1)
						return err
					}

					packages = compactor.FindPackages(options)
					file := compactor.FindPackage(options, path)

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
	os.Printf(os.Notice, "[INFO] Running compilation on each package\n")

	for _, item := range packages {
		process(options, item)
	}

	// Keep process alive
	if watch || server {
		<-exit
	}

	// Shutdown
	shutdown(options)

}
