package javascript

import (
	"regexp"
	"strings"

	"github.com/mateussouzaweb/compactor/src/plugins/generic"
	"github.com/mateussouzaweb/compactor/src/processor"
	"github.com/mateussouzaweb/compactor/src/system"
)

// Resolve processor
func Resolve(options *processor.Options, file *processor.File) (string, error) {

	destination := options.ToDestination(file.Path)
	destination = options.ToExtension(destination, ".js")

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
	fileMap = fileMap + ".js.map"

	related = append(related, processor.Related{
		Type:       "source-map",
		Dependency: true,
		Source:     "",
		Path:       system.File(fileMap),
		File:       processor.GetFile(fileMap),
	})

	// Detect imports
	regex := regexp.MustCompile(`import ?((.+) ?from ?)?("(.+)"|'(.+)');?`)
	matches := regex.FindAllStringSubmatch(file.Content, -1)
	extensions := []string{".js", ".mjs"}

	for _, match := range matches {
		source := match[0]
		path := strings.Trim(match[3], `'"`)
		filePath := system.Resolve(path, extensions, system.Dir(file.Path))

		if processor.GetFile(filePath).Path != "" {
			related = append(related, processor.Related{
				Type:       "import",
				Dependency: false,
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

	files := []string{file.Path}

	for _, related := range file.Related {
		if related.File.Exists && related.Type == "import" {
			files = append(files, related.File.Path)
		}
	}

	args := []string{}
	args = append(args, files...)
	args = append(args, "--output", file.Destination)

	if options.ShouldCompress(file.Path) {
		args = append(args, "--compress", "--comments")
	} else {
		args = append(args, "--beautify")
	}

	if options.ShouldGenerateSourceMap(file.Path) {
		args = append(args, "--source-map", strings.Join([]string{
			"includeSources",
			"base='" + options.Destination.Path + "'",
			"filename='" + file.File + ".map'",
			"url='" + file.File + ".map'",
		}, ","))
	}

	_, err := system.Exec("terser", args...)
	if err != nil {
		return err
	}

	return nil
}

// Plugin return the compactor plugin instance
func Plugin() *processor.Plugin {
	return &processor.Plugin{
		Namespace:  "javascript",
		Extensions: []string{".js", ".mjs"},
		Init:       generic.Init,
		Shutdown:   generic.Shutdown,
		Resolve:    Resolve,
		Related:    Related,
		Transform:  Transform,
		Optimize:   generic.Optimize,
	}
}
