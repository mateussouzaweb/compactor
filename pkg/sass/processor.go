package sass

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// Sass processor
func Processor(action *compactor.Action, bundle *compactor.Bundle, logger *compactor.Logger) error {

	if action.IsDelete() {
		return generic.DeleteProcessor(bundle, logger, []string{".map"})
	}

	files := bundle.GetFiles()
	target := bundle.GetDestination()
	multiple := []string{}

	for _, file := range files {

		if !bundle.IsToMultipleDestinations() {
			multiple = append(multiple, file)
			continue
		}

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

	if bundle.IsToMultipleDestinations() {
		return nil
	}

	// TODO: that is wrong
	content, err := compactor.ReadFiles(multiple)

	if err != nil {
		return err
	}

	perm, err := compactor.GetPermission(multiple[0])

	if err != nil {
		return err
	}

	err = compactor.WriteFile(target, content, perm)

	if err != nil {
		return err
	}

	destination := target
	destination = bundle.ToExtension(destination, "js")

	args := []string{
		target + ":" + destination,
	}

	if bundle.ShouldCompress(target) {
		args = append(args, "--style", "compressed")
	}

	if bundle.ShouldGenerateSourceMap(target) {
		args = append(args, "--source-map", "--embed-sources")
	}

	_, err = compactor.ExecCommand(
		"sass",
		args...,
	)

	if err == nil {
		logger.AddProcessed(target)
	}

	return nil
}
