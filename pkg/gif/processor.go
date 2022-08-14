package gif

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// Transform processor
func Transform(options *compactor.Options, file *compactor.File) error {

	err := os.Copy(file.Path, file.Destination)

	if err != nil {
		return err
	}

	return nil
}

// Optimize processor
func Optimize(options *compactor.Options, file *compactor.File) error {

	if !options.ShouldCompress(file.Path) {
		return nil
	}

	_, err := os.Exec(
		"gifsicle",
		"-03",
		file.Destination,
		"-o", file.Destination,
	)

	if err != nil {
		return err
	}

	return nil
}

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:  "gif",
		Extensions: []string{".gif"},
		Init:       generic.Init,
		Shutdown:   generic.Shutdown,
		Resolve:    generic.Resolve,
		Related:    generic.Related,
		Transform:  Transform,
		Optimize:   Optimize,
	}
}
