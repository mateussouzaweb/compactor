package javascript

import (
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
)

// Javascript processor
func Processor(context *compactor.Context, options *compactor.Options) error {

	sourceOptions := strings.Join([]string{
		"includeSources",
		"filename='" + context.File + ".map'",
		"url='" + context.File + ".map'",
	}, ",")

	_, err := compactor.ExecCommand(
		"uglifyjs",
		context.Source,
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
