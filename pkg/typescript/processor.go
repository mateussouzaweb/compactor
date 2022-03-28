package typescript

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/javascript"
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
	var found []compactor.Related
	return found, nil
}

// Find user defined TypeScript config file
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

// Rename destination file
func RenameDestination(from string, to string) error {

	err := os.Rename(from, to)

	if err != nil {
		return err
	}

	if os.Exist(from + ".map") {

		err = os.Rename(from+".map", to+".map")

		if err != nil {
			return err
		}

		fromName := os.File(from)
		toName := os.File(to)

		err = os.Replace(to,
			"sourceMappingURL="+fromName+".map",
			"sourceMappingURL="+toName+".map",
		)

		if err != nil {
			return err
		}

		err = os.Replace(to+".map",
			"\"file\":\""+fromName+"\"",
			"\"file\":\""+toName+"\"",
		)

	}

	return err
}

// Execute processor
func Execute(bundle *compactor.Bundle) error {

	hash, err := bundle.GetChecksum()

	if err != nil {
		return err
	}

	destination := bundle.ToDestination(bundle.Item.Path)
	destination = bundle.ToHashed(destination, hash)
	destination = bundle.ToExtension(destination, ".js")

	files := []string{bundle.Item.Path}

	for _, related := range bundle.Item.Related {
		if related.Item.Exists && related.Type == "import" {
			files = append(files, related.Item.Path)
		}
		if related.Item.Exists && related.Type == "require" {
			files = append(files, related.Item.Path)
		}
	}

	config := TSConfig{
		CompilerOptions: make(map[string]interface{}),
		Exclude:         make([]string, 0),
		Extends:         FindConfig(bundle.Source.Path),
		Files:           make([]string, 0),
		Include:         files,
		References:      make([]string, 0),
	}

	config.CompilerOptions["outDir"] = os.Dir(destination)
	config.CompilerOptions["removeComments"] = true
	config.CompilerOptions["skipLibCheck"] = true
	config.CompilerOptions["emitDeclarationOnly"] = false
	config.CompilerOptions["noEmit"] = false

	if bundle.ShouldGenerateSourceMap(bundle.Item.Path) {
		config.CompilerOptions["sourceMap"] = true
		config.CompilerOptions["inlineSources"] = true
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

	// Rename file to hashed version if necessary
	output := bundle.ToNonHashed(destination, hash)

	if destination != output {
		err = RenameDestination(output, destination)

		if err != nil {
			return err
		}
	}

	bundle.Processed(bundle.Item.Path)

	return nil
}

// Optimize processor
func Optimize(bundle *compactor.Bundle) error {

	if !bundle.ShouldCompress(bundle.Item.Path) {
		return nil
	}

	hash, err := bundle.GetChecksum()

	if err != nil {
		return err
	}

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

	_, err = os.Exec(
		"terser",
		args...,
	)

	bundle.Optimized(bundle.Item.Path)

	return err
}

// Plugin return the compactor plugin instance
func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Namespace:  "typescript",
		Extensions: []string{".ts", ".tsx", ".mts", ".js", ".jsx", ".mjs"},
		Init:       Init,
		Related:    Related,
		Execute:    Execute,
		Optimize:   Optimize,
		Delete:     javascript.Delete,
		Resolve:    javascript.Resolve,
	}
}
