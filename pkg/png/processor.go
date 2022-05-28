package png

import (
	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
	"github.com/mateussouzaweb/compactor/pkg/webp"
)

// Init processor
func Init(options *compactor.Options) error {

	err := os.NodeRequire("optipng", "optipng-bin")

	if err != nil {
		return err
	}

	return os.NodeRequire("cwebp", "cwebp-bin")
}

// Related processor
func Related(options *compactor.Options, file *compactor.File) ([]compactor.Related, error) {

	var related []compactor.Related

	// Add possible progressive image
	filePath := file.Path + ".webp"
	related = append(related, compactor.Related{
		Type:       "alternative",
		Dependency: true,
		Source:     "",
		Path:       os.File(filePath),
		File:       compactor.GetFile(filePath),
	})

	return related, nil
}

// Transform processor
func Transform(options *compactor.Options, file *compactor.File) error {

	err := os.Copy(file.Path, file.Destination)

	if err != nil {
		return err
	}

	return nil
}

// Optimize processor
func Optimize(options *compactor.Options, file *compactor.File) error {

	if options.ShouldCompress(file.Path) {
		_, err := os.Exec(
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
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:  "png",
		Extensions: []string{".png"},
		Init:       Init,
		Resolve:    generic.Resolve,
		Related:    Related,
		Transform:  Transform,
		Optimize:   Optimize,
	}
}
