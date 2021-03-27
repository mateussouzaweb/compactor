package typescript

import (
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
)

// Typescript processor
func Processor(context *compactor.Context) error {

	context.Destination = strings.Replace(
		context.Destination, ".tsx", ".js", 1,
	)
	context.Destination = strings.Replace(
		context.Destination, ".ts", ".js", 1,
	)

	// Compile
	_, err := compactor.ExecCommand(
		"tsc",
		context.Source,
		"--outFile", context.Destination,
		"--target", "ES6",
		"--removeComments",
		"--sourceMap",
		"--inlineSources",
	)

	if err != nil {
		return err
	}

	// Minify
	file := strings.Replace(
		context.File, ".ts", ".js", 1,
	)
	sourceOptions := strings.Join([]string{
		"includeSources",
		"filename='" + file + ".map'",
		"url='" + file + ".map'",
		"content='" + context.Destination + ".map'",
	}, ",")

	_, err = compactor.ExecCommand(
		"uglifyjs",
		context.Destination,
		"--output", context.Destination,
		"--compress",
		"--comments",
		"--source-map", sourceOptions,
	)

	if err == nil {
		context.Processed = true
	}

	return err
}
