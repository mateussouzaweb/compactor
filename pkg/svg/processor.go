package svg

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// Init processor
func InitProcessor(bundle *compactor.Bundle) error {
	return os.NodeRequire("svgo", "svgo")
}

// SVG minify
func Minify(content string) (string, error) {

	config := os.TemporaryFile("svgo.config.js")
	file := os.TemporaryFile("svgo.svg")
	defer os.Delete(config)
	defer os.Delete(file)

	settings := "module.exports = {plugins: [{ name: 'removeViewBox', active: false }]}"
	err := os.Write(config, settings, 0775)

	if err != nil {
		return content, err
	}

	err = os.Write(file, content, 0775)

	if err != nil {
		return content, err
	}

	_, err = os.Exec(
		"svgo",
		"--quiet",
		"--config", config,
		"--input", file,
		"--output", file,
	)

	if err != nil {
		return content, err
	}

	content, err = os.Read(file)

	return content, err
}

// SVG processor
func RunProcessor(bundle *compactor.Bundle) error {

	// TODO: to multiple, merge svgs as array and join data
	if bundle.ShouldOutputToMany() {

		for _, item := range bundle.Items {

			if !item.Exists {
				continue
			}

			content := item.Content
			var err error

			if bundle.ShouldCompress(item.Path) {
				content, err = Minify(content)
				if err != nil {
					return err
				}
			}

			destination := bundle.ToDestination(item.Path)
			err = os.Write(destination, content, item.Permission)

			if err != nil {
				return err
			}

			bundle.Processed(item.Path)

		}

		return nil
	}

	content := ""
	for _, item := range bundle.Items {
		if item.Exists {
			content += item.Content
		}
	}

	destination := bundle.ToDestination(bundle.Destination.File)
	var err error

	if bundle.ShouldCompress(destination) {
		content, err = Minify(content)
		if err != nil {
			return err
		}
	}

	perm := bundle.Items[0].Permission
	err = os.Write(destination, content, perm)

	if err == nil {
		bundle.Written(destination)
	}

	return err
}

func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Extensions: []string{".svg"},
		Init:       InitProcessor,
		Run:        RunProcessor,
		Delete:     generic.DeleteProcessor,
		Resolve:    generic.ResolveProcessor,
	}
}
