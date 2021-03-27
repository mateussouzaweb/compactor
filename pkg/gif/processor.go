package gif

import (
	"github.com/mateussouzaweb/compactor/compactor"
)

// GIF processor
func Processor(context *compactor.Context, options *compactor.Options) error {

	err := compactor.CopyFile(context.Source, context.Destination)

	if err != nil {
		return err
	}

	_, err = compactor.ExecCommand(
		"gifsicle",
		"-03",
		context.Destination,
		"-o", context.Destination,
	)

	if err == nil {
		context.Processed = true
	}

	return err
}
