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
		Namespace:    "xml",
		Extensions:   []string{".xml"},
		Init:         generic.Init,
		Dependencies: generic.Dependencies,
		Execute:      generic.Execute,
		Optimize:     Optimize,
		Delete:       generic.Delete,
		Resolve:      generic.Resolve,
	}
}
