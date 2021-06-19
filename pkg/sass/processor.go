package sass

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/css"
)

// Sass processor
func RunProcessor(bundle *compactor.Bundle) error {

	// TODO: to multiple, simulate a sass file with @imports
	for _, item := range bundle.Items {

		if !item.Exists {
			continue
		}

		destination := bundle.ToDestination(item.Path)
		destination = bundle.ToHashed(destination, item.Checksum)
		destination = bundle.ToExtension(destination, ".css")

		args := []string{
			item.Path + ":" + destination,
		}

		if bundle.ShouldCompress(item.Path) {
			args = append(args, "--style", "compressed")
		}

		if bundle.ShouldGenerateSourceMap(item.Path) {
			args = append(args, "--source-map", "--embed-sources")
		}

		_, err := os.Exec(
			"sass",
			args...,
		)

		if err != nil {
			return err
		}

		bundle.Processed(item.Path)

	}

	return nil
}

func Plugin() *compactor.Plugin {

	os.NodeRequire("sass", "sass")

	return &compactor.Plugin{
		Extensions: []string{".sass", ".scss"},
		Run:        RunProcessor,
		Delete:     css.DeleteProcessor,
		Resolve:    css.ResolveProcessor,
	}
}
