package html

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/pkg/generic"
)

// HTML minify
func Minify(content string) (string, error) {

	file, err := ioutil.TempFile("", "minify.*.html")

	if err != nil {
		return content, err
	}

	defer os.Remove(file.Name())

	_, err = file.WriteString(content)

	if err != nil {
		return content, err
	}

	_, err = compactor.ExecCommand(
		"html-minifier",
		"--output", file.Name(),
		"--collapse-whitespace",
		"--conservative-collapse",
		"--remove-comments",
		"--remove-script-type-attributes",
		"--remove-style-link-type-attributes",
		"--use-short-doctype",
		"--minify-urls", "true",
		"--minify-css", "true",
		"--minify-js", "true",
		"--ignore-custom-fragments", "/{{[{]?(.*?)[}]?}}/",
		file.Name(),
	)

	if err != nil {
		return content, err
	}

	content, err = compactor.ReadFile(file.Name())

	return content, err
}

// HTML ReplaceFormats method
func ReplaceFormats(content string) string {

	content = strings.ReplaceAll(content, ".scss", ".css")
	content = strings.ReplaceAll(content, ".sass", ".css")
	content = strings.ReplaceAll(content, ".ts", ".js")
	content = strings.ReplaceAll(content, ".tsx", ".js")

	return content
}

// HTML processor
func Processor(action *compactor.Action, bundle *compactor.Bundle, logger *compactor.Logger) error {

	if action.IsDelete() {
		return generic.DeleteProcessor(bundle, logger, []string{})
	}

	files := bundle.GetFiles()

	if bundle.IsToMultipleDestinations() {

		for _, file := range files {

			content, perm, err := compactor.ReadFileAndPermission(file)

			if err != nil {
				return err
			}

			content = ReplaceFormats(content)

			if bundle.ShouldCompress(file) {
				content, err = Minify(content)
				if err != nil {
					return err
				}
			}

			destination := bundle.ToDestination(file)
			err = compactor.WriteFile(destination, content, perm)

			if err != nil {
				return err
			}

			logger.AddProcessed(file)

		}

		return nil
	}

	content, perm, err := compactor.ReadFilesAndPermission(files)

	if err != nil {
		return err
	}

	content = ReplaceFormats(content)
	destination := bundle.GetDestination()

	if bundle.ShouldCompress(destination) {
		content, err = Minify(content)
		if err != nil {
			return err
		}
	}

	err = compactor.WriteFile(destination, content, perm)

	if err == nil {
		logger.AddWritten(destination)
	}

	return err
}
