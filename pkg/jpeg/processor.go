package jpeg

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
	"github.com/mateussouzaweb/compactor/pkg/webp"
)

// Init processor
func InitProcessor(bundle *compactor.Bundle) error {

	err := os.NodeRequire("jpegoptim", "jpegoptim-bin")

	if err != nil {
		return err
	}

	return os.NodeRequire("cwebp", "cwebp-bin")
}

// JPEG processor
func RunProcessor(bundle *compactor.Bundle) error {

	for _, item := range bundle.Items {

		if !item.Exists {
			continue
		}

		destination := bundle.ToDestination(item.Path)
		err := os.Copy(item.Path, destination)

		if err != nil {
			return err
		}

		if bundle.ShouldCompress(item.Path) {
			_, err = os.Exec(
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

		if bundle.ShouldGenerateProgressive(item.Path) {
			err = webp.CreateCopy(item.Path, destination, 75)
		}

		if err != nil {
			return err
		}

		bundle.Processed(item.Path)

	}

	return nil
}

// DeleteProcessor
func DeleteProcessor(bundle *compactor.Bundle) error {

	err := generic.DeleteProcessor(bundle)

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

func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Extensions: []string{".jpeg", ".jpg"},
		Init:       InitProcessor,
		Run:        RunProcessor,
		Delete:     DeleteProcessor,
		Resolve:    generic.ResolveProcessor,
	}
}
