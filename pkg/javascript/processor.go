package javascript

import (
	"regexp"
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// Init processor
func Init(options *compactor.Options) error {
	return os.NodeRequire("terser", "terser")
}

// Resolve processor
func Resolve(options *compactor.Options, file *compactor.File) (string, error) {

	hash := file.Checksum
	destination := options.ToDestination(file.Path)
	destination = options.ToHashed(destination, hash)
	destination = options.ToExtension(destination, ".js")

	return destination, nil
}

// Related processor
func Related(options *compactor.Options, file *compactor.File) ([]compactor.Related, error) {

	var related []compactor.Related

	// Add possible source map
	filemap := strings.TrimSuffix(file.Path, file.Extension)
	filemap = filemap + ".js.map"

	related = append(related, compactor.Related{
		Type:       "source-map",
		Dependency: true,
		Source:     "",
		Path:       os.File(filemap),
		File:       compactor.GetFile(filemap),
	})

	// Detect imports
	regex := regexp.MustCompile(`import ?((.+) ?from ?)?("(.+)"|'(.+)');?`)
	matches := regex.FindAllStringSubmatch(file.Content, -1)
	extensions := []string{".js", ".mjs"}

	for _, match := range matches {
		source := match[0]
		path := strings.Trim(match[3], `'"`)
		filepath := os.Resolve(path, extensions, os.Dir(file.Path))

		if compactor.GetFile(filepath).Path != "" {
			related = append(related, compactor.Related{
				Type:       "import",
				Dependency: false,
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

	_, err := os.Exec(
		"terser",
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
		Namespace:  "javascript",
		Extensions: []string{".js", ".mjs"},
		Init:       Init,
		Resolve:    Resolve,
		Related:    Related,
		Transform:  Transform,
		Optimize:   generic.Optimize,
	}
}
