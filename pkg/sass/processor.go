package sass

import (
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
)

// Sass processor
func Processor(context *compactor.Context, options *compactor.Options) error {

	context.Destination = strings.Replace(
		context.Destination, ".scss", ".css", 1,
	)
	context.Destination = strings.Replace(
		context.Destination, ".sass", ".css", 1,
	)

	args := []string{
		context.Source + ":" + context.Destination,
	}

	if options.ShouldCompress(context) {
		args = append(args, "--style", "compressed")
	}

	if options.ShouldGenerateSourceMap(context) {
		args = append(args, "--source-map", "--embed-sources")
	}

	_, err := compactor.ExecCommand(
		"sass",
		args...,
	)

	if err == nil {
		context.Processed = true
	}

	return err
}
