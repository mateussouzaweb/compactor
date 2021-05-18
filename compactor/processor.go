package compactor

// Processor struct
type Processor func(action *Action, bundle *Bundle, logger *Logger) error

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
	bundle.Extension = ""
	bundle.Destination.File = ""
	bundle.Source.Files = []string{}

	return &bundle
}

// RetrieveBundles get a list of all registered bundles
func RetrieveBundles() Bundles {
	return _bundles
}

// RetrieveBundleFor retrieve the related bundle of the file
func RetrieveBundleFor(file string) *Bundle {

	for _, bundle := range _bundles {
		if bundle.ContainsFile(file) {
			return bundle
		}
	}

	bundle := NewBundle()
	bundle.Extension = bundle.CleanExtension(file)
	bundle.Destination.File = bundle.CleanPath(file)
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

	logger := Logger{}
	action := Action{Type: "PROCESS"}

	// Determine action based on processable list
	files := bundle.GetFiles()
	if len(files) == 0 {
		action.Type = "DELETE"
	}

	// Make sure folder exists to avoid issues
	destination := bundle.GetDestination()
	err := EnsureDirectory(destination)

	if err != nil {
		return logger, err
	}

	// Process by extension
	processors := RetrieveProcessors(bundle.Extension)

	if len(processors) == 0 {
		processors = RetrieveProcessors("*")
	}

	// Execute processors
	for _, callback := range processors {
		err = callback(&action, bundle, &logger)
		if err != nil {
			return logger, err
		}
	}

	return logger, err
}
