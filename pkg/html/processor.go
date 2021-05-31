package html

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/mateussouzaweb/compactor/compactor"
	"github.com/mateussouzaweb/compactor/pkg/css"
	"github.com/mateussouzaweb/compactor/pkg/generic"
	"github.com/mateussouzaweb/compactor/pkg/javascript"
	"github.com/mateussouzaweb/compactor/pkg/sass"
	"github.com/mateussouzaweb/compactor/pkg/typescript"
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

// ExtractAttribute find the value of the attribute
func ExtractAttribute(html string, attribute string, defaultValue string) string {

	regex := regexp.MustCompile(attribute + `=[\"\']([^"']*)[\"\']`)
	match := regex.FindStringSubmatch(html)
	value := ""

	if match != nil {
		value = match[1]
	}

	if value == "" {
		value = defaultValue
	}

	return value
}

// HTML Format method
func Format(content string) (string, error) {

	var err error
	var file string

	script := `(?i)<script(.+)?>(.+)?</script>`
	regex := regexp.MustCompile(script)
	matches := regex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {

		code := match[0]
		src := ExtractAttribute(code, "src", "")

		// Ignore protocol scripts, only handle relative and absolute paths
		if strings.Contains(src, "://") {
			continue
		}

		file = src
		extension := compactor.CleanExtension(src)

		if extension == "js" {
			file, err = javascript.CorrectPath(src)
		} else {
			file, err = typescript.CorrectPath(src)
		}

		if err != nil {
			return content, err
		}

		content = strings.Replace(content, src, file, 1)

	}

	link := `(?i)<link(.+)?\/?>`
	regex = regexp.MustCompile(link)
	matches = regex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {

		code := match[0]
		rel := ExtractAttribute(code, "rel", "")
		href := ExtractAttribute(code, "href", "")
		as := ExtractAttribute(code, "as", "")

		if rel == "" || href == "" {
			continue
		}
		if rel != "stylesheet" && (rel == "preload" && as == "") {
			continue
		}

		file = href
		extension := compactor.CleanExtension(href)

		if extension == "css" {
			file, err = css.CorrectPath(href)
		} else if extension == "sass" || extension == "scss" {
			file, err = sass.CorrectPath(href)
		} else if extension == "js" {
			file, err = javascript.CorrectPath(href)
		} else if extension == "ts" {
			file, err = typescript.CorrectPath(href)
		}

		if err != nil {
			return content, err
		}

		content = strings.Replace(content, href, file, 1)

	}

	return content, nil
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

			content, err = Format(content)

			if err != nil {
				return err
			}

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

	content, err = Format(content)

	if err != nil {
		return err
	}

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
