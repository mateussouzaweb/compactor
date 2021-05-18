package javascript

import (
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
)

// Javascript processor
func Processor(bundle *compactor.Bundle, logger *compactor.Logger) error {

	files := bundle.GetFiles()
	target, isDir := bundle.GetDestination()
	multiple := []string{}

	for _, file := range files {

		if !isDir {
			multiple = append(multiple, file)
			continue
		}

		destination := bundle.ToDestination(file)

		args := []string{}
		args = append(args, file)
		args = append(args, "--output", destination)

		if bundle.ShouldCompress(file) {
			args = append(args, "--compress", "--comments")
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

	if isDir {
		return nil
	}

	args := []string{}
	args = append(args, multiple...)
	args = append(args, "--output", target)

	if bundle.ShouldCompress(target) {
		args = append(args, "--compress", "--comments")
	}

	if bundle.ShouldGenerateSourceMap(target) {
		name := bundle.CleanName(target)
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
		logger.AddWritten(target)
	}

	return err
}
