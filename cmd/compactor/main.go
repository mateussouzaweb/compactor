package main

import (
	"os"
	"os/signal"
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

	// Read options from arguments
	context := readContext()
	options := context.Options

	// Print information
	if context.Version {
		cli.Printf("", "Compactor version 0.3.2\n")
		return
	}

	cli.Printf(cli.Purple, ":::| COMPACTOR - 0.3.2 |:::\n")
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
	if context.DebugMode {

		cli.Printf(cli.Purple, "[DEBUG] --- RUNTIME SETTINGS ---\n")
		cli.Printf(cli.Notice, "[DEBUG] Source ==> %+v\n", options.Source)
		cli.Printf(cli.Notice, "[DEBUG] Destination ==> %+v\n", options.Destination)
		cli.Printf(cli.Notice, "[DEBUG] Compress ==> %+v\n", options.Compress)
		cli.Printf(cli.Notice, "[DEBUG] SourceMap ==> %+v\n", options.SourceMap)
		cli.Printf(cli.Notice, "[DEBUG] Progressive ==> %+v\n", options.Progressive)
		cli.Printf(cli.Notice, "[DEBUG] Watch ==> %+v\n", context.WatchMode)
		cli.Printf(cli.Notice, "[DEBUG] Server ==> %+v\n", context.ServerMode)
		cli.Printf(cli.Notice, "[DEBUG] Server Port ==> %+v\n", context.ServerPort)

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
	if context.ServerMode {

		go func() {
			cli.Printf(cli.Notice, "[INFO] Starting server at \033[1m%s\033[0m\n", "http://localhost:"+context.ServerPort)
			server.Start(
				options.Destination.Path,
				context.ServerPort,
				func(uri string) error {
					cli.Printf(cli.Notice, "[GET] %s\n", uri)
					return nil
				},
			)
		}()

	}

	// Watch mode
	if context.WatchMode {

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
	if context.WatchMode || context.ServerMode {
		<-exit
	}

	// Shutdown
	shutdown(options)

}
