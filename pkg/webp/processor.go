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
func Processor(bundle *compactor.Bundle, logger *compactor.Logger) error {

	files := bundle.GetFiles()

	for _, file := range files {

		destination := bundle.ToDestination(file)
		err := compactor.CopyFile(file, destination)

		if err != nil {
			return err
		}

		logger.AddProcessed(file)

	}

	return nil
}
