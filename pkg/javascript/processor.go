package javascript

import (
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// Init processor
func Init(bundle *compactor.Bundle) error {
	return os.NodeRequire("terser", "terser")
}

// Dependencies processor
func Dependencies(item *compactor.Item) ([]string, error) {

	// TODO: implement
	// item.Content

	return []string{}, nil
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
	destination = bundle.ToExtension(destination, ".js")

	files := []string{}

	for _, item := range bundle.Items {
		if item.Exists {
			files = append(files, item.Path)
		}
	}

	args := []string{}
	args = append(args, files...)
	args = append(args, "--output", destination)

	if bundle.ShouldCompress(destination) {
		args = append(args, "--compress", "--comments")
	} else {
		args = append(args, "--beautify")
	}

	if bundle.ShouldGenerateSourceMap(destination) {
		file := os.File(destination)
		args = append(args, "--source-map", strings.Join([]string{
			"includeSources",
			"base='" + bundle.Destination.Path + "'",
			"filename='" + file + ".map'",
			"url='" + file + ".map'",
		}, ","))
	}

	_, err = os.Exec(
		"terser",
		args...,
	)

	if err == nil {
		bundle.Written(destination)
	}

	return err
}

// Delete processor
func Delete(bundle *compactor.Bundle) error {

	err := generic.Delete(bundle)

	if err != nil {
		return err
	}

	for _, deleted := range bundle.Logs.Deleted {

		extra := bundle.ToExtension(deleted, ".js.map")

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
	destination = bundle.ToExtension(destination, ".js")

	return destination, nil
}

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:    "javascript",
		Extensions:   []string{".js", ".mjs"},
		Init:         Init,
		Dependencies: Dependencies,
		Execute:      Execute,
		Optimize:     generic.Optimize,
		Delete:       Delete,
		Resolve:      Resolve,
	}
}
