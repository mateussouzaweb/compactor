package css

import (
	"regexp"
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// Init processor
func Init(bundle *compactor.Bundle) error {
	return os.NodeRequire("sass", "sass")
}

// Related processor
func Related(item *compactor.Item) ([]compactor.Related, error) {

	var related []compactor.Related

	// Add possible source map
	extension := os.Extension(item.Path)
	file := strings.TrimSuffix(item.Path, extension)
	file = file + ".css.map"

	related = append(related, compactor.Related{
		Type:       "source-map",
		Dependency: true,
		Source:     "",
		Path:       os.File(file),
		Item:       compactor.Get(file),
	})

	// Detect imports
	regex := regexp.MustCompile(`@import ("(.+)"|'(.+)');?`)
	matches := regex.FindAllStringSubmatch(item.Content, -1)

	for _, match := range matches {
		source := match[0]
		sourcePath := strings.Trim(match[1], `'"`)

		path := sourcePath
		if os.Extension(path) == "" {
			path += item.Extension
		}

		file := os.Resolve(path, os.Dir(item.Path))
		related = append(related, compactor.Related{
			Type:       "import",
			Dependency: true,
			Source:     source,
			Path:       sourcePath,
			Item:       compactor.Get(file),
		})
	}

	return related, nil
}

// Resolve processor
func Resolve(path string, item *compactor.Item) (string, error) {

	destination, err := generic.Resolve(path, item)

	if err != nil {
		return destination, err
	}

	bundle := compactor.GetBundle(path)
	hash := bundle.Item.Checksum

	destination = bundle.ToHashed(destination, hash)
	destination = bundle.ToExtension(destination, ".css")

	return destination, nil
}

// Execute processor
func Execute(bundle *compactor.Bundle) error {

	content := bundle.Item.Content
	hash := bundle.Item.Checksum

	destination := bundle.ToDestination(bundle.Item.Path)
	destination = bundle.ToHashed(destination, hash)
	destination = bundle.ToExtension(destination, ".css")

	perm := bundle.Item.Permission
	err := os.Write(destination, content, perm)

	if err != nil {
		return err
	}

	args := []string{
		destination + ":" + destination,
	}

	if bundle.ShouldCompress(bundle.Item.Path) {
		args = append(args, "--style", "compressed")
	}

	if bundle.ShouldGenerateSourceMap(bundle.Item.Path) {
		args = append(args, "--source-map", "--embed-sources")
	}

	_, err = os.Exec(
		"sass",
		args...,
	)

	if err != nil {
		return err
	}

	return nil
}

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:  "css",
		Extensions: []string{".css"},
		Init:       Init,
		Related:    Related,
		Resolve:    Resolve,
		Execute:    Execute,
		Optimize:   generic.Optimize,
		Delete:     generic.Delete,
	}
}
