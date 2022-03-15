package typescript

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
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
func InitProcessor(bundle *compactor.Bundle) error {

	err := os.NodeRequire("tsc", "typescript")

	if err != nil {
		return err
	}

	return os.NodeRequire("terser", "terser")
}

// Find user defined TypeScript config file
func FindConfig(path string) string {

	if len(path) <= 1 {
		return ""
	}
	if os.Exist(filepath.Join(path, "jsconfig.json")) {
		return filepath.Join(path, "jsconfig.json")
	}
	if os.Exist(filepath.Join(path, "tsconfig.json")) {
		return filepath.Join(path, "tsconfig.json")
	}

	return FindConfig(os.Dir(path))
}

// Typescript processor
func RunProcessor(bundle *compactor.Bundle) error {

	userConfig := FindConfig(bundle.Source.Path)

	// TODO: to multiple, simulate a typescript file with requires/imports
	for _, item := range bundle.Items {

		if !item.Exists {
			continue
		}

		destination := bundle.ToDestination(item.Path)
		destination = bundle.ToHashed(destination, item.Checksum)
		destination = bundle.ToExtension(destination, ".js")

		config := TSConfig{
			CompilerOptions: make(map[string]interface{}),
			Exclude:         make([]string, 0),
			Extends:         userConfig,
			Files:           make([]string, 0),
			Include:         []string{bundle.ToSource(item.Path)},
			References:      make([]string, 0),
		}

		config.CompilerOptions["outFile"] = destination
		config.CompilerOptions["removeComments"] = true
		config.CompilerOptions["skipLibCheck"] = true
		config.CompilerOptions["emitDeclarationOnly"] = false
		config.CompilerOptions["noEmit"] = false

		if bundle.ShouldGenerateSourceMap(item.Path) {
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

		// Compress
		if bundle.ShouldCompress(item.Path) {

			args = []string{
				destination,
				"--output", destination,
				"--compress",
				"--comments",
			}

			if bundle.ShouldGenerateSourceMap(item.Path) {

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

			if err != nil {
				return err
			}

		}

		bundle.Processed(item.Path)

	}

	return nil
}

// DeleteProcessor
func DeleteProcessor(bundle *compactor.Bundle) error {

	err := generic.DeleteProcessor(bundle)

	if err != nil {
		return err
	}

	for _, deleted := range bundle.Logs.Deleted {

		extra := bundle.ToExtension(deleted, ".js.map")

		if !os.Exist(extra) {
			continue
		}

		err := os.Delete(extra)
		if err != nil {
			return err
		}

	}

	return err
}

func Plugin() *compactor.Plugin {
	return &compactor.Plugin{
		Extensions: []string{".ts", ".tsx"},
		Init:       InitProcessor,
		Run:        RunProcessor,
		Delete:     javascript.DeleteProcessor,
		Resolve:    javascript.ResolveProcessor,
	}
}
