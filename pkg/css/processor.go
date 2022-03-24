package css

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// Init processor
func Init(bundle *compactor.Bundle) error {
	return os.NodeRequire("sass", "sass")
}

// Execute processor
func Execute(bundle *compactor.Bundle) error {

	content := bundle.GetContent()
	hash, err := os.Checksum(content)

	if err != nil {
		return err
	}

	destination := bundle.ToDestination(bundle.Destination.File)
	destination = bundle.ToHashed(destination, hash)
	destination = bundle.ToExtension(destination, ".css")

	perm := bundle.GetPermission()
	err = os.Write(destination, content, perm)

	if err == nil {
		bundle.Processed(destination)
	}

	args := []string{
		destination + ":" + destination,
	}

	if bundle.ShouldCompress(destination) {
		args = append(args, "--style", "compressed")
	}

	if bundle.ShouldGenerateSourceMap(destination) {
		args = append(args, "--source-map", "--embed-sources")
	}

	_, err = os.Exec(
		"sass",
		args...,
	)

	return err
}

// Delete processor
func Delete(bundle *compactor.Bundle) error {

	err := generic.Delete(bundle)

	if err != nil {
		return err
	}

	for _, deleted := range bundle.Logs.Deleted {

		extra := bundle.ToExtension(deleted, ".css.map")

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

// Resolve processor
func Resolve(path string) (string, error) {

	destination, err := generic.Resolve(path)

	if err != nil {
		return destination, err
	}

	bundle := compactor.GetBundle(path)
	content := bundle.GetContent()

	hash, err := os.Checksum(content)

	if err != nil {
		return destination, err
	}

	destination = bundle.ToHashed(destination, hash)
	destination = bundle.ToExtension(destination, ".css")

	return destination, nil
}

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:    "css",
		Extensions:   []string{".css"},
		Init:         Init,
		Dependencies: generic.Dependencies,
		Execute:      Execute,
		Optimize:     generic.Optimize,
		Delete:       Delete,
		Resolve:      Resolve,
	}
}
