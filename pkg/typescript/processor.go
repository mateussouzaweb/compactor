package typescript

import (
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
)

// Typescript processor
func Processor(bundle *compactor.Bundle, logger *compactor.Logger) error {

	files := bundle.GetFiles()
	target, isDir := bundle.GetDestination()
	multiple := []string{}

	for _, file := range files {

		// Compile each file individually
		destination := bundle.ToDestination(file)
		destination = strings.Replace(destination, ".tsx", ".js", 1)
		destination = strings.Replace(destination, ".ts", ".js", 1)

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
		_, err := compactor.ExecCommand(
			"tsc",
			args...,
		)

		if err != nil {
			return err
		}

		if !isDir {
			multiple = append(multiple, destination)
			continue
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

				name := bundle.CleanName(destination)
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

	if isDir {
		return nil
	}

	destination := target
	destination = strings.Replace(destination, ".tsx", ".js", 1)
	destination = strings.Replace(destination, ".ts", ".js", 1)

	args := []string{}
	args = append(args, multiple...)
	args = append(args, "--output", destination)

	// Compress
	if bundle.ShouldCompress(target) {
		args = append(args, "--compress", "--comments")
	}

	// SourceMap
	if bundle.ShouldGenerateSourceMap(target) {

		maps := ""
		for _, file := range multiple {
			maps += "," + file + ".map"
		}
		maps = strings.TrimLeft(maps, ",")

		name := bundle.CleanName(destination)
		sourceOptions := strings.Join([]string{
			"includeSources",
			"filename='" + name + ".map'",
			"url='" + name + ".map'",
			"content='" + maps + "'",
		}, ",")

		args = append(args, "--source-map", sourceOptions)

	}

	_, err := compactor.ExecCommand(
		"uglifyjs",
		args...,
	)

	if err == nil {
		logger.AddProcessed(target)
	}

	return err
}
