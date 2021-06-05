package webp

import (
	"fmt"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// CreateCopy make a WEBP copy of a image file from almost any format
func CreateCopy(source string, destination string, quality int) error {

	_, err := os.Exec(
		"cwebp",
		"-q", fmt.Sprintf("%d", quality),
		destination,
		"-o", destination+".webp",
	)

	return err
}

// WEBP processor
func RunProcessor(bundle *compactor.Bundle) error {

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

func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Extensions: []string{".webp"},
		Run:        RunProcessor,
		Delete:     generic.DeleteProcessor,
		Resolve:    generic.ResolveProcessor,
	}
}
