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

// XML processor
func RunProcessor(bundle *compactor.Bundle) error {

	// TODO: to multiple, merge xmls as array and join data
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

func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Extensions: []string{".xml"},
		Run:        RunProcessor,
		Delete:     generic.DeleteProcessor,
		Resolve:    generic.ResolveProcessor,
	}
}
