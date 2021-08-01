package typescript

import (
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/os"
	"github.com/mateussouzaweb/compactor/pkg/generic"
	"github.com/mateussouzaweb/compactor/pkg/javascript"
)

// Init processor
func InitProcessor(bundle *compactor.Bundle) error {

	err := os.NodeRequire("tsc", "typescript")

	if err != nil {
		return err
	}

	return os.NodeRequire("uglifyjs", "uglify-js")
}

// Typescript processor
func RunProcessor(bundle *compactor.Bundle) error {

	// TODO: to multiple, simulate a typescript file with requires/imports
	for _, item := range bundle.Items {

		if !item.Exists {
			continue
		}

		destination := bundle.ToDestination(item.Path)
		destination = bundle.ToHashed(destination, item.Checksum)
		destination = bundle.ToExtension(destination, ".js")

		args := []string{
			item.Path,
			"--outFile", destination,
			"--target", "ES2017",
			"--module", "None",
			"--allowUmdGlobalAccess",
			"--skipLibCheck",
			"--allowJs",
			"--removeComments",
		}

		if bundle.ShouldGenerateSourceMap(item.Path) {
			args = append(args, "--sourceMap", "--inlineSources")
		}

		// Compile
		_, err := os.Exec(
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
				"uglifyjs",
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
