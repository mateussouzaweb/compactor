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
	)

	if err != nil {
		return err
	}

	// Minify
	_, err = compactor.ExecCommand(
		"uglifyjs",
		context.Destination,
		"--output", context.Destination,
		"--compress",
		"--comments",
		"--source-map", "content='"+context.Destination+".map'",
	)

	if err == nil {
		context.Processed = true
	}

	return err
}
