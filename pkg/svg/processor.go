package svg

import (
	"github.com/mateussouzaweb/compactor/compactor"
)

// Svg processor
func Processor(context *compactor.Context) error {

	_, err := compactor.ExecCommand(
		"svgo",
		"--quiet",
		"--input", context.Source,
		"--output", context.Destination,
	)

	if err == nil {
		context.Processed = true
	}

	return err
}
