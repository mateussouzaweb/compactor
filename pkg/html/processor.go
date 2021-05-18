package html

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
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

// HTML processor
func Processor(bundle *compactor.Bundle, logger *compactor.Logger) error {

	files := bundle.GetFiles()
	target, isDir := bundle.GetDestination()
	result := ""

	for _, file := range files {

		content, err := compactor.ReadFile(file)

		if err != nil {
			return err
		}

		content = strings.ReplaceAll(content, ".scss", ".css")
		content = strings.ReplaceAll(content, ".sass", ".css")
		content = strings.ReplaceAll(content, ".ts", ".js")
		content = strings.ReplaceAll(content, ".tsx", ".js")

		if bundle.ShouldCompress(file) {
			content, err = Minify(content)
			if err != nil {
				return err
			}
		}

		if !isDir {
			result += content
			continue
		}

		destination := bundle.ToDestination(file)
		perm, err := compactor.GetPermission(file)

		if err != nil {
			return err
		}

		err = compactor.WriteFile(destination, content, perm)

		if err != nil {
			return err
		}

		logger.AddProcessed(file)

	}

	if isDir {
		return nil
	}

	perm, err := compactor.GetPermission(files[0])

	if err != nil {
		return err
	}

	err = compactor.WriteFile(target, result, perm)

	if err == nil {
		logger.AddWritten(target)
	}

	return err
}
