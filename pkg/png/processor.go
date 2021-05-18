package png

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/pkg/webp"
)

// PNG processor
func Processor(bundle *compactor.Bundle, logger *compactor.Logger) error {

	files := bundle.GetFiles()

	for _, file := range files {

		destination := bundle.ToDestination(file)
		err := compactor.CopyFile(file, destination)

		if err != nil {
			return err
		}

		if bundle.ShouldCompress(file) {
			_, err = compactor.ExecCommand(
				"optipng",
				"--quiet",
				destination,
			)

			if err != nil {
				return err
			}
		}

		if bundle.ShouldGenerateProgressive(file) {
			err = webp.CreateCopy(file, destination, 75)
		}

		if err != nil {
			return err
		}

		logger.AddProcessed(file)

	}

	return nil
}
