package gif

import (
	"github.com/mateussouzaweb/compactor/compactor"
)

// GIF processor
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
				"gifsicle",
				"-03",
				destination,
				"-o", destination,
			)
		}

		if err != nil {
			return err
		}

		logger.AddProcessed(file)

	}

	return nil
}
