package jpeg

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
	"github.com/mateussouzaweb/compactor/pkg/webp"
)

// Init processor
func Init(bundle *compactor.Bundle) error {

	err := os.NodeRequire("jpegoptim", "jpegoptim-bin")

	if err != nil {
		return err
	}

	return os.NodeRequire("cwebp", "cwebp-bin")
}

// Related processor
func Related(item *compactor.Item) ([]compactor.Related, error) {

	var related []compactor.Related

	related = append(related, compactor.Related{
		Type:       "alternative",
		Dependency: true,
		Source:     "",
		Path:       item.File + ".webp",
		Item:       compactor.Get(item.Path + ".webp"),
	})

	return related, nil
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

	if bundle.ShouldCompress(bundle.Item.Path) {
		destination := bundle.ToDestination(bundle.Item.Path)
		_, err := os.Exec(
			"jpegoptim",
			"--quiet",
			"--strip-all",
			"--all-progressive",
			"--overwrite",
			destination,
		)

		if err != nil {
			return err
		}
	}

	if bundle.ShouldGenerateProgressive(bundle.Item.Path) {
		destination := bundle.ToDestination(bundle.Item.Path)
		err := webp.CreateCopy(bundle.Item.Path, destination, 75)

		if err != nil {
			return err
		}
	}

	return nil
}

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:  "jpeg",
		Extensions: []string{".jpeg", ".jpg"},
		Init:       Init,
		Related:    Related,
		Execute:    Execute,
		Optimize:   Optimize,
		Delete:     generic.Delete,
		Resolve:    generic.Resolve,
	}
}
