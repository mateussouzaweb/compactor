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

// Execute processor
func Execute(bundle *compactor.Bundle) error {

	hash := bundle.Item.Checksum
	destination := bundle.ToDestination(bundle.Item.Path)
	destination = bundle.ToHashed(destination, hash)
	destination = bundle.ToExtension(destination, ".css")

	args := []string{
		bundle.Item.Path + ":" + destination,
	}

	if bundle.ShouldCompress(bundle.Item.Path) {
		args = append(args, "--style", "compressed")
	}

	if bundle.ShouldGenerateSourceMap(bundle.Item.Path) {
		args = append(args, "--source-map", "--embed-sources")
	}

	_, err := os.Exec(
		"sass",
		args...,
	)

	if err != nil {
		return err
	}

	return nil
}

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:  "sass",
		Extensions: []string{".sass", ".scss", ".css"},
		Init:       Init,
		Related:    css.Related,
		Execute:    Execute,
		Optimize:   generic.Optimize,
		Delete:     generic.Delete,
		Resolve:    css.Resolve,
	}
}
