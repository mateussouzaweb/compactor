package webp

import (
	"fmt"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// Init processor
func Init(bundle *compactor.Bundle) error {
	return os.NodeRequire("cwebp", "cwebp-bin")
}

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

// Execute processor
func Execute(bundle *compactor.Bundle) error {

	destination := bundle.ToDestination(bundle.Item.Path)
	err := os.Copy(bundle.Item.Path, destination)

	if err != nil {
		return err
	}

	bundle.Processed(bundle.Item.Path)

	return nil
}

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:  "webp",
		Extensions: []string{".webp"},
		Init:       Init,
		Related:    generic.Related,
		Execute:    Execute,
		Optimize:   generic.Optimize,
		Delete:     generic.Delete,
		Resolve:    generic.Resolve,
	}
}
