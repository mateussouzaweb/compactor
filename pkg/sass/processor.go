package sass

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/css"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// Init processor
func Init(bundle *compactor.Bundle) error {
	return os.NodeRequire("sass", "sass")
}

// Dependencies processor
func Dependencies(item *compactor.Item) ([]string, error) {

	// TODO: implement
	// item.Content

	return []string{}, nil
}

// Execute processor
func Execute(bundle *compactor.Bundle) error {

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

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:    "sass",
		Extensions:   []string{".sass", ".scss", ".css"},
		Init:         Init,
		Dependencies: Dependencies,
		Execute:      Execute,
		Optimize:     generic.Optimize,
		Delete:       css.Delete,
		Resolve:      css.Resolve,
	}
}
