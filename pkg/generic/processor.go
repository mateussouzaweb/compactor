package generic

import (
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
)

// Init processor
func Init(bundle *compactor.Bundle) error {
	return nil
}

// Related detect the dependencies of the file
func Related(item *compactor.Item) ([]compactor.Related, error) {
	var found []compactor.Related
	return found, nil
}

// Resolve return the clean bundle destination path for given source file path
func Resolve(path string, item *compactor.Item) (string, error) {

	bundle := compactor.GetBundle(path)
	destination := bundle.ToDestination(bundle.Item.Path)
	destination = bundle.CleanPath(destination)

	if strings.HasPrefix(path, "/") {
		destination = "/" + destination
	}

	return destination, nil
}

// Execute create generic copy of file(s) content to destination
func Execute(bundle *compactor.Bundle) error {

	content := bundle.Item.Content
	perm := bundle.Item.Permission

	destination := bundle.ToDestination(bundle.Item.Path)
	err := os.Write(destination, content, perm)

	if err != nil {
		return err
	}

	return nil
}

// Optimize apply optimizations into the destination file
func Optimize(bundle *compactor.Bundle) error {
	return nil
}

// Delete remove the destination file(s)
func Delete(bundle *compactor.Bundle) error {

	toDelete := []string{}

	// Item file name
	destination := bundle.ToDestination(bundle.Item.Path)
	checksum := bundle.ToHashed(destination, bundle.Item.Checksum)
	previous := bundle.ToHashed(destination, bundle.Item.Previous)
	toDelete = append(toDelete, destination, checksum, previous)

	// Related dependencies
	for _, related := range bundle.Item.Related {
		if related.Dependency {

			// Variation from related item path
			destination := bundle.ToDestination(related.Item.Path)
			toDelete = append(toDelete, destination)

			// Variations from related item checksum
			hashed := bundle.ToHashed(destination, related.Item.Checksum)
			previous := bundle.ToHashed(destination, related.Item.Previous)
			toDelete = append(toDelete, destination, hashed, previous)

			// Variations from bundle item checksum
			hashedFromItem := bundle.ToHashed(destination, bundle.Item.Checksum)
			previousFromItem := bundle.ToHashed(destination, bundle.Item.Previous)
			toDelete = append(toDelete, hashedFromItem, previousFromItem)

		}
	}

	for _, file := range toDelete {

		if !os.Exist(file) {
			continue
		}

		err := os.Delete(file)

		if err != nil {
			return err
		}

	}

	return nil
}

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:  "generic",
		Extensions: []string{},
		Init:       Init,
		Related:    Related,
		Resolve:    Resolve,
		Execute:    Execute,
		Optimize:   Optimize,
		Delete:     Delete,
	}
}
