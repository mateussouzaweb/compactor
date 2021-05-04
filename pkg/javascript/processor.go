package javascript

import (
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
)

// Javascript processor
func Processor(context *compactor.Context, options *compactor.Options) error {

	args := []string{
		context.Source,
		"--output", context.Destination,
	}

	if options.ShouldCompress(context) {
		args = append(args, "--compress", "--comments")
	}

	if options.ShouldGenerateSourceMap(context) {
		args = append(args, "--source-map", strings.Join([]string{
			"includeSources",
			"filename='" + context.File + ".map'",
			"url='" + context.File + ".map'",
		}, ","))
	}

	_, err := compactor.ExecCommand(
		"uglifyjs",
		args...,
	)

	if err == nil {
		context.Processed = true
	}

	return err
}
