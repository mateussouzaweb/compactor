package svg

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// SVG minify
func Minify(content string) (string, error) {

	// TODO: Viewbox removal causing bugs
	// _, err = compactor.ExecCommand(
	// 	"svgo",
	// 	"--quiet",
	// 	"--input", target,
	// 	"--output", target,
	// )

	return content, nil
}

// SVG processor
func Processor(action *compactor.Action, bundle *compactor.Bundle, logger *compactor.Logger) error {

	if action.IsDelete() {
		return generic.DeleteProcessor(bundle, logger, []string{})
	}

	files := bundle.GetFiles()

	if bundle.IsToMultipleDestinations() {

		for _, file := range files {

			content, perm, err := compactor.ReadFileAndPermission(file)

			if err != nil {
				return err
			}

			if bundle.ShouldCompress(file) {
				content, err = Minify(content)
				if err != nil {
					return err
				}
			}

			destination := bundle.ToDestination(file)
			err = compactor.WriteFile(destination, content, perm)

			if err != nil {
				return err
			}

			logger.AddProcessed(file)

		}

		return nil
	}

	destination := bundle.GetDestination()
	content, perm, err := compactor.ReadFilesAndPermission(files)

	if err != nil {
		return err
	}

	if bundle.ShouldCompress(destination) {
		content, err = Minify(content)
		if err != nil {
			return err
		}
	}

	err = compactor.WriteFile(destination, content, perm)

	if err == nil {
		logger.AddWritten(destination)
	}

	return err
}
