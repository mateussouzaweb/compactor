package svg

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// SVG minify
func Minify(content string) (string, error) {

	// TODO: Viewbox removal causing bugs
	// _, err = os.Exec(
	// 	"svgo",
	// 	"--quiet",
	// 	"--input", target,
	// 	"--output", target,
	// )

	return content, nil
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
		Run:        RunProcessor,
		Delete:     generic.DeleteProcessor,
		Resolve:    generic.ResolveProcessor,
	}
}
