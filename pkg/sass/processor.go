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

		hash, err := compactor.GetChecksum([]string{file})

		if err != nil {
			return err
		}

		destination := bundle.ToDestination(file)
		destination = bundle.ToHashed(destination, hash)
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

		_, err = compactor.ExecCommand(
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


// CorrectPath fix the path for given src
func CorrectPath(src string) (string, error) {

	bundle := compactor.RetrieveBundleFor(src)

	if bundle.IsToMultipleDestinations() {

		source := bundle.ToSource(src)
		hash, err := compactor.GetChecksum([]string{source})

		if err != nil {
			return "", err
		}

		destination := bundle.ToDestination(src)
		destination = bundle.ToHashed(destination, hash)
		destination = bundle.ToExtension(destination, "css")
		destination = bundle.CleanPath(destination)

		if src[0] == '/' {
			destination = "/" + destination
		}

		return destination, nil
	}

	files := bundle.GetFiles()
	hash, err := compactor.GetChecksum(files)

	if err != nil {
		return "", err
	}

	destination := bundle.GetDestination()
	destination = bundle.ToHashed(destination, hash)
	destination = bundle.ToExtension(destination, "css")
	destination = bundle.CleanPath(destination)

	if src[0] == '/' {
		destination = "/" + destination
	}

	return destination, nil
}