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

// Related processor
func Related(item *compactor.Item) ([]compactor.Related, error) {
	var found []compactor.Related
	return found, nil
}

// Execute processor
func Execute(bundle *compactor.Bundle) error {

	hash, err := bundle.GetChecksum()

	if err != nil {
		return err
	}

	destination := bundle.ToDestination(bundle.Item.Path)
	destination = bundle.ToHashed(destination, hash)
	destination = bundle.ToExtension(destination, ".js")

	files := []string{bundle.Item.Path}

	for _, related := range bundle.Item.Related {
		if related.Item.Exists && related.Type == "import" {
			files = append(files, related.Item.Path)
		}
		if related.Item.Exists && related.Type == "require" {
			files = append(files, related.Item.Path)
		}
	}

	args := []string{}
	args = append(args, files...)
	args = append(args, "--output", destination)

	if bundle.ShouldCompress(bundle.Item.Path) {
		args = append(args, "--compress", "--comments")
	} else {
		args = append(args, "--beautify")
	}

	if bundle.ShouldGenerateSourceMap(bundle.Item.Path) {
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

	if err != nil {
		return err
	}

	bundle.Written(destination)

	if bundle.ShouldCompress(bundle.Item.Path) {
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
	hash, err := bundle.GetChecksum()

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
		Namespace:  "javascript",
		Extensions: []string{".js", ".mjs"},
		Init:       Init,
		Related:    Related,
		Execute:    Execute,
		Optimize:   generic.Optimize,
		Delete:     Delete,
		Resolve:    Resolve,
	}
}
