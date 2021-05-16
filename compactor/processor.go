package compactor

import (
	"path/filepath"
)

// Processor struct
type Processor func(files []string, bundle *Bundle, logger *Logger) error

// Processors struct
type Processors []Processor

// ProcessorsMap struct
type ProcessorsMap map[string]Processors

// Bundles struct
type Bundles []*Bundle

// Default bundle used as model
var Default = &Bundle{}

// Variables
var _processors = ProcessorsMap{}
var _bundles = Bundles{}

// RegisterProcessor register a new processor for the extension
func RegisterProcessor(extension string, processor Processor) {

	if _, ok := _processors[extension]; !ok {
		_processors[extension] = Processors{}
	}

	_processors[extension] = append(_processors[extension], processor)

}

// RemoveProcessors removes all processors for the extension
func RemoveProcessors(extension string) {

	if _, ok := _processors[extension]; ok {
		_processors[extension] = Processors{}
	}

}

// RetrieveProcessors for given extension
func RetrieveProcessors(extension string) Processors {

	if _, ok := _processors[extension]; ok {
		return _processors[extension]
	}

	return Processors{}
}

// NewBundle create a new bundle instance from default bundle
func NewBundle() *Bundle {

	bundle := *Default
	bundle.Files = []string{}

	return &bundle
}

// RetrieveBundles get a list of all registered bundles
func RetrieveBundles() Bundles {
	return _bundles
}

// RetrieveBundleFor retrieve the related bundle of the file
func RetrieveBundleFor(file string) *Bundle {

	for _, bundle := range _bundles {

		compare := bundle.CleanPath(file)
		for _, fileInPackage := range bundle.Files {

			if fileInPackage == compare {
				return bundle
			}

			match, err := filepath.Match(compare, fileInPackage)

			if err != nil {
				continue
			}
			if match {
				return bundle
			}

		}

	}

	bundle := NewBundle()
	bundle.Target = bundle.CleanPath(file)
	bundle.AddFile(file)

	RegisterBundle(bundle)

	return bundle
}

// RegisterBundle register a bundle into the index
func RegisterBundle(bundle *Bundle) {
	_bundles = append(_bundles, bundle)
}

// Process package by running processors
func Process(bundle *Bundle) (Logger, error) {

	// Make sure folder exists to avoid issues
	err := EnsureDirectory(bundle.DestinationPath(bundle.Target))
	logger := Logger{}

	if err != nil {
		return logger, err
	}

	// Retrieve processable file list and do basic logging on files
	files := []string{}
	for _, file := range bundle.GetFiles() {
		if bundle.ShouldSkip(file) {
			logger.AddSkipped(file)
		} else if bundle.ShouldIgnore(file) {
			logger.AddIgnored(file)
		} else {
			files = append(files, file)
		}
	}

	if len(files) == 0 {
		return logger, nil
	}

	processors := RetrieveProcessors(bundle.CleanExtension(files[0]))

	if len(processors) == 0 {
		processors = RetrieveProcessors("*")
	}

	// Extension processors
	for _, callback := range processors {
		err = callback(files, bundle, &logger)
		if err != nil {
			return logger, err
		}
	}

	return logger, err
}
