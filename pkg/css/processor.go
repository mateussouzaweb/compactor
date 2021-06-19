package css

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// CSS processor
func RunProcessor(bundle *compactor.Bundle) error {

	if bundle.ShouldOutputToMany() {

		for _, item := range bundle.Items {

			if !item.Exists {
				continue
			}

			destination := bundle.ToDestination(item.Path)
			destination = bundle.ToHashed(destination, item.Checksum)
			destination = bundle.ToExtension(destination, ".css")

			args := []string{
				item.Path + ":" + destination,
			}

			if bundle.ShouldCompress(item.Path) {
				args = append(args, "--style", "compressed")
			}

			if bundle.ShouldGenerateSourceMap(item.Path) {
				args = append(args, "--source-map", "--embed-sources")
			}

			_, err := os.Exec(
				"sass",
				args...,
			)

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

	hash, err := os.Checksum(content)

	if err != nil {
		return err
	}

	destination := bundle.ToDestination(bundle.Destination.File)
	destination = bundle.ToHashed(destination, hash)
	destination = bundle.ToExtension(destination, ".css")
	perm := bundle.Items[0].Permission
	err = os.Write(destination, content, perm)

	if err != nil {
		return err
	}

	args := []string{
		destination + ":" + destination,
	}

	if bundle.ShouldCompress(destination) {
		args = append(args, "--style", "compressed")
	}

	if bundle.ShouldGenerateSourceMap(destination) {
		args = append(args, "--source-map", "--embed-sources")
	}

	_, err = os.Exec(
		"sass",
		args...,
	)

	if err == nil {
		bundle.Processed(destination)
	}

	return nil
}

// CSS delete processor
func DeleteProcessor(bundle *compactor.Bundle) error {

	err := generic.DeleteProcessor(bundle)

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

// ResolveProcessor fix the path for given file path
func ResolveProcessor(path string) (string, error) {

	destination, err := generic.ResolveProcessor(path)

	if err != nil {
		return destination, err
	}

	bundle := compactor.GetBundleFor(path)

	if bundle.ShouldOutputToMany() {

		source := bundle.ToSource(path)
		item := compactor.Get(source)

		destination := bundle.ToHashed(path, item.Checksum)
		destination = bundle.ToExtension(destination, ".css")

		return destination, nil
	}

	content := ""
	for _, item := range bundle.Items {
		if item.Exists {
			content += item.Content
		}
	}

	hash, err := os.Checksum(content)

	if err != nil {
		return destination, err
	}

	destination = bundle.ToHashed(destination, hash)
	destination = bundle.ToExtension(destination, ".css")

	return destination, nil
}

func Plugin() *compactor.Plugin {

	os.NodeRequire("sass", "sass")

	return &compactor.Plugin{
		Extensions: []string{".css"},
		Run:        RunProcessor,
		Delete:     DeleteProcessor,
		Resolve:    ResolveProcessor,
	}
}
