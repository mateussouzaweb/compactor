package svg

import (
	"github.com/mateussouzaweb/compactor/src/plugins/generic"
	"github.com/mateussouzaweb/compactor/src/processor"
	"github.com/mateussouzaweb/compactor/src/system"
)

// Minify SVG content
func Minify(content string) (string, error) {

	config := system.TemporaryFile("svgo.config.js")
	file := system.TemporaryFile("svgo.svg")
	defer system.Delete(config)
	defer system.Delete(file)

	settings := "module.exports = {plugins: [{ name: 'removeViewBox', active: false }]}"
	err := system.Write(config, settings, 0775)
	if err != nil {
		return content, err
	}

	err = system.Write(file, content, 0775)
	if err != nil {
		return content, err
	}

	_, err = system.Exec(
		"svgo",
		"--quiet",
		"--config", config,
		"--input", file,
		"--output", file,
	)

	if err != nil {
		return content, err
	}

	content, err = system.Read(file)

	return content, err
}

// Optimize processor
func Optimize(options *processor.Options, file *processor.File) error {

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
	err = system.Write(destination, content, perm)
	if err != nil {
		return err
	}

	return nil
}

// Plugin return the compactor plugin instance
func Plugin() *processor.Plugin {
	return &processor.Plugin{
		Namespace:  "svg",
		Extensions: []string{".svg"},
		Init:       generic.Init,
		Shutdown:   generic.Shutdown,
		Resolve:    generic.Resolve,
		Related:    generic.Related,
		Transform:  generic.Transform,
		Optimize:   Optimize,
	}
}
