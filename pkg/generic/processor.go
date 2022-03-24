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

// Dependencies detect the dependencies of the file
func Dependencies(item *compactor.Item) ([]string, error) {
	return []string{}, nil
}

// Execute create generic copy of file(s) content to destination
func Execute(bundle *compactor.Bundle) error {

	content := bundle.GetContent()
	perm := bundle.GetPermission()

	destination := bundle.ToDestination(bundle.Destination.File)
	err := os.Write(destination, content, perm)

	if err == nil {
		bundle.Written(destination)
	}

	return err
}

// Optimize apply optimizations into the destination file
func Optimize(bundle *compactor.Bundle) error {
	return nil
}

// Delete remove the destination file(s)
func Delete(bundle *compactor.Bundle) error {

	toDelete := []string{}

	for _, item := range bundle.Items {

		destination := bundle.ToDestination(item.Path)
		hashed := bundle.ToHashed(destination, item.Checksum)
		previous := bundle.ToHashed(destination, item.Previous)
		toDelete = append(toDelete, destination, hashed, previous)

	}

	content := bundle.GetContent()
	hash, err := os.Checksum(content)

	if err != nil {
		return err
	}

	destination := bundle.ToDestination(bundle.Destination.File)
	hashed := bundle.ToHashed(destination, hash)
	toDelete = append(toDelete, destination, hashed)

	for _, file := range toDelete {

		if !os.Exist(file) {
			continue
		}

		err := os.Delete(file)
		if err != nil {
			return err
		}

		bundle.Deleted(file)

	}

	return nil
}

// Resolve return the clean bundle destination path for given source file path
func Resolve(path string) (string, error) {

	bundle := compactor.GetBundle(path)
	destination := bundle.ToDestination(bundle.Destination.File)
	destination = bundle.CleanPath(destination)

	if strings.HasPrefix(path, "/") {
		destination = "/" + destination
	}

	return destination, nil
}

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:    "generic",
		Extensions:   []string{},
		Init:         Init,
		Dependencies: Dependencies,
		Execute:      Execute,
		Optimize:     Optimize,
		Delete:       Delete,
		Resolve:      Resolve,
	}
}
