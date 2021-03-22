package generic

import "github.com/mateussouzaweb/compactor/compactor"

func Processor(context *compactor.Context) error {

	err := compactor.CopyFile(context.Source, context.Destination)

	if err == nil {
		context.Processed = true
	}

	return err
}
