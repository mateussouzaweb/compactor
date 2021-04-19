package html

import (
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
)

// HTML processor
func Processor(context *compactor.Context, options *compactor.Options) error {

	content, perm, err := compactor.ReadFileAndPermission(context.Source)

	if err != nil {
		return err
	}

	content = strings.ReplaceAll(content, ".scss", ".css")
	content = strings.ReplaceAll(content, ".sass", ".css")
	content = strings.ReplaceAll(content, ".ts", ".js")
	content = strings.ReplaceAll(content, ".tsx", ".js")

	err = compactor.WriteFile(context.Destination, content, perm)

	if err != nil {
		return err
	}

	if options.ShouldCompress(context) {
		_, err = compactor.ExecCommand(
			"html-minifier",
			"--output", context.Destination,
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
			context.Destination,
		)
	}

	if err == nil {
		context.Processed = true
	}

	return err
}
