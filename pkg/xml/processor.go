package xml

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/xml"
)

// XML minify
func Minify(content string) (string, error) {

	m := minify.New()
	m.AddFunc("generic", xml.Minify)

	content, err := m.String("generic", content)

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
