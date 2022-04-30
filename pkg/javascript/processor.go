package javascript

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// Init processor
func Init(bundle *compactor.Bundle) error {
	return os.NodeRequire("terser", "terser")
}

// Related processor
func Related(item *compactor.Item) ([]compactor.Related, error) {

	var related []compactor.Related

	// Add possible source map
	extension := os.Extension(item.Path)
	file := strings.TrimSuffix(item.Path, extension)
	file = file + ".js.map"

	related = append(related, compactor.Related{
		Type:       "source-map",
		Dependency: true,
		Source:     "",
		Path:       os.File(file),
		Item:       compactor.Get(file),
	})

	// Add possible type declaration
	file = item.Path + ".d"
	related = append(related, compactor.Related{
		Type:       "declaration",
		Dependency: true,
		Source:     "",
		Path:       os.File(file),
		Item:       compactor.Get(file),
	})

	// Detect imports
	regex := regexp.MustCompile(`import (.+) from ("(.+)"|'(.+)');?`)
	matches := regex.FindAllStringSubmatch(item.Content, -1)

	for _, match := range matches {
		source := match[0]
		path := strings.Trim(match[2], `'"`)
		file := filepath.Join(os.Dir(item.Path), path)

		if os.Extension(file) == "" {
			file += item.Extension
		}

		if !os.Exist(file) {
			file = strings.TrimSuffix(file, os.Extension(file))
			file = file + item.Extension
		}

		related = append(related, compactor.Related{
			Type:       "import",
			Dependency: false,
			Source:     source,
			Path:       path,
			Item:       compactor.Get(file),
		})
	}

	return related, nil
}

// Execute processor
func Execute(bundle *compactor.Bundle) error {

	hash := bundle.Item.Checksum
	destination := bundle.ToDestination(bundle.Item.Path)
	destination = bundle.ToHashed(destination, hash)
	destination = bundle.ToExtension(destination, ".js")

	files := []string{bundle.Item.Path}

	for _, related := range bundle.Item.Related {
		if related.Item.Exists && related.Type == "import" {
			files = append(files, related.Item.Path)
		}
	}

	args := []string{}
	args = append(args, files...)
	args = append(args, "--output", destination)

	if bundle.ShouldCompress(bundle.Item.Path) {
		args = append(args, "--compress", "--comments")
	} else {
		args = append(args, "--beautify")
	}

	if bundle.ShouldGenerateSourceMap(bundle.Item.Path) {
		file := os.File(destination)
		args = append(args, "--source-map", strings.Join([]string{
			"includeSources",
			"base='" + bundle.Destination.Path + "'",
			"filename='" + file + ".map'",
			"url='" + file + ".map'",
		}, ","))
	}

	_, err := os.Exec(
		"terser",
		args...,
	)

	if err != nil {
		return err
	}

	return nil
}

// Resolve processor
func Resolve(path string) (string, error) {

	destination, err := generic.Resolve(path)

	if err != nil {
		return destination, err
	}

	bundle := compactor.GetBundle(path)
	hash := bundle.Item.Checksum

	destination = bundle.ToHashed(destination, hash)
	destination = bundle.ToExtension(destination, ".js")

	return destination, nil
}

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:  "javascript",
		Extensions: []string{".js", ".mjs"},
		Init:       Init,
		Related:    Related,
		Execute:    Execute,
		Optimize:   generic.Optimize,
		Delete:     generic.Delete,
		Resolve:    Resolve,
	}
}
