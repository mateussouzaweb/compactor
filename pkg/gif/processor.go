package gif

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// GIF processor
func Processor(action *compactor.Action, bundle *compactor.Bundle, logger *compactor.Logger) error {

	if action.IsDelete() {
		return generic.DeleteProcessor(bundle, logger, []string{})
	}

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
