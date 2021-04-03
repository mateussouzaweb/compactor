package jpeg

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/pkg/webp"
)

// JPEG processor
func Processor(context *compactor.Context, options *compactor.Options) error {

	err := compactor.CopyFile(context.Source, context.Destination)

	if err != nil {
		return err
	}

	if options.ShouldCompress(context) {
		_, err = compactor.ExecCommand(
			"jpegoptim",
			"--quiet",
			"--strip-all",
			"--all-progressive",
			"--overwrite",
			context.Destination,
		)

		if err != nil {
			return err
		}
	}

	if options.GenerateProgressive(context) {
		err = webp.CreateCopy(context.Source, context.Destination, 75)
	}

	if err == nil {
		context.Processed = true
	}

	return err
}
