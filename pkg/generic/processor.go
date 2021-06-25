package generic

import (
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
)

// Init processor
func InitProcessor(bundle *compactor.Bundle) error {
	return nil
}

// RunProcessor create generic copy of file(s) to destination
func RunProcessor(bundle *compactor.Bundle) error {

	if bundle.ShouldOutputToMany() {

		for _, item := range bundle.Items {

			if !item.Exists {
				continue
			}

			destination := bundle.ToDestination(item.Path)
			err := os.Copy(item.Path, destination)

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
	perm := bundle.Items[0].Permission
	err := os.Write(destination, content, perm)

	if err == nil {
		bundle.Written(destination)
	}

	return err
}

// DeleteProcessor remove the destination file(s)
func DeleteProcessor(bundle *compactor.Bundle) error {

	toDelete := []string{}

	if bundle.ShouldOutputToMany() {
		for _, item := range bundle.Items {

			destination := bundle.ToDestination(item.Path)
			hashed := bundle.ToHashed(destination, item.Checksum)
			previous := bundle.ToHashed(destination, item.Previous)
			toDelete = append(toDelete, destination, hashed, previous)

		}
	} else {

		content := ""
		for _, item := range bundle.Items {
			if item.Exists {
				content += item.Content
			}
		}

		hash, err := os.Checksum(content)

		if err != nil {
			return err
		}

		destination := bundle.ToDestination(bundle.Destination.File)
		hashed := bundle.ToHashed(destination, hash)
		toDelete = append(toDelete, destination, hashed)

	}

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

// ResolveProcessor fix the path for given file path
func ResolveProcessor(path string) (string, error) {

	bundle := compactor.GetBundleFor(path)

	if bundle.ShouldOutputToMany() {

		destination := bundle.ToDestination(path)
		destination = bundle.CleanPath(destination)

		if strings.HasPrefix(path, "/") {
			destination = "/" + destination
		}

		return destination, nil
	}

	destination := bundle.ToDestination(bundle.Destination.File)
	destination = bundle.CleanPath(destination)

	if strings.HasPrefix(path, "/") {
		destination = "/" + destination
	}

	return destination, nil
}

func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Extensions: []string{},
		Init:       InitProcessor,
		Run:        RunProcessor,
		Delete:     DeleteProcessor,
		Resolve:    ResolveProcessor,
	}
}
