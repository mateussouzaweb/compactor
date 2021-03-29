package compactor

import (
	"path/filepath"
	"strings"
)

// Extension struct
type Extension string

// Options struct
type Options struct {
	Source      string
	Destination string
	Watch       bool
	Minify      bool
	SourceMap   bool
	Compress    bool
	Progressive bool
	Include     []string
	Exclude     []string
	Maps        []string
}

// Context struct
type Context struct {
	File        string
	Extension   string
	Path        string
	Source      string
	Destination string
	Processed   bool
	Skipped     bool
}

// Processor struct
type Processor func(context *Context, options *Options) error

// Processors struct
type Processors []Processor

// Instance struct
type Instance map[Extension]Processors

// Instance var
var _processors = Instance{}

// Add processor to the instance
func Add(extension Extension, processor Processor) {

	if _, ok := _processors[extension]; !ok {
		_processors[extension] = Processors{}
	}

	_processors[extension] = append(_processors[extension], processor)

}

// Remove processor from the instance
func Remove(extension Extension) {
	if _, ok := _processors[extension]; ok {
		_processors[extension] = Processors{}
	}
}

// Process file
func Process(file string, options *Options) (*Context, error) {

	var err error
	var match bool

	context := &Context{
		File:        filepath.Base(file),
		Extension:   strings.TrimLeft(filepath.Ext(file), "."),
		Path:        strings.Replace(file, options.Source, "", 1),
		Source:      file,
		Destination: strings.Replace(file, options.Source, options.Destination, 1),
	}

	// Check include mode
	if !MatchList(context.Path, options.Include, true) {
		context.Skipped = true
		return context, err
	}

	// Check exclude mode
	if MatchList(context.Path, options.Exclude, false) {
		context.Skipped = true
		return context, err
	}

	// Make sure folder exists to avoid issues
	err = EnsureDirectory(context.Destination)

	if err != nil {
		return context, err
	}

	// Extension processors
	for extension, extensionProcessors := range _processors {
		if context.Extension == string(extension) && len(extensionProcessors) > 0 {

			for _, processor := range extensionProcessors {
				err = processor(context, options)
				if err != nil {
					return context, err
				}
			}

			match = true
			break
		}
	}

	// Generic processors
	if !match {
		for _, processor := range _processors["*"] {
			err = processor(context, options)
			if err != nil {
				return context, err
			}
		}
	}

	return context, err
}

// Check if filename match glob list pattern
func MatchList(filename string, list []string, defaultValue bool) bool {

	match := defaultValue

	for _, pattern := range list {

		result, err := filepath.Match(pattern, filename)

		if err != nil {
			continue
		}

		if result {
			match = true
			break
		}

	}

	return match
}
