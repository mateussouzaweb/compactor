package png

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/pkg/webp"
)

// PNG processor
func Processor(context *compactor.Context, options *compactor.Options) error {

	err := compactor.CopyFile(context.Source, context.Destination)

	if err != nil {
		return err
	}

	if options.Compress {
		_, err = compactor.ExecCommand(
			"optipng",
			"--quiet",
			context.Destination,
		)

		if err != nil {
			return err
		}
	}

	if options.Progressive {
		err = webp.CreateCopy(context.Source, context.Destination, 75)
	}

	if err == nil {
		context.Processed = true
	}

	return err
}
