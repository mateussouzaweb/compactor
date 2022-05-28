package svg

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// Minify SVG content
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

// Optimize processor
func Optimize(options *compactor.Options, file *compactor.File) error {

	if !options.ShouldCompress(file.Path) {
		return nil
	}

	content := file.Content
	content, err := Minify(content)

	if err != nil {
		return err
	}

	destination := file.Destination
	perm := file.Permission
	err = os.Write(destination, content, perm)

	if err != nil {
		return err
	}

	return nil
}

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:  "svg",
		Extensions: []string{".svg"},
		Init:       generic.Init,
		Resolve:    generic.Resolve,
		Related:    generic.Related,
		Transform:  generic.Transform,
		Optimize:   Optimize,
	}
}
