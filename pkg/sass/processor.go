package sass

import (
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
)

// Sass processor.
// Will also compress CSS, so no need to minify
func Processor(context *compactor.Context, options *compactor.Options) error {

	context.Destination = strings.Replace(
		context.Destination, ".scss", ".css", 1,
	)
	context.Destination = strings.Replace(
		context.Destination, ".sass", ".css", 1,
	)

	_, err := compactor.ExecCommand(
		"sass",
		context.Source+":"+context.Destination,
		"--style", "compressed",
		"--source-map", "--embed-sources",
	)

	if err == nil {
		context.Processed = true
	}

	return err
}
