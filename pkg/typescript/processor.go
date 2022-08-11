package typescript

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
)

var _tsConfig *TSConfig
var _tsConfigFile string

// Init processor
func Init(options *compactor.Options) error {
	return InitConfig(options.Source.Path)
}

// InitConfig find and read tsconfig file from given path
func InitConfig(path string) error {

	var err error
	_tsConfigFile = FindConfig(path)
	_tsConfig, err = ReadConfig(_tsConfigFile)

	return err
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

	related = append(related, compactor.Related{
		Type:       "source-map",
		Dependency: true,
		Source:     "",
		Path:       os.File(fileMap),
		File:       compactor.GetFile(fileMap),
	})

	// Add possible type declaration
	declaration := file.Path + ".d"
	related = append(related, compactor.Related{
		Type:       "declaration",
		Dependency: true,
		Source:     "",
		Path:       os.File(declaration),
		File:       compactor.GetFile(declaration),
	})

	// Detect imports
	regex := regexp.MustCompile(`import ?((.+) ?from ?)?("(.+)"|'(.+)');?`)
	matches := regex.FindAllStringSubmatch(file.Content, -1)
	extensions := []string{".js", ".mjs", ".jsx", ".ts", ".mts", ".tsx"}

	for _, match := range matches {
		source := match[0]
		path := strings.Trim(match[3], `'"`)
		thePath := FindRealPath(path)
		filePath := os.Resolve(thePath, extensions, os.Dir(file.Path))

		if compactor.GetFile(filePath).Path != "" {
			related = append(related, compactor.Related{
				Type:       "import",
				Dependency: false,
				Source:     source,
				Path:       path,
				File:       compactor.GetFile(filePath),
			})
		}
	}

	return related, nil
}

// FindRealPath transform TSConfig paths into real path values
func FindRealPath(path string) string {

	paths, ok := _tsConfig.CompilerOptions["paths"].(map[string]interface{})

	if ok {
		for key, values := range paths {
			find := strings.Trim(key, "*")
			value := fmt.Sprintf("%v", values.([]interface{})[0])
			replace := strings.Trim(value, "*")
			path = strings.Replace(path, find, replace, 1)
		}
	}

	return path
}

// Transform processor
func Transform(options *compactor.Options, file *compactor.File) error {

	// Copy from user config file
	config := *_tsConfig

	if config.CompilerOptions == nil {
		config.CompilerOptions = make(map[string]interface{})
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
	err := RunTranspiler(Transpiler{
		Config:      &config,
		File:        file.Path,
		Content:     file.Content,
		Destination: file.Destination,
	})

	if err != nil {
		return err
	}

	// Update paths after transpile code with correct final destinations
	for _, related := range file.Related {
		if related.File.Exists && related.Type == "import" {

			relativePath := os.Relative(os.Dir(file.Destination), related.File.Destination)
			if !strings.HasPrefix(relativePath, "../") {
				relativePath = "./" + relativePath
			}

			oldSource := related.Source
			firstIndex := strings.LastIndex(oldSource, related.Path)
			lastIndex := firstIndex + len(related.Path)
			newSource := oldSource[:firstIndex] + relativePath + oldSource[lastIndex:]

			err = os.Replace(
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
func Optimize(options *compactor.Options, file *compactor.File) error {

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

	_, err := os.Exec(
		"terser",
		args...,
	)

	return err
}

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:  "typescript",
		Extensions: []string{".js", ".mjs", ".jsx", ".ts", ".mts", ".tsx"},
		Init:       Init,
		Resolve:    Resolve,
		Related:    Related,
		Transform:  Transform,
		Optimize:   Optimize,
	}
}
