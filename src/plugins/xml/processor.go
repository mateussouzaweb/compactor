package xml

import (
	"github.com/mateussouzaweb/compactor/src/plugins/generic"
	"github.com/mateussouzaweb/compactor/src/processor"
	"github.com/mateussouzaweb/compactor/src/system"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/xml"
)

// XML minify
func Minify(content string) (string, error) {

	m := minify.New()
	m.AddFunc("generic", xml.Minify)

	content, err := m.String("generic", content)

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
		Namespace:  "xml",
		Extensions: []string{".xml"},
		Init:       generic.Init,
		Shutdown:   generic.Shutdown,
		Resolve:    generic.Resolve,
		Related:    generic.Related,
		Transform:  generic.Transform,
		Optimize:   Optimize,
	}
}
