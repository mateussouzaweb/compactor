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

// Transform processor
func Transform(options *compactor.Options, file *compactor.File) error {

	err := os.Copy(file.Path, file.Destination)

	if err != nil {
		return err
	}

	return nil
}

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:  "webp",
		Extensions: []string{".webp"},
		Init:       generic.Init,
		Shutdown:   generic.Shutdown,
		Resolve:    generic.Resolve,
		Related:    generic.Related,
		Transform:  Transform,
		Optimize:   generic.Optimize,
	}
}
