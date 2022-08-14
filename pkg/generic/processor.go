package generic

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
)

// Init processor
func Init(options *compactor.Options) error {
	return nil
}

// Shutdown processor
func Shutdown(options *compactor.Options) error {
	return nil
}

// Resolve returns the clean file destination path for given source file
func Resolve(options *compactor.Options, file *compactor.File) (string, error) {
	destination := options.ToDestination(file.Path)
	return destination, nil
}

// Related detects the dependencies of the file
func Related(options *compactor.Options, file *compactor.File) ([]compactor.Related, error) {
	var found []compactor.Related
	return found, nil
}

// Transform creates generic copy of file(s) content to destination
func Transform(options *compactor.Options, file *compactor.File) error {

	content := file.Content
	perm := file.Permission
	destination := file.Destination
	err := os.Write(destination, content, perm)

	if err != nil {
		return err
	}

	return nil
}

// Optimize apply optimizations into the destination file
func Optimize(options *compactor.Options, file *compactor.File) error {
	return nil
}

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
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
