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

	destination := bundle.ToDestination(bundle.Item.Path)
	err := os.Copy(bundle.Item.Path, destination)

	if err != nil {
		return err
	}

	return nil
}

// Optimize processor
func Optimize(bundle *compactor.Bundle) error {

	if !bundle.ShouldCompress(bundle.Item.Path) {
		return nil
	}

	destination := bundle.ToDestination(bundle.Item.Path)
	_, err := os.Exec(
		"gifsicle",
		"-03",
		destination,
		"-o", destination,
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
		Init:       Init,
		Related:    generic.Related,
		Execute:    Execute,
		Optimize:   Optimize,
		Delete:     generic.Delete,
		Resolve:    generic.Resolve,
	}
}
