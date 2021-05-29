package css

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// CSS processor
func Processor(action *compactor.Action, bundle *compactor.Bundle, logger *compactor.Logger) error {

	if action.IsDelete() {
		return generic.DeleteProcessor(bundle, logger, []string{"css.map"})
	}

	files := bundle.GetFiles()

	if bundle.IsToMultipleDestinations() {

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

	content, perm, err := compactor.ReadFilesAndPermission(files)

	if err != nil {
		return err
	}

	destination := bundle.GetDestination()
	err = compactor.WriteFile(destination, content, perm)

	if err != nil {
		return err
	}

	args := []string{
		destination + ":" + destination,
	}

	if bundle.ShouldCompress(destination) {
		args = append(args, "--style", "compressed")
	}

	if bundle.ShouldGenerateSourceMap(destination) {
		args = append(args, "--source-map", "--embed-sources")
	}

	_, err = compactor.ExecCommand(
		"sass",
		args...,
	)

	if err == nil {
		logger.AddProcessed(destination)
	}

	return nil
}
