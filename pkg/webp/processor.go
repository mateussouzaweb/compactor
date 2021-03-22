package webp

import (
	"fmt"

	"github.com/mateussouzaweb/compactor/compactor"
)

// CreateCopy make a WEBP copy of a image file from almost any format
func CreateCopy(source string, destination string, quality int) error {

	if compactor.ExistFile(source + ".webp") {
		return nil
	}

	_, err := compactor.ExecCommand(
		"cwebp",
		"-q", fmt.Sprintf("%d", quality),
		destination,
		"-o", destination+".webp",
	)

	return err
}

// WEBP processor
func Processor(context *compactor.Context) error {

	err := compactor.CopyFile(context.Source, context.Destination)

	if err == nil {
		context.Processed = true
	}

	return err
}
