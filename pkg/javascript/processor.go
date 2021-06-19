package javascript

import (
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// Javascript processor
func RunProcessor(bundle *compactor.Bundle) error {

	if bundle.ShouldOutputToMany() {

		for _, item := range bundle.Items {

			if !item.Exists {
				continue
			}

			destination := bundle.ToDestination(item.Path)
			destination = bundle.ToHashed(destination, item.Checksum)
			destination = bundle.ToExtension(destination, ".js")

			args := []string{}
			args = append(args, item.Path)
			args = append(args, "--output", destination)

			if bundle.ShouldCompress(item.Path) {
				args = append(args, "--compress", "--comments")
			} else {
				args = append(args, "--beautify")
			}

			if bundle.ShouldGenerateSourceMap(item.Path) {
				args = append(args, "--source-map", strings.Join([]string{
					"includeSources",
					"filename='" + item.File + ".map'",
					"url='" + item.File + ".map'",
				}, ","))
			}

			_, err := os.Exec(
				"uglifyjs",
				args...,
			)

			if err != nil {
				return err
			}

			bundle.Processed(item.Path)

		}

		return nil
	}

	content := ""
	files := []string{}

	for _, item := range bundle.Items {
		if item.Exists {
			content += item.Content
			files = append(files, item.Path)
		}
	}

	hash, err := os.Checksum(content)

	if err != nil {
		return err
	}

	destination := bundle.ToDestination(bundle.Destination.File)
	destination = bundle.ToHashed(destination, hash)
	destination = bundle.ToExtension(destination, ".js")

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
			"filename='" + file + ".map'",
			"url='" + file + ".map'",
		}, ","))
	}

	_, err = os.Exec(
		"uglifyjs",
		args...,
	)

	if err == nil {
		bundle.Written(destination)
	}

	return err
}

// DeleteProcessor
func DeleteProcessor(bundle *compactor.Bundle) error {

	err := generic.DeleteProcessor(bundle)

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

// ResolveProcessor fix the path for given file path
func ResolveProcessor(path string) (string, error) {

	destination, err := generic.ResolveProcessor(path)

	if err != nil {
		return destination, err
	}

	bundle := compactor.GetBundleFor(path)

	if bundle.ShouldOutputToMany() {

		source := bundle.ToSource(path)
		item := compactor.Get(source)

		destination := bundle.ToHashed(path, item.Checksum)
		destination = bundle.ToExtension(destination, ".js")

		return destination, nil
	}

	content := ""
	for _, item := range bundle.Items {
		if item.Exists {
			content += item.Content
		}
	}

	hash, err := os.Checksum(content)

	if err != nil {
		return destination, err
	}

	destination = bundle.ToHashed(destination, hash)
	destination = bundle.ToExtension(destination, ".js")

	return destination, nil
}

func Plugin() *compactor.Plugin {

	os.NodeRequire("uglifyjs", "uglify-js")

	return &compactor.Plugin{
		Extensions: []string{".js", ".mjs"},
		Run:        RunProcessor,
		Delete:     DeleteProcessor,
		Resolve:    ResolveProcessor,
	}
}
