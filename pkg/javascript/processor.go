package javascript

import (
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// Javascript processor
func Processor(action *compactor.Action, bundle *compactor.Bundle, logger *compactor.Logger) error {

	if action.IsDelete() {
		return generic.DeleteProcessor(bundle, logger, []string{"js.map"})
	}

	files := bundle.GetFiles()

	if bundle.IsToMultipleDestinations() {

		for _, file := range files {

			destination := bundle.ToDestination(file)

			args := []string{}
			args = append(args, file)
			args = append(args, "--output", destination)

			if bundle.ShouldCompress(file) {
				args = append(args, "--compress", "--comments")
			} else {
				args = append(args, "--beautify")
			}

			if bundle.ShouldGenerateSourceMap(file) {
				name := bundle.CleanName(destination)
				args = append(args, "--source-map", strings.Join([]string{
					"includeSources",
					"filename='" + name + ".map'",
					"url='" + name + ".map'",
				}, ","))
			}

			_, err := compactor.ExecCommand(
				"uglifyjs",
				args...,
			)

			if err != nil {
				return err
			}

			logger.AddProcessed(file)

		}

		return nil
	}

	destination := bundle.GetDestination()

	args := []string{}
	args = append(args, files...)
	args = append(args, "--output", destination)

	if bundle.ShouldCompress(destination) {
		args = append(args, "--compress", "--comments")
	} else {
		args = append(args, "--beautify")
	}

	if bundle.ShouldGenerateSourceMap(destination) {
		name := bundle.CleanName(destination)
		args = append(args, "--source-map", strings.Join([]string{
			"includeSources",
			"filename='" + name + ".map'",
			"url='" + name + ".map'",
		}, ","))
	}

	_, err := compactor.ExecCommand(
		"uglifyjs",
		args...,
	)

	if err == nil {
		logger.AddWritten(destination)
	}

	return err
}
