package typescript

import (
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
)

// Typescript processor
func Processor(context *compactor.Context, options *compactor.Options) error {

	context.Destination = strings.Replace(
		context.Destination, ".tsx", ".js", 1,
	)
	context.Destination = strings.Replace(
		context.Destination, ".ts", ".js", 1,
	)

	args := []string{
		context.Source,
		"--outFile", context.Destination,
		"--target", "ES6",
		"--removeComments",
	}

	if options.GenerateSourceMap(context) {
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

	// Compress
	if options.ShouldCompress(context) {

		args = []string{
			context.Destination,
			"--output", context.Destination,
			"--compress",
			"--comments",
		}

		if options.GenerateSourceMap(context) {

			file := strings.Replace(
				context.File, ".ts", ".js", 1,
			)

			sourceOptions := strings.Join([]string{
				"includeSources",
				"filename='" + file + ".map'",
				"url='" + file + ".map'",
				"content='" + context.Destination + ".map'",
			}, ",")

			args = append(args, "--source-map", sourceOptions)

		}

		_, err = compactor.ExecCommand(
			"uglifyjs",
			args...,
		)

	}

	if err == nil {
		context.Processed = true
	}

	return err
}
