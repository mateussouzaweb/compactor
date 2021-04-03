package compactor

import (
	"path/filepath"
	"strings"
)

// Extension struct
type Extension string

// Context struct
type Context struct {
	File        string
	Extension   string
	Path        string
	Source      string
	Destination string
	Processed   bool
	Skipped     bool
	Ignored     bool
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

	if options.ShouldIgnore(context) {
		context.Ignored = true
		return context, err
	}

	if options.ShouldSkip(context) {
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
