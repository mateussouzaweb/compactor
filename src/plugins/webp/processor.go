package webp

import (
	"fmt"

	"github.com/mateussouzaweb/compactor/src/plugins/generic"
	"github.com/mateussouzaweb/compactor/src/processor"
	"github.com/mateussouzaweb/compactor/src/system"
)

// CreateCopy make a WEBP copy of a image file from almost any format
func CreateCopy(source string, destination string, quality int) error {

	_, err := system.Exec(
		"cwebp",
		"-q", fmt.Sprintf("%d", quality),
		destination,
		"-o", destination+".webp",
	)

	return err
}

// Transform processor
func Transform(options *processor.Options, file *processor.File) error {

	err := system.Copy(file.Path, file.Destination)
	if err != nil {
		return err
	}

	return nil
}

// Plugin return the compactor plugin instance
func Plugin() *processor.Plugin {
	return &processor.Plugin{
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
