package svg

import (
	"io/ioutil"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// SVG minify
func Minify(content string) (string, error) {

	config, err := ioutil.TempFile("", "svgo.*.config.js")
	defer os.Delete(config.Name())

	if err != nil {
		return content, err
	}

	_, err = config.WriteString("module.exports = {plugins: [{ name: 'removeViewBox', active: false }]}")

	if err != nil {
		return content, err
	}

	file, err := ioutil.TempFile("", "svgo.*.svg")
	defer os.Delete(file.Name())

	if err != nil {
		return content, err
	}

	_, err = file.WriteString(content)

	if err != nil {
		return content, err
	}

	_, err = os.Exec(
		"svgo",
		"--quiet",
		"--config", config.Name(),
		"--input", file.Name(),
		"--output", file.Name(),
	)

	if err != nil {
		return content, err
	}

	content, err = os.Read(file.Name())

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

	os.NodeRequire("svgo", "svgo")

	return &compactor.Plugin{
		Extensions: []string{".svg"},
		Run:        RunProcessor,
		Delete:     generic.DeleteProcessor,
		Resolve:    generic.ResolveProcessor,
	}
}
