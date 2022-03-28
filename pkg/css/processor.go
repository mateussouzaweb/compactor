package css

import (
	"path/filepath"
	"regexp"

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

	var found []compactor.Related

	pattern := regexp.MustCompile("@import \"(.+)\";?")
	matches := pattern.FindAllStringSubmatch(item.Content, -1)

	for _, match := range matches {
		source := match[0]
		file := filepath.Join(os.Dir(item.Path), match[1])

		if os.Extension(file) == "" {
			file += item.Extension
		}

		found = append(found, compactor.Related{
			Type:   "import",
			Source: source,
			Item:   compactor.Get(file),
		})
	}

	return found, nil
}

// Execute processor
func Execute(bundle *compactor.Bundle) error {

	content := bundle.GetContent()
	hash, err := bundle.GetChecksum()

	if err != nil {
		return err
	}

	destination := bundle.ToDestination(bundle.Item.Path)
	destination = bundle.ToHashed(destination, hash)
	destination = bundle.ToExtension(destination, ".css")

	perm := bundle.GetPermission()
	err = os.Write(destination, content, perm)

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

	bundle.Processed(bundle.Item.Path)

	if bundle.ShouldCompress(bundle.Item.Path) {
		bundle.Optimized(bundle.Item.Path)
	}

	return nil
}

// Delete processor
func Delete(bundle *compactor.Bundle) error {

	err := generic.Delete(bundle)

	if err != nil {
		return err
	}

	for _, deleted := range bundle.Logs.Deleted {

		extra := bundle.ToExtension(deleted, ".css.map")

		if !os.Exist(extra) {
			continue
		}

		err := os.Delete(extra)
		if err != nil {
			return err
		}

	}

	return err
}

// Resolve processor
func Resolve(path string) (string, error) {

	destination, err := generic.Resolve(path)

	if err != nil {
		return destination, err
	}

	bundle := compactor.GetBundle(path)
	hash, err := bundle.GetChecksum()

	if err != nil {
		return destination, err
	}

	destination = bundle.ToHashed(destination, hash)
	destination = bundle.ToExtension(destination, ".css")

	return destination, nil
}

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:  "css",
		Extensions: []string{".css"},
		Init:       Init,
		Related:    Related,
		Execute:    Execute,
		Optimize:   generic.Optimize,
		Delete:     Delete,
		Resolve:    Resolve,
	}
}
