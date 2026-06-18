package generic

import (
	"github.com/mateussouzaweb/compactor/src/processor"
	"github.com/mateussouzaweb/compactor/src/system"
)

// Init processor
func Init(options *processor.Options) error {
	return nil
}

// Shutdown processor
func Shutdown(options *processor.Options) error {
	return nil
}

// Resolve returns the clean file destination path for given source file
func Resolve(options *processor.Options, file *processor.File) (string, error) {
	destination := options.ToDestination(file.Path)
	return destination, nil
}

// Related detects the dependencies of the file
func Related(options *processor.Options, file *processor.File) ([]processor.Related, error) {
	var found []processor.Related
	return found, nil
}

// Transform creates generic copy of file(s) content to destination
func Transform(options *processor.Options, file *processor.File) error {

	content := file.Content
	perm := file.Permission
	destination := file.Destination

	err := system.Write(destination, content, perm)
	if err != nil {
		return err
	}

	return nil
}

// Optimize apply optimizations into the destination file
func Optimize(options *processor.Options, file *processor.File) error {
	return nil
}

// Plugin return the compactor plugin instance
func Plugin() *processor.Plugin {
	return &processor.Plugin{
		Namespace:  "generic",
		Extensions: []string{},
		Init:       Init,
		Shutdown:   Shutdown,
		Resolve:    Resolve,
		Related:    Related,
		Transform:  Transform,
		Optimize:   Optimize,
	}
}
