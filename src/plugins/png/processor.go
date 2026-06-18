package png

import (
	"github.com/mateussouzaweb/compactor/src/plugins/generic"
	"github.com/mateussouzaweb/compactor/src/plugins/webp"
	"github.com/mateussouzaweb/compactor/src/processor"
	"github.com/mateussouzaweb/compactor/src/system"
)

// Related processor
func Related(options *processor.Options, file *processor.File) ([]processor.Related, error) {

	var related []processor.Related

	// Add possible progressive image
	filePath := file.Path + ".webp"
	related = append(related, processor.Related{
		Type:       "alternative",
		Dependency: true,
		Source:     "",
		Path:       system.File(filePath),
		File:       processor.GetFile(filePath),
	})

	return related, nil
}

// Transform processor
func Transform(options *processor.Options, file *processor.File) error {

	err := system.Copy(file.Path, file.Destination)
	if err != nil {
		return err
	}

	return nil
}

// Optimize processor
func Optimize(options *processor.Options, file *processor.File) error {

	if options.ShouldCompress(file.Path) {
		_, err := system.Exec(
			"optipng",
			"--quiet",
			file.Destination,
		)

		if err != nil {
			return err
		}
	}

	if options.ShouldGenerateProgressive(file.Path) {
		err := webp.CreateCopy(file.Path, file.Destination, 75)
		if err != nil {
			return err
		}
	}

	return nil
}

// Plugin return the compactor plugin instance
func Plugin() *processor.Plugin {
	return &processor.Plugin{
		Namespace:  "png",
		Extensions: []string{".png"},
		Init:       generic.Init,
		Shutdown:   generic.Shutdown,
		Resolve:    generic.Resolve,
		Related:    Related,
		Transform:  Transform,
		Optimize:   Optimize,
	}
}
