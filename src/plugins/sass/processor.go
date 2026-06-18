package sass

import (
	"regexp"
	"strings"

	"github.com/mateussouzaweb/compactor/src/plugins/generic"
	"github.com/mateussouzaweb/compactor/src/processor"
	"github.com/mateussouzaweb/compactor/src/system"
)

var _service TranspilerService

// Init processor
func Init(options *processor.Options) error {

	err := _service.Init()
	if err != nil {
		return err
	}

	return nil
}

// Shutdown processor
func Shutdown(options *processor.Options) error {
	return _service.Shutdown()
}

// Resolve processor
func Resolve(options *processor.Options, file *processor.File) (string, error) {

	destination := options.ToDestination(file.Path)
	destination = options.ToExtension(destination, ".css")

	if options.Destination.Hashed {
		hash := file.Checksum[len(file.Checksum)-1]
		destination = options.ToHashed(destination, hash)
	}

	return destination, nil
}

// Related processor
func Related(options *processor.Options, file *processor.File) ([]processor.Related, error) {

	var related []processor.Related

	// Add possible source map
	fileMap := strings.TrimSuffix(file.Path, file.Extension)
	fileMap = fileMap + ".css.map"

	related = append(related, processor.Related{
		Type:       "source-map",
		Dependency: true,
		Source:     "",
		Path:       system.File(fileMap),
		File:       processor.GetFile(fileMap),
	})

	// Detect imports
	regex := regexp.MustCompile(`@import ?("(.+)"|'(.+)');?`)
	matches := regex.FindAllStringSubmatch(file.Content, -1)
	extensions := []string{".scss", ".sass", ".css"}

	for _, match := range matches {
		source := match[0]
		path := strings.Trim(match[1], `'"`)
		filePath := system.Resolve(path, extensions, system.Dir(file.Path))

		if processor.GetFile(filePath).Path != "" {
			related = append(related, processor.Related{
				Type:       "import",
				Dependency: true,
				Source:     source,
				Path:       path,
				File:       processor.GetFile(filePath),
			})
		}
	}

	return related, nil
}

// Transform processor
func Transform(options *processor.Options, file *processor.File) error {

	// Create config
	config := &SassConfig{
		Style: "expanded",
	}

	if options.ShouldCompress(file.Path) {
		config.Style = "compressed"
	}
	if options.ShouldGenerateSourceMap(file.Path) {
		config.SourceMap = true
		config.SourceMapIncludeSources = true
	}

	// Run transpilation
	err := _service.Execute(config, file)
	if err != nil {
		return err
	}

	return nil
}

// Plugin return the compactor plugin instance
func Plugin() *processor.Plugin {
	return &processor.Plugin{
		Namespace:  "sass",
		Extensions: []string{".sass", ".scss", ".css"},
		Init:       Init,
		Shutdown:   Shutdown,
		Resolve:    Resolve,
		Related:    Related,
		Transform:  Transform,
		Optimize:   generic.Optimize,
	}
}
