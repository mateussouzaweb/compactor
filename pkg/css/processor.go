package css

import (
	"regexp"
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// Init processor
func Init(options *compactor.Options) error {
	return os.NodeRequire("sass", "sass")
}

// Resolve processor
func Resolve(options *compactor.Options, file *compactor.File) (string, error) {

	hash := file.Checksum
	destination := options.ToDestination(file.Path)
	destination = options.ToHashed(destination, hash)
	destination = options.ToExtension(destination, ".css")

	return destination, nil
}

// Related processor
func Related(options *compactor.Options, file *compactor.File) ([]compactor.Related, error) {

	var related []compactor.Related

	// Add possible source map
	fileMap := strings.TrimSuffix(file.Path, file.Extension)
	fileMap = fileMap + ".css.map"

	related = append(related, compactor.Related{
		Type:       "source-map",
		Dependency: true,
		Source:     "",
		Path:       os.File(fileMap),
		File:       compactor.GetFile(fileMap),
	})

	// Detect imports
	regex := regexp.MustCompile(`@import ?("(.+)"|'(.+)');?`)
	matches := regex.FindAllStringSubmatch(file.Content, -1)
	extensions := []string{".scss", ".sass", ".css"}

	for _, match := range matches {
		source := match[0]
		path := strings.Trim(match[1], `'"`)
		filepath := os.Resolve(path, extensions, os.Dir(file.Path))

		if compactor.GetFile(filepath).Path != "" {
			related = append(related, compactor.Related{
				Type:       "import",
				Dependency: true,
				Source:     source,
				Path:       path,
				File:       compactor.GetFile(filepath),
			})
		}
	}

	return related, nil
}

// Transform processor
func Transform(options *compactor.Options, file *compactor.File) error {

	args := []string{
		file.Path + ":" + file.Destination,
	}

	if options.ShouldCompress(file.Path) {
		args = append(args, "--style", "compressed")
	}

	if options.ShouldGenerateSourceMap(file.Path) {
		args = append(args, "--source-map", "--embed-sources")
	}

	_, err := os.Exec(
		"sass",
		args...,
	)

	if err != nil {
		return err
	}

	return nil
}

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:  "css",
		Extensions: []string{".css"},
		Init:       Init,
		Resolve:    Resolve,
		Related:    Related,
		Transform:  Transform,
		Optimize:   generic.Optimize,
	}
}
