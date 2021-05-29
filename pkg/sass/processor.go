package sass

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// Sass processor
func Processor(action *compactor.Action, bundle *compactor.Bundle, logger *compactor.Logger) error {

	if action.IsDelete() {
		return generic.DeleteProcessor(bundle, logger, []string{"css.map"})
	}

	files := bundle.GetFiles()

	for _, file := range files {

		destination := bundle.ToDestination(file)
		destination = bundle.ToExtension(destination, "css")

		args := []string{
			file + ":" + destination,
		}

		if bundle.ShouldCompress(file) {
			args = append(args, "--style", "compressed")
		}

		if bundle.ShouldGenerateSourceMap(file) {
			args = append(args, "--source-map", "--embed-sources")
		}

		_, err := compactor.ExecCommand(
			"sass",
			args...,
		)

		if err != nil {
			return err
		}

		logger.AddProcessed(file)

	}

	return nil
}
