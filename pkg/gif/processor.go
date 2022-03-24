package gif

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// Init processor
func Init(bundle *compactor.Bundle) error {
	return os.NodeRequire("gifsicle", "gifsicle")
}

// Execute processor
func Execute(bundle *compactor.Bundle) error {

	for _, item := range bundle.Items {

		if !item.Exists {
			continue
		}

		destination := bundle.ToDestination(item.Path)
		err := os.Copy(item.Path, destination)

		if err != nil {
			return err
		}

		bundle.Processed(item.Path)

	}

	return nil
}

// Optimize processor
func Optimize(bundle *compactor.Bundle) error {

	for _, item := range bundle.Items {

		if !item.Exists || !bundle.ShouldCompress(item.Path) {
			continue
		}

		destination := bundle.ToDestination(item.Path)

		_, err := os.Exec(
			"gifsicle",
			"-03",
			destination,
			"-o", destination,
		)

		if err != nil {
			return err
		}

	}

	return nil
}

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:    "gif",
		Extensions:   []string{".gif"},
		Init:         Init,
		Dependencies: generic.Dependencies,
		Execute:      Execute,
		Optimize:     Optimize,
		Delete:       generic.Delete,
		Resolve:      generic.Resolve,
	}
}
