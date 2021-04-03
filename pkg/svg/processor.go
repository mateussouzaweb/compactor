package svg

import (
	"github.com/mateussouzaweb/compactor/compactor"
)

// Svg processor
func Processor(context *compactor.Context, options *compactor.Options) error {

	var err error

	if options.ShouldCompress(context) {
		_, err = compactor.ExecCommand(
			"svgo",
			"--quiet",
			"--input", context.Source,
			"--output", context.Destination,
		)
	} else {
		err = compactor.CopyFile(context.Source, context.Destination)
	}

	if err == nil {
		context.Processed = true
	}

	return err
}
