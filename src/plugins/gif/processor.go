package gif

import (
	"github.com/mateussouzaweb/compactor/src/plugins/generic"
	"github.com/mateussouzaweb/compactor/src/processor"
	"github.com/mateussouzaweb/compactor/src/system"
)

// Transform processor
func Transform(options *processor.Options, file *processor.File) error {

	err := system.Copy(file.Path, file.Destination)
	if err != nil {
		return err
	}

	return nil
}

// Optimize processor
func Optimize(options *processor.Options, file *processor.File) error {

	if !options.ShouldCompress(file.Path) {
		return nil
	}

	_, err := system.Exec(
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
func Plugin() *processor.Plugin {
	return &processor.Plugin{
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
