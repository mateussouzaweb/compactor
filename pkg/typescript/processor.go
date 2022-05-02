package typescript

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

var _tsConfig *TSConfig
var _tsConfigFile string

// Init processor
func Init(bundle *compactor.Bundle) error {

	err := os.NodeRequire("tsc", "typescript")

	if err != nil {
		return err
	}

	err = os.NodeRequire("terser", "terser")

	if err != nil {
		return err
	}

	return InitConfig(bundle.Source.Path)
}

// InitConfig find and read tsconfig file from given path
func InitConfig(path string) error {

	var err error
	_tsConfigFile = FindConfig(path)
	_tsConfig, err = ReadConfig(_tsConfigFile)

	return err
}

// Related processor
func Related(item *compactor.Item) ([]compactor.Related, error) {

	var related []compactor.Related

	// Read config if not loaded yet
	if _tsConfigFile == "" {
		err := InitConfig(item.Root)

		if err != nil {
			return related, err
		}
	}

	// Add possible source map
	extension := os.Extension(item.Path)
	file := strings.TrimSuffix(item.Path, extension)
	file = file + ".js.map"

	related = append(related, compactor.Related{
		Type:       "source-map",
		Dependency: true,
		Source:     "",
		Path:       os.File(file),
		Item:       compactor.Get(file),
	})

	// Add possible type declaration
	file = item.Path + ".d"
	related = append(related, compactor.Related{
		Type:       "declaration",
		Dependency: true,
		Source:     "",
		Path:       os.File(file),
		Item:       compactor.Get(file),
	})

	// Detect imports
	regex := regexp.MustCompile(`import ?((.+) ?from ?)?("(.+)"|'(.+)');?`)
	matches := regex.FindAllStringSubmatch(item.Content, -1)

	for _, match := range matches {
		source := match[0]
		path := strings.Trim(match[3], `'"`)
		thePath := FindRealPath(path)

		file := os.Resolve(thePath, os.Dir(item.Path))
		related = append(related, compactor.Related{
			Type:       "import",
			Dependency: false,
			Source:     source,
			Path:       path,
			Item:       compactor.Get(file),
		})
	}

	return related, nil
}

// Resolve processor
func Resolve(path string, item *compactor.Item) (string, error) {

	path = FindRealPath(path)
	file := os.Resolve(path, os.Dir(item.Path))

	bundle := compactor.GetBundle(file)
	hash := bundle.Item.Checksum

	destination := bundle.ToDestination(bundle.Item.Path)
	destination = bundle.ToHashed(destination, hash)
	destination = bundle.ToExtension(destination, ".js")
	destination = bundle.CleanPath(destination)
	destination = "/" + destination

	return destination, nil
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

	// Append extension if not declared
	if os.Extension(path) == "" {
		path += ".ts"
	}

	return path
}

// Execute processor
func Execute(bundle *compactor.Bundle) error {

	// Copy from user config file and make sure output is present
	config := *_tsConfig
	config.CompilerOptions["emitDeclarationOnly"] = false
	config.CompilerOptions["noEmit"] = false
	config.CompilerOptions["noEmitOnError"] = true

	// To compile correctly we need to force noLib, noResolve and isolatedModules
	config.CompilerOptions["noLib"] = true
	config.CompilerOptions["noResolve"] = true
	config.CompilerOptions["isolatedModules"] = true

	// Enable source maps
	if bundle.ShouldGenerateSourceMap(bundle.Item.Path) {
		config.CompilerOptions["sourceMap"] = true
		config.CompilerOptions["inlineSources"] = true
	}

	// Run transpilation
	hash := bundle.Item.Checksum
	destination := bundle.ToDestination(bundle.Item.Path)
	destination = bundle.ToHashed(destination, hash)
	destination = bundle.ToExtension(destination, ".js")

	err := RunTranspiler(Transpiler{
		File:        bundle.Item.Path,
		Content:     bundle.Item.Content,
		Options:     &config,
		Destination: destination,
	})

	if err != nil {
		return err
	}

	// Update paths after transpile code with correct final destinations
	for _, related := range bundle.Item.Related {
		if related.Item.Exists && related.Type == "import" {

			fromPath := related.Path
			toPath, err := Resolve(fromPath, bundle.Item)

			if err != nil {
				return err
			}

			oldSource := related.Source
			firstIndex := strings.LastIndex(oldSource, fromPath)
			lastIndex := firstIndex + len(fromPath)
			newSource := oldSource[:firstIndex] + toPath + oldSource[lastIndex:]

			err = os.Replace(
				destination,
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
func Optimize(bundle *compactor.Bundle) error {

	if !bundle.ShouldCompress(bundle.Item.Path) {
		return nil
	}

	hash := bundle.Item.Checksum
	destination := bundle.ToDestination(bundle.Item.Path)
	destination = bundle.ToHashed(destination, hash)
	destination = bundle.ToExtension(destination, ".js")

	args := []string{
		destination,
		"--output", destination,
		"--compress",
		"--comments",
	}

	if bundle.ShouldGenerateSourceMap(bundle.Item.Path) {

		file := os.File(destination)
		sourceOptions := strings.Join([]string{
			"includeSources",
			"filename='" + file + ".map'",
			"url='" + file + ".map'",
			"content='" + destination + ".map'",
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
		Extensions: []string{".ts", ".tsx", ".mts", ".js", ".jsx", ".mjs"},
		Init:       Init,
		Related:    Related,
		Resolve:    Resolve,
		Execute:    Execute,
		Optimize:   Optimize,
		Delete:     generic.Delete,
	}
}
