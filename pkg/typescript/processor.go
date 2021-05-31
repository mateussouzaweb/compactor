package typescript

import (
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// Typescript processor
func Processor(action *compactor.Action, bundle *compactor.Bundle, logger *compactor.Logger) error {

	if action.IsDelete() {
		return generic.DeleteProcessor(bundle, logger, []string{"js.map"})
	}

	files := bundle.GetFiles()

	for _, file := range files {

		hash, err := compactor.GetChecksum([]string{file})

		if err != nil {
			return err
		}

		destination := bundle.ToDestination(file)
		destination = bundle.ToHashed(destination, hash)
		destination = bundle.ToExtension(destination, "js")

		args := []string{
			file,
			"--outFile", destination,
			"--target", "ES6",
			"--removeComments",
		}

		if bundle.ShouldGenerateSourceMap(file) {
			args = append(args, "--sourceMap", "--inlineSources")
		}

		// Compile
		_, err = compactor.ExecCommand(
			"tsc",
			args...,
		)

		if err != nil {
			return err
		}

		// Compress
		if bundle.ShouldCompress(file) {

			args = []string{
				destination,
				"--output", destination,
				"--compress",
				"--comments",
			}

			if bundle.ShouldGenerateSourceMap(file) {

				name := compactor.CleanName(destination)
				sourceOptions := strings.Join([]string{
					"includeSources",
					"filename='" + name + ".map'",
					"url='" + name + ".map'",
					"content='" + destination + ".map'",
				}, ",")

				args = append(args, "--source-map", sourceOptions)

			}

			_, err = compactor.ExecCommand(
				"uglifyjs",
				args...,
			)

			if err != nil {
				return err
			}

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
		destination = bundle.ToExtension(destination, "js")
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
	destination = bundle.ToExtension(destination, "js")
	destination = bundle.CleanPath(destination)

	if src[0] == '/' {
		destination = "/" + destination
	}

	return destination, nil
}
