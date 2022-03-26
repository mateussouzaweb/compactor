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

// Execute processor
func Execute(bundle *compactor.Bundle) error {

	destination := bundle.ToDestination(bundle.Item.Path)
	err := os.Copy(bundle.Item.Path, destination)

	if err != nil {
		return err
	}

	bundle.Processed(bundle.Item.Path)

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

	if bundle.ShouldCompress(bundle.Item.Path) || bundle.ShouldGenerateProgressive(bundle.Item.Path) {
		bundle.Optimized(bundle.Item.Path)
	}

	return nil
}

// Delete processor
func Delete(bundle *compactor.Bundle) error {

	err := generic.Delete(bundle)

	if err != nil {
		return err
	}

	for _, deleted := range bundle.Logs.Deleted {

		extension := os.Extension(deleted)
		extra := bundle.ToExtension(deleted, extension+".webp")

		if !os.Exist(extra) {
			continue
		}

		err := os.Delete(extra)
		if err != nil {
			return err
		}

	}

	return err
}

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:  "jpeg",
		Extensions: []string{".jpeg", ".jpg"},
		Init:       Init,
		Related:    generic.Related,
		Execute:    Execute,
		Optimize:   Optimize,
		Delete:     Delete,
		Resolve:    generic.Resolve,
	}
}
