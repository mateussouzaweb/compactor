package typescript

import (
	"encoding/json"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// TSConfig struct
type TSConfig struct {
	CompilerOptions map[string]interface{} `json:"compilerOptions,omitempty"`
	Exclude         []string               `json:"exclude,omitempty"`
	Extends         string                 `json:"extends,omitempty"`
	Files           []string               `json:"files,omitempty"`
	Include         []string               `json:"include,omitempty"`
	References      []string               `json:"references,omitempty"`
}

// Init processor
func Init(bundle *compactor.Bundle) error {

	err := os.NodeRequire("tsc", "typescript")

	if err != nil {
		return err
	}

	return os.NodeRequire("terser", "terser")
}

// Related processor
func Related(item *compactor.Item) ([]compactor.Related, error) {

	var related []compactor.Related

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
	regex := regexp.MustCompile(`import (.+) from ("(.+)"|'(.+)');?`)
	matches := regex.FindAllStringSubmatch(item.Content, -1)

	for _, match := range matches {
		source := match[0]
		sourcePath := strings.Trim(match[2], `'"`)

		path := sourcePath
		if os.Extension(path) == "" {
			path += item.Extension
		}

		file := os.Resolve(path, os.Dir(item.Path))

		if !os.Exist(file) {
			path = strings.TrimSuffix(path, os.Extension(path))
			path = path + item.Extension
			file = os.Resolve(path, os.Dir(item.Path))
		}

		related = append(related, compactor.Related{
			Type:       "import",
			Dependency: false,
			Source:     source,
			Path:       sourcePath,
			Item:       compactor.Get(file),
		})
	}

	return related, nil
}

// FindConfig locate the user defined TypeScript config file
func FindConfig(path string) string {

	if os.Exist(filepath.Join(path, "jsconfig.json")) {
		return filepath.Join(path, "jsconfig.json")
	}
	if os.Exist(filepath.Join(path, "tsconfig.json")) {
		return filepath.Join(path, "tsconfig.json")
	}
	if len(path) <= 1 {
		return ""
	}

	return FindConfig(os.Dir(path))
}

// FindFiles retrieve the final list of processable items
func FindFiles(item *compactor.Item) []string {

	files := []string{item.Path}
	result := []string{}
	found := make(map[string]bool)

	for _, related := range item.Related {
		if related.Item.Exists && related.Type == "import" {
			files = append(files, related.Item.Path)
			files = append(files, FindFiles(related.Item)...)
		}
	}

	for _, file := range files {
		if _, ok := found[file]; !ok {
			found[file] = true
			result = append(result, file)
		}
	}

	return result
}

// Execute processor
func Execute(bundle *compactor.Bundle) error {

	files := FindFiles(bundle.Item)
	compilerOptions := make(map[string]interface{})

	// Make sure output is present and set destination
	compilerOptions["baseUrl"] = bundle.ToSource(bundle.Item.Folder)
	compilerOptions["outDir"] = bundle.ToDestination(bundle.Item.Folder)
	compilerOptions["skipLibCheck"] = true
	compilerOptions["emitDeclarationOnly"] = false
	compilerOptions["noEmit"] = false
	compilerOptions["noEmitOnError"] = true

	// To compile correctly we need to force isolatedModules and noResolve
	compilerOptions["noResolve"] = true
	compilerOptions["isolatedModules"] = true

	if bundle.ShouldGenerateSourceMap(bundle.Item.Path) {
		compilerOptions["sourceMap"] = true
		compilerOptions["inlineSources"] = true
	}

	config := TSConfig{
		CompilerOptions: compilerOptions,
		Exclude:         make([]string, 0),
		Extends:         FindConfig(bundle.Source.Path),
		Files:           make([]string, 0),
		Include:         files,
		References:      make([]string, 0),
	}

	configJson, err := json.Marshal(config)

	if err != nil {
		return err
	}

	configFile := os.TemporaryFile("tsconfig.json")
	defer os.Delete(configFile)

	err = os.Write(configFile, string(configJson), 0775)

	if err != nil {
		return err
	}

	// Compile
	args := []string{
		"--build",
		configFile,
	}

	_, err = os.Exec(
		"tsc",
		args...,
	)

	if err != nil {
		return err
	}

	// Rename files to hashed version if necessary
	err = RenameDestination(bundle)

	if err != nil {
		return err
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

// Resolve processor
func Resolve(path string, item *compactor.Item) (string, error) {

	destination, err := generic.Resolve(path, item)

	if err != nil {
		return destination, err
	}

	bundle := compactor.GetBundle(path)
	hash := bundle.Item.Checksum

	destination = bundle.ToHashed(destination, hash)
	destination = bundle.ToExtension(destination, ".js")

	return destination, nil
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
