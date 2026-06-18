package typescript

import (
	"fmt"
	"regexp"
	"strings"

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

	return InitConfig(options.Source.Path)
}

// Shutdown processor
func Shutdown(options *processor.Options) error {
	return _service.Shutdown()
}

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

	// Read config if not loaded yet
	if _tsConfigFile == "" {
		err := InitConfig(file.Root)
		if err != nil {
			return related, err
		}
	}

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

	// Add possible type declaration
	declaration := file.Path + ".d"
	related = append(related, processor.Related{
		Type:       "declaration",
		Dependency: true,
		Source:     "",
		Path:       system.File(declaration),
		File:       processor.GetFile(declaration),
	})

	// Detect imports
	regex := regexp.MustCompile(`import ?((.+) ?from ?)?("(.+)"|'(.+)');?`)
	matches := regex.FindAllStringSubmatch(file.Content, -1)
	extensions := []string{".js", ".mjs", ".jsx", ".ts", ".mts", ".tsx"}

	for _, match := range matches {
		source := match[0]
		path := strings.Trim(match[3], `'"`)
		thePath := FindRealPath(path)
		filePath := system.Resolve(thePath, extensions, system.Dir(file.Path))

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

// FindRealPath transform TSConfig paths into real path values
func FindRealPath(path string) string {

	paths, ok := _tsConfig.CompilerOptions["paths"].(map[string]any)

	if ok {
		for key, values := range paths {
			find := strings.Trim(key, "*")
			value := fmt.Sprintf("%v", values.([]any)[0])
			replace := strings.Trim(value, "*")
			path = strings.Replace(path, find, replace, 1)
		}
	}

	return path
}

// Transform processor
func Transform(options *processor.Options, file *processor.File) error {

	// Copy from user config file
	config := *_tsConfig

	if config.CompilerOptions == nil {
		config.CompilerOptions = make(map[string]any)
	}

	// Make sure output is present
	config.CompilerOptions["emitDeclarationOnly"] = false
	config.CompilerOptions["noEmit"] = false
	config.CompilerOptions["noEmitOnError"] = true

	// To compile correctly we need to force noLib, noResolve and isolatedModules
	config.CompilerOptions["noLib"] = true
	config.CompilerOptions["noResolve"] = true
	config.CompilerOptions["isolatedModules"] = true

	// Enable source maps
	if options.ShouldGenerateSourceMap(file.Path) {
		config.CompilerOptions["sourceMap"] = true
		config.CompilerOptions["inlineSources"] = true
		config.CompilerOptions["sourceRoot"] = ""
	}

	// Run transpilation
	err := _service.Execute(&config, file)
	if err != nil {
		return err
	}

	// Update paths after transpile code with correct final destinations
	for _, related := range file.Related {
		if related.File.Exists && related.Type == "import" {

			relativePath := system.Relative(system.Dir(file.Destination), related.File.Destination)
			if !strings.HasPrefix(relativePath, "../") {
				relativePath = "./" + relativePath
			}

			oldSource := related.Source
			firstIndex := strings.LastIndex(oldSource, related.Path)
			lastIndex := firstIndex + len(related.Path)
			newSource := oldSource[:firstIndex] + relativePath + oldSource[lastIndex:]

			err = system.Replace(
				file.Destination,
				oldSource,
				newSource,
			)

			if err != nil {
				return err
			}

		}
	}

	return nil
}

// Optimize processor
func Optimize(options *processor.Options, file *processor.File) error {

	if !options.ShouldCompress(file.Path) {
		return nil
	}

	args := []string{
		file.Destination,
		"--output", file.Destination,
		"--compress",
		"--comments",
	}

	if options.ShouldGenerateSourceMap(file.Path) {

		sourceOptions := strings.Join([]string{
			"includeSources",
			"filename='" + file.File + ".map'",
			"url='" + file.File + ".map'",
			"content='" + file.Destination + ".map'",
		}, ",")

		args = append(args, "--source-map", sourceOptions)

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
		Namespace:  "typescript",
		Extensions: []string{".js", ".mjs", ".jsx", ".ts", ".mts", ".tsx"},
		Init:       Init,
		Shutdown:   Shutdown,
		Resolve:    Resolve,
		Related:    Related,
		Transform:  Transform,
		Optimize:   Optimize,
	}
}
