package main

import (
	"flag"
	"path/filepath"
	"strings"

	"github.com/mateussouzaweb/compactor/src/processor"
)

// Context struct
type Context struct {
	Version     bool
	DebugMode   bool
	WatchMode   bool
	ServerMode  bool
	ServerPort  string
	Source      string
	Destination string
	Options     *processor.Options
}

// trueOrFalse returns if given value is likely to be a true or false flag
func trueOrFalse(value string) bool {
	value = strings.ToLower(value)
	if value == "true" || value == "t" || value == "1" {
		return true
	}
	return false
}

// Read options from flags and arguments
func readContext() *Context {

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

	return &Context{
		Version:     version,
		DebugMode:   debug,
		WatchMode:   watch,
		ServerMode:  serverMode,
		ServerPort:  serverPort,
		Source:      source,
		Destination: destination,
		Options:     options,
	}
}
