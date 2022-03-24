package svg

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// Init processor
func Init(bundle *compactor.Bundle) error {
	return os.NodeRequire("svgo", "svgo")
}

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
func Optimize(bundle *compactor.Bundle) error {

	destination := bundle.ToDestination(bundle.Destination.File)

	if !bundle.ShouldCompress(destination) {
		return nil
	}

	content := bundle.GetContent()
	content, err := Minify(content)

	if err != nil {
		return err
	}

	perm := bundle.GetPermission()
	err = os.Write(destination, content, perm)

	return err
}

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:    "svg",
		Extensions:   []string{".svg"},
		Init:         Init,
		Dependencies: generic.Dependencies,
		Execute:      generic.Execute,
		Optimize:     Optimize,
		Delete:       generic.Delete,
		Resolve:      generic.Resolve,
	}
}
