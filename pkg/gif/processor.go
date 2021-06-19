package gif

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// GIF processor
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
				"gifsicle",
				"-03",
				destination,
				"-o", destination,
			)
		}

		if err != nil {
			return err
		}

		bundle.Processed(item.Path)

	}

	return nil
}

func Plugin() *compactor.Plugin {

	os.NodeRequire("gifsicle", "gifsicle")

	return &compactor.Plugin{
		Extensions: []string{".gif"},
		Run:        RunProcessor,
		Delete:     generic.DeleteProcessor,
		Resolve:    generic.ResolveProcessor,
	}
}
