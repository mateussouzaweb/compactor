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
func Optimize(bundle *compactor.Bundle) error {

	if !bundle.ShouldCompress(bundle.Item.Path) {
		return nil
	}

	content := bundle.Item.Content
	content, err := Minify(content)

	if err != nil {
		return err
	}

	destination := bundle.ToDestination(bundle.Item.Path)
	perm := bundle.Item.Permission
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
		Related:    generic.Related,
		Resolve:    generic.Resolve,
		Execute:    generic.Execute,
		Optimize:   Optimize,
		Delete:     generic.Delete,
	}
}
